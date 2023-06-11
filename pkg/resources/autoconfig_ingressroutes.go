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
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
	traefikv1alpha1 "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func AutoConfigTraefikIngressTLS(s *mailv1alpha1.MailServer) *traefikv1alpha1.IngressRoute {
	return &traefikv1alpha1.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "traefik.containo.us/v1alpha1",
			Kind:       "IngressRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("autoconfig", s.Spec.Domain, "tls"),
			Namespace: s.Namespace,
			Labels:    Labels(s, "https-route"),
		},
		Spec: traefikv1alpha1.IngressRouteSpec{
			EntryPoints: []string{
				s.Spec.Traefik.Entrypoints.HTTPS,
			},
			Routes: []traefikv1alpha1.Route{
				{
					Match: "Host(`autoconfig." + s.Spec.Domain + "`) || Host(`autodiscover." + s.Spec.Domain + "`)",
					Kind:  "Rule",
					Services: []traefikv1alpha1.Service{
						{
							LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
								Name: Normalize("autoconfig", s.Spec.Domain),
								Port: intstr.FromInt(80),
							},
						},
					},
				},
			},
			TLS: &traefikv1alpha1.TLS{
				SecretName: Normalize("autoconfig", s.Spec.Domain, "tls"),
			},
		},
	}
}

func AutoConfigTraefikIngress(s *mailv1alpha1.MailServer) *traefikv1alpha1.IngressRoute {
	return &traefikv1alpha1.IngressRoute{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "traefik.containo.us/v1alpha1",
			Kind:       "IngressRoute",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("autoconfig", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "http-route"),
		},
		Spec: traefikv1alpha1.IngressRouteSpec{
			EntryPoints: []string{
				s.Spec.Traefik.Entrypoints.HTTP,
			},
			Routes: []traefikv1alpha1.Route{
				{
					Match: "Host(`autoconfig." + s.Spec.Domain + "`) || Host(`autodiscover." + s.Spec.Domain + "`)",
					Kind:  "Rule",
					Middlewares: []traefikv1alpha1.MiddlewareRef{
						{
							Name: Normalize("redirect-to-https", s.Spec.Domain),
						},
					},
					Services: []traefikv1alpha1.Service{
						{
							LoadBalancerSpec: traefikv1alpha1.LoadBalancerSpec{
								Name:           Normalize("autoconfig", s.Spec.Domain),
								Port:           intstr.FromInt(80),
								PassHostHeader: P(true),
							},
						},
					},
				},
			},
		},
	}
}

func AutoConfigRedirectToHTTPS(s *mailv1alpha1.MailServer) *traefikv1alpha1.Middleware {
	return &traefikv1alpha1.Middleware{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "traefik.containo.us/v1alpha1",
			Kind:       "Middleware",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("redirect-to-https", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "redirect-to-https"),
		},
		Spec: traefikv1alpha1.MiddlewareSpec{
			RedirectScheme: &dynamic.RedirectScheme{
				Scheme:    "https",
				Permanent: true,
			},
		},
	}
}
