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
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerCert(s *mv1alpha1.MailServer) *cmv1.Certificate {
	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize(s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "certificate"),
		},
		Spec: cmv1.CertificateSpec{
			CommonName: s.Spec.Domain,
			DNSNames: []string{
				s.Spec.Domain,
				"mail." + s.Spec.Domain,
			},
			SecretName: Normalize(s.Spec.Domain, "tls"),
			IssuerRef:  s.Spec.IssuerRef,
		},
	}
}
