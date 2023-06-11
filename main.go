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

package kube_mailserver

import (
	"context"
	"flag"
	"os"

	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"go.linka.cloud/grpc/logger"
	"go.linka.cloud/k8s"
	dnsv1alpha1 "go.linka.cloud/k8s/dns/api/v1alpha1"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/kubernetes"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
	"go.linka.cloud/kube-mailserver/controllers"
	// +kubebuilder:scaffold:imports
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	// utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(k8s.AddToScheme(scheme))

	utilruntime.Must(mailv1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme

	utilruntime.Must(dnsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(cmv1.AddToScheme(scheme))
	utilruntime.Must(traefikv1alpha1.AddToScheme(scheme))
}

func Main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(logger.C(ctx).Logr())

	setupLog := ctrl.Log.WithName("setup")

	restConfig := ctrl.GetConfigOrDie()

	goClient := kubernetes.NewForConfigOrDie(restConfig)

	discoveryClient := discovery.NewDiscoveryClientForConfigOrDie(restConfig)

	// check if cert-manager, k8s-dns-manager and traefik are installed
	gvks := []GroupVersionKind{
		{
			Group:    "cert-manager.io",
			Version:  "v1",
			Kind:     "Certificate",
			Required: true,
		},
		{
			Group:    "dns.linka.cloud",
			Version:  "v1alpha1",
			Kind:     "DNSRecord",
			Required: true,
		},
		{
			Group:   "traefik.containo.us",
			Version: "v1alpha1",
			Kind:    "IngressRoute",
		},
		{
			Group:   "traefik.containo.us",
			Version: "v1alpha1",
			Kind:    "Middleware",
		},
	}
	if err := CheckGroupVersionKinds(ctrl.LoggerInto(ctx, setupLog), discoveryClient, gvks...); err != nil {
		setupLog.Error(err, "unable to check if cert-manager, k8s-dns-manager and traefik are installed")
		os.Exit(1)
	}

	mgr, err := ctrl.NewManager(restConfig, ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "ee12b95d.mail.linka.cloud",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.MailServerReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		RestConfig: restConfig,
		GoClient:   goClient,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MailServer")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
