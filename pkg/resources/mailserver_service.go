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
	corev1 "go.linka.cloud/k8s/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerService(s *mv1alpha1.MailServer) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("mail", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "service"),
		},
		Spec: &corev1.ServiceSpec{
			Selector:              Labels(s, "server"),
			Type:                  k8s.Ref(corev1.ServiceTypeLoadBalancer),
			LoadBalancerIP:        k8s.Ref(s.Spec.LoadBalancerIP.String()),
			LoadBalancerClass:     s.Spec.LoadBalancerClass,
			ExternalTrafficPolicy: k8s.Ref(corev1.ServiceExternalTrafficPolicyTypeLocal),
			Ports: []corev1.ServicePort{
				{
					Name:       k8s.Ref("smtp"),
					Port:       k8s.Ref[int32](25),
					TargetPort: intstr.IntOrString{IntVal: 25},
				},
				{
					Name:       k8s.Ref("imap"),
					Port:       k8s.Ref[int32](143),
					TargetPort: intstr.IntOrString{IntVal: 143},
				},
				{
					Name:       k8s.Ref("esmtp-implicit"),
					Port:       k8s.Ref[int32](465),
					TargetPort: intstr.IntOrString{IntVal: 465},
				},
				{
					Name:       k8s.Ref("esmtp-explicit"),
					Port:       k8s.Ref[int32](587),
					TargetPort: intstr.IntOrString{IntVal: 587},
				},
				{
					Name:       k8s.Ref("imap-implicit"),
					Port:       k8s.Ref[int32](993),
					TargetPort: intstr.IntOrString{IntVal: 993},
				},
				{
					Name:       k8s.Ref("pop3"),
					Port:       k8s.Ref[int32](110),
					TargetPort: intstr.IntOrString{IntVal: 110},
				},
				{
					Name:       k8s.Ref("pop3s"),
					Port:       k8s.Ref[int32](995),
					TargetPort: intstr.IntOrString{IntVal: 995},
				},
				{
					Name:       k8s.Ref("sieve"),
					Port:       k8s.Ref[int32](4190),
					TargetPort: intstr.IntOrString{IntVal: 4190},
				},
			},
		},
	}
}
