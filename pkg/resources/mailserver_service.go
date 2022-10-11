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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerService(s *mv1alpha1.MailServer) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("mail", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "service"),
		},
		Spec: corev1.ServiceSpec{
			Selector:              Labels(s, "server"),
			Type:                  corev1.ServiceTypeLoadBalancer,
			LoadBalancerIP:        string(V(s.Spec.LoadBalancerIP)),
			LoadBalancerClass:     s.Spec.LoadBalancerClass,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeLocal,
			Ports: []corev1.ServicePort{
				{
					Name:       "smtp",
					Port:       25,
					TargetPort: intstr.IntOrString{IntVal: 25},
				},
				{
					Name:       "imap",
					Port:       143,
					TargetPort: intstr.IntOrString{IntVal: 143},
				},
				{
					Name:       "esmtp-implicit",
					Port:       465,
					TargetPort: intstr.IntOrString{IntVal: 465},
				},
				{
					Name:       "esmtp-explicit",
					Port:       587,
					TargetPort: intstr.IntOrString{IntVal: 587},
				},
				{
					Name:       "imap-implicit",
					Port:       993,
					TargetPort: intstr.IntOrString{IntVal: 993},
				},
				{
					Name:       "pop3",
					Port:       110,
					TargetPort: intstr.IntOrString{IntVal: 110},
				},
				{
					Name:       "pop3s",
					Port:       995,
					TargetPort: intstr.IntOrString{IntVal: 995},
				},
				{
					Name:       "sieve",
					Port:       4190,
					TargetPort: intstr.IntOrString{IntVal: 4190},
				},
			},
		},
	}
}
