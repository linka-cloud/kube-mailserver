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

	corev1 "go.linka.cloud/k8s/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerCredentials(s *mv1alpha1.MailServer, password string) *corev1.Secret {
	if password == "" {
		password = RandomPassword()
	}
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("postmaster", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "credentials"),
		},
		Data: map[string][]byte{
			"email":    []byte(fmt.Sprintf("postmaster@%s", s.Spec.Domain)),
			"password": []byte(password),
		},
	}
}
