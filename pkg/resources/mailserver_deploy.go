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

	"go.linka.cloud/k8s"
	appsv1 "go.linka.cloud/k8s/apps/v1"
	corev1 "go.linka.cloud/k8s/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerDeploy(s *mv1alpha1.MailServer) *appsv1.Deployment {
	labels := Labels(s, "server")
	deploy := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        Normalize("mail", s.Spec.Domain),
			Namespace:   s.Namespace,
			Labels:      labels,
			Annotations: s.Spec.Annotations,
		},
		Spec: &appsv1.DeploymentSpec{
			Replicas: s.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: Labels(s, "server"),
			},
			Strategy: &s.Spec.Strategy,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      labels,
					Annotations: s.Spec.Annotations,
				},
				Spec: &corev1.PodSpec{
					Affinity:                  s.Spec.Affinity,
					SecurityContext:           s.Spec.SecurityContext,
					TopologySpreadConstraints: s.Spec.TopologySpreadConstraints,
					Tolerations:               s.Spec.Tolerations,
					NodeSelector:              s.Spec.NodeSelector,
					Hostname:                  k8s.Ref("mail"),
					RestartPolicy:             k8s.Ref(corev1.RestartPolicyAlways),
					InitContainers: []corev1.Container{
						{
							Name:            k8s.Ref("setup"),
							Image:           &s.Spec.Image,
							ImagePullPolicy: k8s.Ref(corev1.PullIfNotPresent),
							Command:         []string{"/bin/bash"},
							Args:            []string{"-c", fmt.Sprintf(`( listmailuser|grep -s $POSTMASTER_EMAIL || (echo "Creating Postmaster email $POSTMASTER_EMAIL" && addmailuser $POSTMASTER_EMAIL $POSTMASTER_PASSWORD)) && ( test -f /tmp/docker-mailserver/opendkim/keys/${MAIL_DOMAIN}/mail.private || (echo "Generating DKIM Private Key" && open-dkim) )`)},
							Env: append(
								[]corev1.EnvVar{
									{
										Name: k8s.Ref("POSTMASTER_EMAIL"),
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: k8s.Ref(Normalize("postmaster", s.Spec.Domain)),
												},
												Key: k8s.Ref("email"),
											},
										},
									},
									{
										Name: k8s.Ref("POSTMASTER_PASSWORD"),
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: k8s.Ref(Normalize("postmaster", s.Spec.Domain)),
												},
												Key: k8s.Ref("password"),
											},
										},
									},
									{
										Name:  k8s.Ref("MAIL_DOMAIN"),
										Value: &s.Spec.Domain,
									},
								},
								s.Spec.Env...,
							),
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: k8s.Ref(Normalize("config", s.Spec.Domain)),
										},
									},
								},
							},
							Resources:       &s.Spec.Resources,
							SecurityContext: &mailServerDeploySecurityContext,
							VolumeMounts:    mailServerDeployVolumeMounts(s),
						},
					},
					Containers: []corev1.Container{
						{
							Name:            k8s.Ref("mailserver"),
							Image:           &s.Spec.Image,
							ImagePullPolicy: k8s.Ref(corev1.PullIfNotPresent),
							// Resources:       &corev1.ResourceRequirements{},
							EnvFrom: []corev1.EnvFromSource{
								{
									SecretRef: &corev1.SecretEnvSource{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: k8s.Ref(Normalize("config", s.Spec.Domain)),
										},
									},
								},
							},
							Env: []corev1.EnvVar{
								{
									Name:  k8s.Ref("SSL_TYPE"),
									Value: k8s.Ref("manual"),
								},
								{
									Name:  k8s.Ref("SSL_CERT_PATH"),
									Value: k8s.Ref("/etc/mailserver/ssl/tls.crt"),
								},
								{
									Name:  k8s.Ref("SSL_KEY_PATH"),
									Value: k8s.Ref("/etc/mailserver/ssl/tls.key"),
								},
							},
							SecurityContext: &mailServerDeploySecurityContext,
							VolumeMounts:    mailServerDeployVolumeMounts(s),
							Ports: []corev1.ContainerPort{
								{
									Name:          k8s.Ref("smtp"),
									ContainerPort: k8s.Ref[int32](25),
								},
								{
									Name:          k8s.Ref("imap"),
									ContainerPort: k8s.Ref[int32](143),
								},
								{
									Name:          k8s.Ref("esmtp-implicit"),
									ContainerPort: k8s.Ref[int32](465),
								},
								{
									Name:          k8s.Ref("esmtp-explicit"),
									ContainerPort: k8s.Ref[int32](587),
								},
								{
									Name:          k8s.Ref("imap-implicit"),
									ContainerPort: k8s.Ref[int32](993),
								},
								{
									Name:          k8s.Ref("pop3"),
									ContainerPort: k8s.Ref[int32](110),
								},
								{
									Name:          k8s.Ref("pop3s"),
									ContainerPort: k8s.Ref[int32](995),
								},
								{
									Name:          k8s.Ref("sieve"),
									ContainerPort: k8s.Ref[int32](4190),
								},
							},
						},
					},
					Volumes: volumes(s),
				},
			},
		},
	}
	if s.Spec.Features.LDAP.Enabled && s.Spec.Features.LDAP.Nameserver != nil {
		deploy.Spec.Template.Spec.DNSPolicy = k8s.Ref(corev1.DNSNone)
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
)

func mailServerDeployVolumeMounts(s *mv1alpha1.MailServer) []corev1.VolumeMount {
	mounts := []corev1.VolumeMount{
		{
			Name:      k8s.Ref("mail-data"),
			MountPath: k8s.Ref("/var/mail"),
			SubPath:   k8s.Ref("volumes/maildata"),
		},
		{
			Name:      k8s.Ref("mail-data"),
			MountPath: k8s.Ref("/var/mail-state"),
			SubPath:   k8s.Ref("volumes/mailstate"),
		},
		{
			Name:      k8s.Ref("mail-data"),
			MountPath: k8s.Ref("/tmp/docker-mailserver"),
			SubPath:   k8s.Ref("config"),
		},
		{
			Name:      k8s.Ref("certs"),
			MountPath: k8s.Ref("/etc/mailserver/ssl/"),
			ReadOnly:  k8s.Ref(true),
		},
	}
	if s.Spec.Features.LDAP.Enabled {
		mounts = append(mounts, corev1.VolumeMount{
			Name:      k8s.Ref(ConfigOverride),
			MountPath: k8s.Ref("/tmp/docker-mailserver/" + PostfixGroups),
			SubPath:   k8s.Ref(PostfixGroups),
			ReadOnly:  k8s.Ref(true),
		})
	}
	return append(mounts, s.Spec.VolumeMounts...)
}

func volumes(s *mv1alpha1.MailServer) []corev1.Volume {
	vols := []corev1.Volume{
		{
			Name: k8s.Ref("certs"),
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: k8s.Ref(Normalize(s.Spec.Domain, "tls")),
				},
			},
		},
		{
			Name: k8s.Ref("mail-data"),
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: k8s.Ref(Normalize(s.Spec.Domain, "data")),
				},
			},
		},
	}
	if s.Spec.Features.LDAP.Enabled {
		vols = append(vols, corev1.Volume{
			Name: k8s.Ref(ConfigOverride),
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: k8s.Ref(Normalize(ConfigOverride, s.Spec.Domain)),
					},
				},
			},
		})
	}
	return append(vols, s.Spec.Volumes...)
}
