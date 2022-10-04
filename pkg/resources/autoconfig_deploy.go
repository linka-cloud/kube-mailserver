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
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func AutoConfigDeploy(s *mailv1alpha1.MailServer) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        Normalize("autoconfig", s.Spec.Domain),
			Namespace:   s.Namespace,
			Labels:      Labels(s, "autoconfig"),
			Annotations: s.Spec.AutoConfig.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: s.Spec.Replicas,
			Strategy: s.Spec.AutoConfig.Strategy,
			Selector: &metav1.LabelSelector{
				MatchLabels: Labels(s, "autoconfig"),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      Labels(s, "autoconfig"),
					Annotations: s.Spec.AutoConfig.Annotations,
				},
				Spec: corev1.PodSpec{
					ServiceAccountName:        s.Spec.AutoConfig.ServiceAccountName,
					Affinity:                  s.Spec.AutoConfig.Affinity,
					SecurityContext:           s.Spec.AutoConfig.SecurityContext,
					TopologySpreadConstraints: s.Spec.AutoConfig.TopologySpreadConstraints,
					Tolerations:               s.Spec.AutoConfig.Tolerations,
					NodeSelector:              s.Spec.AutoConfig.NodeSelector,
					RestartPolicy:             corev1.RestartPolicyAlways,
					Containers: []corev1.Container{
						{
							Name:      "autoconfig",
							Image:     s.Spec.AutoConfig.Image,
							Resources: s.Spec.AutoConfig.Resources,
							Env: []corev1.EnvVar{
								{
									Name:  "DOMAIN",
									Value: s.Spec.Domain,
								},
								{
									Name:  "IMAP_SERVER",
									Value: "mail." + s.Spec.Domain,
								},
								{
									Name:  "SMTP_SERVER",
									Value: "mail." + s.Spec.Domain,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 1323,
									Protocol:      corev1.ProtocolTCP,
								},
							},
						},
					},
				},
			},
		},
	}
}
