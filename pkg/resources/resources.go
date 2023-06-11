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

package resources

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"github.com/sirupsen/logrus"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	appsv1 "go.linka.cloud/k8s/apps/v1"
	corev1 "go.linka.cloud/k8s/core/v1"
	dnsv1alpha1 "go.linka.cloud/k8s/dns/api/v1alpha1"
	networkingv1 "go.linka.cloud/k8s/networking/v1"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

type Config struct {
	MailServer *mailv1alpha1.MailServer
	IP         string
	Password   string
	BindDN     string
	BindPW     string
}

func (config *Config) Resources() *Resources {
	if config.Password == "" {
		config.Password = RandomPassword()
	}
	var ting TraefikIngressRoutes
	if config.MailServer.Spec.Traefik != nil {
		ting = TraefikIngressRoutes{
			Route:          AutoConfigTraefikIngress(config.MailServer),
			RouteTLS:       AutoConfigTraefikIngressTLS(config.MailServer),
			Redirect2HTTPs: AutoConfigRedirectToHTTPS(config.MailServer),
		}
	}
	return &Resources{
		MailServer: &MailServerResources{
			Deployment:     MailServerDeploy(config.MailServer),
			ConfigOverride: ConfigOverrideConfigMap(config.MailServer),
			PVC:            MailServerPVC(config.MailServer),
			Service:        MailServerService(config.MailServer),
			CredsSecret:    MailServerCredentials(config.MailServer, config.Password),
			ConfigSecret:   MailServerConfigSecret(config.MailServer, config.BindDN, config.BindPW),
			Cert:           MailServerCert(config.MailServer),
			DNS: &MailServerDNS{
				A:          MailServerARecord(config.MailServer, config.IP),
				MX:         MailServerMXRecord(config.MailServer),
				DMARC:      MailServerDMARCRecord(config.MailServer),
				SPF:        MailServerSPFRecord(config.MailServer),
				DKIM:       MailServerDKIMRecord(config.MailServer),
				IMAP:       MailServerIMAPRecord(config.MailServer),
				IMAPs:      MailServerIMAPsRecord(config.MailServer),
				Submission: MailServerSubmissionRecord(config.MailServer),
				POP3:       MailServerPOP3Record(config.MailServer),
				POP3s:      MailServerPOP3sRecord(config.MailServer),
			},
		},
		AutoConfig: &AutoConfigResources{
			Cert:                 AutoConfigCert(config.MailServer),
			AutoDiscoverRecord:   AutoConfigSRVRecord(config.MailServer),
			Service:              AutoConfigService(config.MailServer),
			Deployment:           AutoConfigDeploy(config.MailServer),
			TraefikIngressRoutes: ting,
			Ingress:              AutoConfigIngress(config.MailServer),
		},
	}
}

type Resources struct {
	MailServer *MailServerResources
	AutoConfig *AutoConfigResources
}

func (r *Resources) SetSecretsHash() error {
	b, err := r.MailServer.ConfigSecret.Marshal()
	if err != nil {
		return err
	}
	h := sha256.New()
	h.Write(b)
	if r.MailServer.Deployment.Spec.Template.Annotations == nil {
		r.MailServer.Deployment.Spec.Template.Annotations = map[string]string{}
	}
	r.MailServer.Deployment.Spec.Template.Annotations["linka.cloud/kube-mailserver-config-secret-hash"] = hex.EncodeToString(h.Sum(nil))
	return nil
}

type MailServerResources struct {
	CredsSecret    *corev1.Secret
	Deployment     *appsv1.Deployment
	ConfigOverride *corev1.ConfigMap
	PVC            *corev1.PersistentVolumeClaim
	Service        *corev1.Service
	ConfigSecret   *corev1.Secret
	Cert           *cmv1.Certificate
	DNS            *MailServerDNS
}

type MailServerDNS struct {
	A          *dnsv1alpha1.DNSRecord
	MX         *dnsv1alpha1.DNSRecord
	DMARC      *dnsv1alpha1.DNSRecord
	SPF        *dnsv1alpha1.DNSRecord
	DKIM       *dnsv1alpha1.DNSRecord
	IMAP       *dnsv1alpha1.DNSRecord
	IMAPs      *dnsv1alpha1.DNSRecord
	Submission *dnsv1alpha1.DNSRecord
	POP3       *dnsv1alpha1.DNSRecord
	POP3s      *dnsv1alpha1.DNSRecord
}

type AutoConfigResources struct {
	Cert                 *cmv1.Certificate
	AutoDiscoverRecord   *dnsv1alpha1.DNSRecord
	Service              *corev1.Service
	Deployment           *appsv1.Deployment
	TraefikIngressRoutes TraefikIngressRoutes
	Ingress              *networkingv1.Ingress
}

type TraefikIngressRoutes struct {
	Route          *traefikv1alpha1.IngressRoute
	RouteTLS       *traefikv1alpha1.IngressRoute
	Redirect2HTTPs *traefikv1alpha1.Middleware
}

func RandomPassword() string {
	password := make([]byte, 32)
	if _, err := rand.Read(password); err != nil {
		logrus.Fatal(err)
	}
	return base64.RawStdEncoding.EncodeToString(password)
}
