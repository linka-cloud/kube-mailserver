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
	"go.linka.cloud/k8s"
	networkingv1 "go.linka.cloud/k8s/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

const (
	TraefikIngress                  = "traefik.ingress.kubernetes.io"
	TraefikIngressRouter            = TraefikIngress + "/router"
	TraefikIngressRouterEntryPoints = TraefikIngressRouter + ".entrypoints"
	TraefikIngressRouterTLS         = TraefikIngressRouter + ".tls"
)

func AutoConfigIngress(s *mailv1alpha1.MailServer) *networkingv1.Ingress {
	var hosts []string
	var rules []networkingv1.IngressRule
	annotations := make(map[string]string)
	if s.Spec.Traefik != nil {
		annotations[TraefikIngressRouterEntryPoints] = s.Spec.Traefik.Entrypoints.HTTPS
		annotations[TraefikIngressRouterTLS] = "true"
	}
	for k, v := range s.Spec.AutoConfig.Ingress.Annotations {
		annotations[k] = v
	}
	for _, v := range []string{"autoconfig.", "autodiscover."} {
		host := v + s.Spec.Domain
		hosts = append(hosts, host)
		rules = append(rules, networkingv1.IngressRule{
			Host: &host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{
						{
							Path:     k8s.Ref("/"),
							PathType: k8s.Ref(networkingv1.PathTypeImplementationSpecific),
							Backend: &networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: k8s.Ref(Normalize("autoconfig", s.Spec.Domain)),
									Port: &networkingv1.ServiceBackendPort{
										Number: k8s.Ref[int32](80),
									},
								},
							},
						},
					},
				},
			},
		})
	}
	return &networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        Normalize("autoconfig", s.Spec.Domain),
			Namespace:   s.Namespace,
			Labels:      Labels(s, "ingress"),
			Annotations: annotations,
		},
		Spec: &networkingv1.IngressSpec{
			Rules: rules,
			TLS: []networkingv1.IngressTLS{
				{
					Hosts:      hosts,
					SecretName: k8s.Ref(Normalize("autoconfig", s.Spec.Domain, "tls")),
				},
			},
		},
	}
}
