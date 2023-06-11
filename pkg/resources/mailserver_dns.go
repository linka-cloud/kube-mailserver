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
	"bytes"
	"fmt"
	"text/template"

	"github.com/miekg/dns"
	dnsv1alpha1 "go.linka.cloud/k8s/dns/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerARecord(s *mv1alpha1.MailServer, ip string) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("mail", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "mail-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			A: &dnsv1alpha1.ARecord{
				Name:   dns.Fqdn("mail." + s.Spec.Domain),
				Ttl:    s.Spec.DNSTTL,
				Target: ip,
			},
		},
	}
}

func MailServerMXRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("mx", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "mx-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			MX: &dnsv1alpha1.MXRecord{
				Name:       dns.Fqdn(s.Spec.Domain),
				Ttl:        s.Spec.DNSTTL,
				Preference: 10,
				Target:     dns.Fqdn("mail." + s.Spec.Domain),
			},
		},
	}
}

func MailServerDMARCRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	dmarc := fmt.Sprintf("v=DMARC1; p=reject; rua=mailto:dmarc@%[1]s; ruf=mailto:dmarc@%[1]s; fo=0; adkim=r; aspf=r; pct=100; rf=afrf; ri=86400; sp=quarantine", s.Spec.Domain)
	if s.Spec.DMARC != "" {
		dmarc = s.Spec.DMARC
	}
	dmarc = parseDMARC(s)
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("dmarc", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "dmarc-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			TXT: &dnsv1alpha1.TXTRecord{
				Name:    dns.Fqdn("_dmarc." + s.Spec.Domain),
				Ttl:     s.Spec.DNSTTL,
				Targets: []string{dmarc},
			},
		},
	}
}

func MailServerDKIMRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("dkim", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "dkim-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			// TODO(adphi): will be parsed from the domain.key file
			TXT: &dnsv1alpha1.TXTRecord{
				Name:    dns.Fqdn("mail._domainkey." + s.Spec.Domain),
				Ttl:     s.Spec.DNSTTL,
				Targets: []string{"v=DKIM1; k=rsa;"},
			},
		},
	}
}

func MailServerSPFRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("spf", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "spf-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			TXT: &dnsv1alpha1.TXTRecord{
				Name:    dns.Fqdn(s.Spec.Domain),
				Ttl:     s.Spec.DNSTTL,
				Targets: []string{s.Spec.SPF},
			},
		},
	}
}

func MailServerIMAPRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("imap", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "imap-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			SRV: &dnsv1alpha1.SRVRecord{
				Name:     dns.Fqdn("_imap._tcp." + s.Spec.Domain),
				Ttl:      s.Spec.DNSTTL,
				Priority: 10,
				Weight:   10,
				Port:     143,
				Target:   dns.Fqdn("mail." + s.Spec.Domain),
			},
		},
	}
}

func MailServerIMAPsRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("imaps", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "imaps-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			SRV: &dnsv1alpha1.SRVRecord{
				Name:     dns.Fqdn("_imaps._tcp." + s.Spec.Domain),
				Ttl:      s.Spec.DNSTTL,
				Priority: 10,
				Weight:   10,
				Port:     993,
				Target:   dns.Fqdn("mail." + s.Spec.Domain),
			},
		},
	}
}

func MailServerSubmissionRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("submission", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "submission-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			SRV: &dnsv1alpha1.SRVRecord{
				Name:     dns.Fqdn("_submission._tcp." + s.Spec.Domain),
				Ttl:      s.Spec.DNSTTL,
				Priority: 10,
				Weight:   10,
				Port:     587,
				Target:   dns.Fqdn("mail." + s.Spec.Domain),
			},
		},
	}
}

func MailServerPOP3Record(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("pop3", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "pop3-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			SRV: &dnsv1alpha1.SRVRecord{
				Name:     dns.Fqdn("_pop3._tcp." + s.Spec.Domain),
				Ttl:      s.Spec.DNSTTL,
				Priority: 10,
				Weight:   10,
				Port:     110,
				Target:   dns.Fqdn("mail." + s.Spec.Domain),
			},
		},
	}
}

func MailServerPOP3sRecord(s *mv1alpha1.MailServer) *dnsv1alpha1.DNSRecord {
	return &dnsv1alpha1.DNSRecord{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "dns.linka.cloud/v1alpha1",
			Kind:       "DNSRecord",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("pop3s", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "pop3s-record"),
		},
		Spec: dnsv1alpha1.DNSRecordSpec{
			SRV: &dnsv1alpha1.SRVRecord{
				Name:     dns.Fqdn("_pop3s._tcp." + s.Spec.Domain),
				Ttl:      s.Spec.DNSTTL,
				Priority: 10,
				Weight:   10,
				Port:     995,
				Target:   dns.Fqdn("mail." + s.Spec.Domain),
			},
		},
	}
}

func parseDMARC(s *mv1alpha1.MailServer) string {
	t, err := template.New("dmarc").Parse(s.Spec.DMARC)
	if err != nil {
		return ""
	}
	var buff bytes.Buffer
	if err := t.Execute(&buff, s.Spec); err != nil {
		return ""
	}
	return buff.String()
}
