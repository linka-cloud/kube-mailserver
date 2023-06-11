// Copyright 2022 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/miekg/dns"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"go.linka.cloud/k8s"
	appsv1 "go.linka.cloud/k8s/apps/v1"
	corev1 "go.linka.cloud/k8s/core/v1"
	dnsv1alpha1 "go.linka.cloud/k8s/dns/api/v1alpha1"
	networkingv1 "go.linka.cloud/k8s/networking/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/utils/diff"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
	"go.linka.cloud/kube-mailserver/pkg/resources"
)

const (
	Finalizer = "mail.linka.cloud/finalizer"

	ownerKey = ".metadata.controller"

	restartAnnotation = "mail.linka.cloud/restart"

	owner = client.FieldOwner("kube-mailserver")
)

// MailServerReconciler reconciles a MailServer object
type MailServerReconciler struct {
	client.Client
	Scheme     *runtime.Scheme
	GoClient   *kubernetes.Clientset
	RestConfig *rest.Config
}

// +kubebuilder:rbac:groups=mail.linka.cloud,resources=mailservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=mail.linka.cloud,resources=mailservers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=mail.linka.cloud,resources=mailservers/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete;exec
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=dns.linka.cloud,resources=dnsrecords,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cert-manager.io,resources=certificates,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *MailServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var s mailv1alpha1.MailServer
	if err := r.Get(ctx, req.NamespacedName, &s); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch MailServer")
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if !s.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("deleting server")
		if err := r.ReconcileDelete(ctx, &s); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if !hasFinalizer(&s) {
		log.V(5).Info("adding finalizer")
		s.Finalizers = append(s.Finalizers, Finalizer)
		if err := r.Update(ctx, &s); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if s.Status.Domain == "" {
		s.Status.Domain = s.Spec.Domain
		if err := r.Status().Update(ctx, &s); err != nil {
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if s.Status.Domain != s.Spec.Domain {
		return ctrl.Result{}, fmt.Errorf("domain cannot be changed")
	}

	conf := resources.Config{
		MailServer: &s,
	}

	// check for ldap secret
	if s.Spec.Features.LDAP.Enabled {
		var bindCreds corev1.Secret
		if err := r.Get(ctx, client.ObjectKey{Namespace: s.Namespace, Name: s.Spec.Features.LDAP.BindSecret}, &bindCreds); err != nil {
			log.Error(err, "unable to fetch LDAP bind credentials")
			return ctrl.Result{}, err
		}
		dn, ok := bindCreds.Data["bindDN"]
		if !ok {
			return ctrl.Result{}, fmt.Errorf("bindDN not found in LDAP bind credentials secret")
		}
		password, ok := bindCreds.Data["bindPW"]
		if !ok {
			return ctrl.Result{}, fmt.Errorf("bindPW not found in LDAP bind credentials secret")
		}
		conf.BindDN = string(dn)
		conf.BindPW = string(password)
	}

	res := conf.Resources()

	if err := res.SetSecretsHash(); err != nil {
		log.Error(err, "unable to set secrets hash")
		return ctrl.Result{}, err
	}

	if r, ok, err := r.reconcileCredentials(ctx, &s, res); !ok {
		return r, err
	}

	if r, ok, err := r.reconcileResources(ctx, &s, res); !ok {
		return r, err
	}

	if r, ok, err := r.reconcileARecord(ctx, &s, res); !ok {
		return r, err
	}

	if r, ok, err := r.reconcileReplicas(ctx, &s, res); !ok {
		return r, err
	}

	// TODO(adphi): restart deploy when certificates are renewed

	// set mailserver public IP (e.g. `curl ifconfig.me`) in spf record: v=spf1 a mx ip4:$PUBLIC_IP -all
	if r, ok, err := r.reconcileSPF(ctx, &s, res); !ok {
		return r, err
	}

	// get dkim key from deployment
	if r, ok, err := r.reconcileDKIM(ctx, &s, res); !ok {
		return r, err
	}
	return ctrl.Result{}, nil
}

func (r *MailServerReconciler) ReconcileDelete(ctx context.Context, s *mailv1alpha1.MailServer) error {
	// garbage collection should handle cleaning by itself with the resource owner references
	if removeFinalizer(s) {
		return r.Update(ctx, s)
	}
	return nil
}

func (r *MailServerReconciler) reconcileCredentials(ctx context.Context, s *mailv1alpha1.MailServer, res *resources.Resources) (ctrl.Result, bool, error) {
	log := ctrl.LoggerFrom(ctx)

	ps := &corev1.Secret{}
	log.V(5).Info("looking for credentials secret", "name", res.MailServer.CredsSecret.Name)
	if err := r.Get(ctx, client.ObjectKeyFromObject(res.MailServer.CredsSecret), ps); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, false, err
		}
		ps = res.MailServer.CredsSecret
		log.Info("creating credentials secret", "name", res.MailServer.CredsSecret)
		if err := ctrl.SetControllerReference(s, ps, r.Scheme); err != nil {
			return ctrl.Result{}, false, err
		}
		if err := r.Create(ctx, ps); err != nil {
			return ctrl.Result{}, false, err
		}
	} else {
		log.V(5).Info("credentials secret already exists", "name", res.MailServer.CredsSecret)
	}
	return ctrl.Result{}, true, nil
}

func (r *MailServerReconciler) reconcileResources(ctx context.Context, s *mailv1alpha1.MailServer, res *resources.Resources) (ctrl.Result, bool, error) {
	log := ctrl.LoggerFrom(ctx)
	// generate manifests
	log.V(5).Info("generating manifests")

	var mres = []client.Object{
		res.MailServer.ConfigSecret,
		res.MailServer.Cert,
		res.MailServer.PVC,
		res.MailServer.Deployment,
		res.MailServer.ConfigOverride,
		res.MailServer.Service,
		res.MailServer.DNS.MX,
		res.MailServer.DNS.SPF,
		res.MailServer.DNS.DMARC,
		res.MailServer.DNS.IMAP,
		res.MailServer.DNS.IMAPs,
		res.MailServer.DNS.POP3,
		res.MailServer.DNS.POP3s,
		res.MailServer.DNS.Submission,
	}
	ares := []client.Object{
		res.AutoConfig.Cert,
		res.AutoConfig.Deployment,
		res.AutoConfig.Service,
		res.AutoConfig.AutoDiscoverRecord,
	}
	tres := []client.Object{
		res.AutoConfig.TraefikIngressRoutes.Redirect2HTTPs,
		res.AutoConfig.TraefikIngressRoutes.Route,
		res.AutoConfig.TraefikIngressRoutes.RouteTLS,
	}
	autoConfigEnabled := s.Spec.AutoConfig.Enabled == nil || *s.Spec.AutoConfig.Enabled
	if autoConfigEnabled {
		mres = append(mres, ares...)
		if s.Spec.Traefik != nil && s.Spec.Traefik.CRDs {
			mres = append(mres, tres...)
		} else {
			mres = append(mres, res.AutoConfig.Ingress)
		}
	}

	ok := true

	// create/update resources
	for _, v := range mres {
		// skip nil resources, e.g. traefik ingress
		if v == nil {
			continue
		}
		if v, ok := interface{}(v).(interface{ Default() }); ok {
			v.Default()
		}

		log.V(5).Info("reconciling", "resource", v.GetObjectKind().GroupVersionKind().Kind, "resourceName", v.GetName())
		// do not set owner reference for PVCs to preserve the data on deletion
		if _, ok := v.(*corev1.PersistentVolumeClaim); !ok {
			if err := ctrl.SetControllerReference(s, v, r.Scheme); err != nil {
				return ctrl.Result{}, false, err
			}
		}
		want := v.DeepCopyObject().(client.Object)
		got := v.DeepCopyObject().(client.Object)
		if err := r.Get(ctx, client.ObjectKeyFromObject(got), got); err != nil {
			log.Info("creating", "resource", v.GetObjectKind().GroupVersionKind().Kind, "resourceName", v.GetName())
			if !apierrors.IsNotFound(err) {
				return ctrl.Result{}, false, err
			}
			if err := r.Patch(ctx, v, client.Apply, owner); err != nil {
				return ctrl.Result{}, false, err
			}
			continue
		}
		if err := r.Patch(ctx, want, client.Apply, owner, client.DryRunAll); err != nil {
			return ctrl.Result{}, false, err
		}
		want.SetManagedFields(nil)
		got.SetManagedFields(nil)
		if equality.Semantic.DeepDerivative(want, got) {
			log.V(5).Info("no changes", "resource", got.GetObjectKind().GroupVersionKind().Kind, "resourceName", v.GetName())
			continue
		}
		log.Info("diff", "resource", v.GetObjectKind().GroupVersionKind().Kind, "resourceName", v.GetName())
		fmt.Println(diff.ObjectReflectDiff(got, want))
		log.Info("applying", "resource", v.GetObjectKind().GroupVersionKind().Kind, "resourceName", v.GetName())
		if err := r.Patch(ctx, v, client.Apply, owner); err != nil {
			return ctrl.Result{}, false, err
		}
	}

	if !autoConfigEnabled {
		for _, v := range ares {
			if err := r.Delete(ctx, v); err != nil {
				if client.IgnoreNotFound(err) != nil {
					log.Error(err, "unable to delete resource", "kind", v.GetObjectKind().GroupVersionKind().Kind, "name", v.GetName())
					return ctrl.Result{}, false, err
				}
			} else {
				log.Info("deleted resource", "kind", v.GetObjectKind().GroupVersionKind().Kind, "name", v.GetName())
				ok = false
			}
		}
	}
	if !ok {
		return ctrl.Result{}, false, nil
	}
	if s.Status.AutoConfig == nil || *s.Status.AutoConfig != autoConfigEnabled {
		s.Status.AutoConfig = &autoConfigEnabled
		if err := r.Status().Update(ctx, s); err != nil {
			log.Error(err, "unable to update autoconfig status")
			return ctrl.Result{}, false, err
		}
		return ctrl.Result{}, false, nil
	}
	if s.Spec.Traefik == nil || !s.Spec.Traefik.CRDs {
		for _, v := range tres {
			if err := r.Delete(ctx, v); err != nil {
				if client.IgnoreNotFound(err) != nil {
					log.Error(err, "unable to delete resource", "kind", v.GetObjectKind().GroupVersionKind().Kind, "name", v.GetName())
					return ctrl.Result{}, false, err
				}
			} else {
				log.Info("deleted resource", "kind", v.GetObjectKind().GroupVersionKind().Kind, "name", v.GetName())
				ok = false
			}
		}
	} else {
		if err := r.Delete(ctx, res.AutoConfig.Ingress); err != nil {
			if client.IgnoreNotFound(err) != nil {
				log.Error(err, "unable to delete ingress")
				return ctrl.Result{}, false, err
			}
		} else {
			log.Info("deleted ingress")
			ok = false
		}
	}
	if !ok {
		return ctrl.Result{}, false, nil
	}
	if s.Status.Traefik == nil || *s.Status.Traefik != (autoConfigEnabled && s.Spec.Traefik.CRDs) {
		s.Status.Traefik = &s.Spec.Traefik.CRDs
		if err := r.Status().Update(ctx, s); err != nil {
			log.Error(err, "unable to update traefik crds status")
			return ctrl.Result{}, false, err
		}
		return ctrl.Result{}, false, nil
	}
	return ctrl.Result{}, true, nil
}

func (r *MailServerReconciler) restartMailServer(ctx context.Context, s *mailv1alpha1.MailServer, res *resources.Resources) error {
	log := ctrl.LoggerFrom(ctx)
	log.Info("restarting mail server")
	var deploy appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKeyFromObject(res.MailServer.Deployment), &deploy); err != nil {
		return err
	}
	return r.Patch(ctx, &deploy, client.RawPatch(types.MergePatchType, []byte(fmt.Sprintf(`{"spec": {"template": {"metadata": {"annotations": {"%s": "%s"}}}}}`, restartAnnotation, time.Now().UTC().Format(time.RFC3339)))))
}

func (r *MailServerReconciler) reconcileARecord(ctx context.Context, s *mailv1alpha1.MailServer, res *resources.Resources) (ctrl.Result, bool, error) {
	log := ctrl.LoggerFrom(ctx)
	// retrieve load balancer IP
	var svc corev1.Service
	if err := r.Get(ctx, client.ObjectKeyFromObject(res.MailServer.Service), &svc); err != nil {
		return ctrl.Result{}, false, err
	}
	var ip string
	if s.Spec.OverrideIP != nil {
		ip = string(*s.Spec.OverrideIP)
	} else if len(svc.Status.LoadBalancer.Ingress) != 0 {
		ip = k8s.Value(svc.Status.LoadBalancer.Ingress[0].IP)
	} else {
		log.Error(fmt.Errorf("load balancer IP not available yet"), "waiting for load balancer IP")
		return ctrl.Result{}, false, nil
	}
	if s.Status.LoadBalancerIP != ip {
		s.Status.LoadBalancerIP = ip
		if err := r.Status().Update(ctx, s); err != nil {
			return ctrl.Result{}, false, err
		}
		return ctrl.Result{}, false, nil
	}
	rec := &dnsv1alpha1.DNSRecord{}
	if err := r.Get(ctx, client.ObjectKeyFromObject(res.MailServer.DNS.A), rec); err != nil {
		if client.IgnoreNotFound(err) != nil {
			return ctrl.Result{}, false, err
		}
		if ip != "" {
			rec = resources.MailServerARecord(s, ip)
			if err := ctrl.SetControllerReference(s, rec, r.Scheme); err != nil {
				return ctrl.Result{}, false, err
			}
			if err := r.Create(ctx, rec); err != nil {
				return ctrl.Result{}, false, err
			}
		}
		return ctrl.Result{}, false, nil
	}
	if ip != "" && (rec.Spec.A == nil || ip != rec.Spec.A.Target) {
		rec.Spec.A = resources.MailServerARecord(s, ip).Spec.A
		if err := r.Update(ctx, rec); err != nil {
			return ctrl.Result{}, false, err
		}
	}
	return ctrl.Result{}, true, nil
}

func (r *MailServerReconciler) reconcileReplicas(ctx context.Context, s *mailv1alpha1.MailServer, res *resources.Resources) (ctrl.Result, bool, error) {
	log := ctrl.LoggerFrom(ctx)
	// retrieve deployment
	var deploy appsv1.Deployment
	if err := r.Get(ctx, client.ObjectKeyFromObject(res.MailServer.Deployment), &deploy); err != nil {
		log.Error(err, "unable to fetch Deployment")
		return ctrl.Result{}, false, err
	}

	selector, err := metav1.LabelSelectorAsSelector(deploy.Spec.Selector)
	if err != nil {
		log.Error(err, "unable to retrieve Deployment labels")
		return reconcile.Result{}, false, err
	}

	if k8s.Value(deploy.Status.AvailableReplicas) != s.Status.Replicas {
		s.Status.Replicas = k8s.Value(deploy.Status.AvailableReplicas)
		s.Status.Selector = selector.String()
		if err := r.Status().Update(ctx, s); err != nil {
			log.Error(err, "unable to update status")
			return ctrl.Result{}, false, err
		}
		return ctrl.Result{}, false, nil
	}
	return ctrl.Result{}, true, nil
}

func (r *MailServerReconciler) reconcileSPF(ctx context.Context, s *mailv1alpha1.MailServer, res *resources.Resources) (ctrl.Result, bool, error) {
	log := ctrl.LoggerFrom(ctx)
	if s.Spec.SPF != "" {
		return ctrl.Result{}, true, nil
	}
	spf := "v=spf1 a mx ip4:%s -all"
	out, ok, err := r.execDeployOut(ctx, res.MailServer.Deployment, "curl ifconfig.me")
	if !ok {
		return ctrl.Result{}, false, err
	}
	var rec dnsv1alpha1.DNSRecord
	if err := r.Get(ctx, client.ObjectKeyFromObject(res.MailServer.DNS.SPF), &rec); err != nil {
		log.Error(err, "unable to fetch SPF DNSRecord")
		return ctrl.Result{}, false, err
	}
	targets := []string{fmt.Sprintf(spf, strings.TrimSpace(out))}
	if strings.Join(rec.Spec.TXT.Targets, "") == strings.Join(targets, "") {
		return ctrl.Result{}, true, nil
	}
	rec.Spec.TXT.Targets = targets
	if err := r.Update(ctx, &rec); err != nil {
		log.Error(err, "unable to update SPF DNSRecord")
		return ctrl.Result{}, false, err
	}
	return ctrl.Result{}, false, nil
}

func (r *MailServerReconciler) reconcileDKIM(ctx context.Context, s *mailv1alpha1.MailServer, res *resources.Resources) (ctrl.Result, bool, error) {
	log := ctrl.LoggerFrom(ctx)
	// retrieve dkim record
	got := &dnsv1alpha1.DNSRecord{}
	var exists bool
	if err := r.Get(ctx, client.ObjectKeyFromObject(res.MailServer.DNS.DKIM), got); err != nil {
		if client.IgnoreNotFound(err) != nil {
			log.Error(err, "unable to fetch DNSRecord")
			return ctrl.Result{}, false, err
		}
	} else {
		exists = true
	}
	out, ok, err := r.execDeployOut(ctx, res.MailServer.Deployment, fmt.Sprintf("cat /etc/opendkim/keys/%s/mail.txt", s.Spec.Domain))
	if err != nil {
		log.Error(err, "unable to retrieve dkim key")
		return ctrl.Result{}, false, err
	}
	if !ok {
		return ctrl.Result{}, false, nil
	}
	// create dns record
	rr, err := dns.NewRR(out)
	if err != nil {
		log.Error(err, "unable to parse dkim record")
		return ctrl.Result{}, false, err
	}
	if _, ok := rr.(*dns.TXT); !ok {
		log.Error(err, "dkim record is not a TXT record")
		return ctrl.Result{}, false, err
	}
	rec := resources.MailServerDKIMRecord(s)
	rec.Spec.TXT.Targets = rr.(*dns.TXT).Txt
	rec.Default()
	if got.Spec.TXT != nil && got.Spec.TXT.Name == rec.Spec.TXT.Name && strings.Join(got.Spec.TXT.Targets, " ") == strings.Join(rec.Spec.TXT.Targets, " ") {
		return ctrl.Result{}, true, nil
	}
	if err := ctrl.SetControllerReference(s, rec, r.Scheme); err != nil {
		log.Error(err, "unable to set controller reference on dkim record")
		return ctrl.Result{}, false, err
	}

	if !exists {
		if err := r.Create(ctx, rec); err != nil {
			log.Error(err, "unable to create dkim record")
			return ctrl.Result{}, false, err
		}
	} else {
		rec.SetResourceVersion(got.GetResourceVersion())
		if err := r.Update(ctx, rec); err != nil {
			log.Error(err, "unable to update dkim record")
			return ctrl.Result{}, false, err
		}
	}

	return ctrl.Result{}, false, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MailServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	res := []client.Object{
		&corev1.Secret{},
		&appsv1.Deployment{},
		&corev1.Service{},
		&corev1.PersistentVolumeClaim{},
		&networkingv1.Ingress{},
		&dnsv1alpha1.DNSRecord{},
		&cmv1.Certificate{},
		&traefikv1alpha1.IngressRoute{},
		&traefikv1alpha1.Middleware{},
	}
	c := ctrl.NewControllerManagedBy(mgr).
		For(&mailv1alpha1.MailServer{})
	for _, v := range res {
		if err := mgr.GetFieldIndexer().IndexField(context.Background(), v, ownerKey, extractValue); err != nil {
			return err
		}
		c = c.Owns(v)
	}

	return c.Complete(r)
}

func extractValue(rawObj client.Object) []string {
	// grab the owner object
	owner := metav1.GetControllerOf(rawObj)
	if owner == nil {
		return nil
	}

	if owner.APIVersion != mailv1alpha1.GroupVersion.String() || owner.Kind != "MailServer" {
		return nil
	}

	return []string{owner.Name}
}

func hasFinalizer(s *mailv1alpha1.MailServer) bool {
	for _, v := range s.Finalizers {
		if v == Finalizer {
			return true
		}
	}
	return false
}

func removeFinalizer(s *mailv1alpha1.MailServer) bool {
	for i, v := range s.ObjectMeta.Finalizers {
		if v != Finalizer {
			continue
		}
		if len(s.Finalizers) == 1 {
			s.Finalizers = nil
			return true
		}
		s.ObjectMeta.Finalizers = append(s.ObjectMeta.Finalizers[i:], s.ObjectMeta.Finalizers[i+1:]...)
		return true
	}
	return false
}
