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

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func AutoConfigCert(s *mailv1alpha1.MailServer) *cmv1.Certificate {
	return &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("autoconfig", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "autoconfig-certs"),
		},
		Spec: cmv1.CertificateSpec{
			CommonName: "autoconfig." + s.Spec.Domain,
			DNSNames: []string{
				"autoconfig." + s.Spec.Domain,
				"autodiscover." + s.Spec.Domain,
			},
			SecretName: Normalize("autoconfig", s.Spec.Domain, "tls"),
			IssuerRef:  s.Spec.IssuerRef,
		},
	}
}
