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
	"github.com/miekg/dns"
	dnsv1alpha1 "go.linka.cloud/k8s/dns/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func AutoConfigSRVRecord(s *mailv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("autoconfig", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "autoconfig-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			SRV: &dnsv1alpha1.SRVRecord{
				Name:     dns.Fqdn("_autodiscover._tcp." + s.Spec.Domain),
				Ttl:      s.Spec.DNSTTL,
				Priority: 10,
				Weight:   10,
				Port:     443,
				Target:   dns.Fqdn("autodiscover." + s.Spec.Domain),
			},
		},
	}
}
