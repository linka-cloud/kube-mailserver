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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerPVC(s *mv1alpha1.MailServer) *corev1.PersistentVolumeClaim {
	var storageClassName *string
	if s.Spec.Volume.StorageClass != "" {
		storageClassName = &s.Spec.Volume.StorageClass
	}
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize(s.Spec.Domain, "data"),
			Namespace: s.Namespace,
			Labels:    Labels(s, "storage"),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			StorageClassName: storageClassName,
			AccessModes: []corev1.PersistentVolumeAccessMode{
				s.Spec.Volume.AccessMode,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(s.Spec.Volume.Size),
				},
			},
			VolumeMode: P(corev1.PersistentVolumeFilesystem),
		},
	}
}
