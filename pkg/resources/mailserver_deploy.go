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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerDeploy(s *mv1alpha1.MailServer) *appsv1.Deployment {
	labels := Labels(s, "server")
	deploy := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        Normalize("mail", s.Spec.Domain),
			Namespace:   s.Namespace,
			Labels:      labels,
			Annotations: s.Spec.Annotations,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: s.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: Labels(s, "server"),
			},
			Strategy: s.Spec.Strategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: s.Spec.Annotations,
				},
				Spec: corev1.PodSpec{
					Affinity:                  s.Spec.Affinity,
					SecurityContext:           s.Spec.SecurityContext,
					TopologySpreadConstraints: s.Spec.TopologySpreadConstraints,
					Tolerations:               s.Spec.Tolerations,
					NodeSelector:              s.Spec.NodeSelector,
					Hostname:                  "mail",
					RestartPolicy:             corev1.RestartPolicyAlways,
					InitContainers: []corev1.Container{
						{
							Name:            "setup",
							Image:           s.Spec.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"/bin/bash"},
							Args:            []string{"-c", fmt.Sprintf(`( listmailuser|grep -s $POSTMASTER_EMAIL || (echo "Creating Postmaster email $POSTMASTER_EMAIL" && addmailuser $POSTMASTER_EMAIL $POSTMASTER_PASSWORD)) && ( test -f /tmp/docker-mailserver/opendkim/keys/${MAIL_DOMAIN}/mail.private || (echo "Generating DKIM Private Key" && open-dkim) )`)},
							Env: append(
								[]corev1.EnvVar{
									{
										Name: "POSTMASTER_EMAIL",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: Normalize("postmaster", s.Spec.Domain),
												},
												Key: "email",
											},
										},
									},
									{
										Name: "POSTMASTER_PASSWORD",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: Normalize("postmaster", s.Spec.Domain),
												},
												Key: "password",
											},
										},
									},
									{
										Name:  "MAIL_DOMAIN",
										Value: s.Spec.Domain,
									},
								},
								s.Spec.Env...,
							),
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: Normalize("config", s.Spec.Domain),
										},
									},
								},
							},
							Resources:       s.Spec.Resources,
							SecurityContext: &mailServerDeploySecurityContext,
							VolumeMounts:    mailServerDeployVolumeMounts,
						},
					},
					Containers: []corev1.Container{
						{
							Name:            "mailserver",
							Image:           s.Spec.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Resources:       corev1.ResourceRequirements{},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: Normalize("config", s.Spec.Domain),
										},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  "SSL_TYPE",
									Value: "manual",
								},
								{
									Name:  "SSL_CERT_PATH",
									Value: "/etc/mailserver/ssl/tls.crt",
								},
								{
									Name:  "SSL_KEY_PATH",
									Value: "/etc/mailserver/ssl/tls.key",
								},
							},
							SecurityContext: &mailServerDeploySecurityContext,
							VolumeMounts:    mailServerDeployVolumeMounts,
							Ports: []corev1.ContainerPort{
								{
									Name:          "smtp",
									ContainerPort: 25,
								},
								{
									Name:          "imap",
									ContainerPort: 143,
								},
								{
									Name:          "esmtp-implicit",
									ContainerPort: 465,
								},
								{
									Name:          "esmtp-explicit",
									ContainerPort: 587,
								},
								{
									Name:          "imap-implicit",
									ContainerPort: 993,
								},
								{
									Name:          "pop3",
									ContainerPort: 110,
								},
								{
									Name:          "pop3s",
									ContainerPort: 995,
								},
								{
									Name:          "sieve",
									ContainerPort: 4190,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "certs",
							VolumeSource: corev1.VolumeSource{
								Secret: &corev1.SecretVolumeSource{
									SecretName: Normalize(s.Spec.Domain, "tls"),
								},
							},
						},
						{
							Name: "mail-data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: Normalize(s.Spec.Domain, "data"),
								},
							},
						},
					},
				},
			},
		},
	}
	if s.Spec.Features.LDAP.Enabled && s.Spec.Features.LDAP.Nameserver != nil {
		deploy.Spec.Template.Spec.DNSPolicy = corev1.DNSNone
		deploy.Spec.Template.Spec.DNSConfig = &corev1.PodDNSConfig{
			Nameservers: []string{string(*s.Spec.Features.LDAP.Nameserver)},
		}
	}
	return deploy
}

var (
	mailServerDeploySecurityContext = corev1.SecurityContext{
		AllowPrivilegeEscalation: P(false),
		ReadOnlyRootFilesystem:   P(false),
		RunAsUser:                P(int64(0)),
		RunAsGroup:               P(int64(0)),
		RunAsNonRoot:             P(false),
		Privileged:               P(false),
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{
				// file permission capabilities
				"CHOWN",
				"FOWNER",
				"MKNOD",
				"SETGID",
				"SETUID",
				"DAC_OVERRIDE",
				// network capabilities
				"NET_ADMIN", // needed for F2B
				"NET_RAW",   // needed for F2B
				"NET_BIND_SERVICE",
				// miscellaneous  capabilities
				"SYS_CHROOT",
				"SYS_PTRACE",
				"KILL",
			},
			Drop: []corev1.Capability{"ALL"},
		},
	}

	mailServerDeployVolumeMounts = []corev1.VolumeMount{
		{
			Name:      "mail-data",
			MountPath: "/var/mail",
			SubPath:   "volumes/maildata",
		},
		{
			Name:      "mail-data",
			MountPath: "/var/mail-state",
			SubPath:   "volumes/mailstate",
		},
		{
			Name:      "mail-data",
			MountPath: "/tmp/docker-mailserver",
			SubPath:   "config",
		},
		{
			Name:      "certs",
			MountPath: "/etc/mailserver/ssl/",
		},
	}
)
