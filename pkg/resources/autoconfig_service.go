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

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func AutoConfigService(s *mailv1alpha1.MailServer) *corev1.Service {
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("autoconfig", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "service"),
		},
		Spec: &corev1.ServiceSpec{
			Selector: Labels(s, "autoconfig"),
			Ports: []corev1.ServicePort{
				{
					Name:       k8s.Ref("http"),
					Port:       k8s.Ref[int32](80),
					TargetPort: intstr.FromInt(1323),
				},
			},
		},
	}
}
