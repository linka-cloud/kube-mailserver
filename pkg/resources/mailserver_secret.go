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

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

func MailServerConfigSecret(s *mailv1alpha1.MailServer, bindDN, bindPW string) *corev1.Secret {
	config := MailServerConfig{
		OverrideHostname:              "mail." + s.Spec.Domain,
		OneDir:                        P(true),
		AccountProvisioner:            AccountProvisionerFile,
		PostmasterAddress:             "postmaster@" + s.Spec.Domain,
		SSLType:                       "manual",
		SpoofProtection:               V(s.Spec.Features.SpoofProtection),
		EnablePOP3:                    V(s.Spec.Features.POP3),
		EnableClamav:                  V(s.Spec.Features.Clamav),
		EnableAmavis:                  V(s.Spec.Features.Amavis),
		EnableFail2ban:                V(s.Spec.Features.Fail2ban),
		EnableManageSieve:             V(s.Spec.Features.ManageSieve),
		PostscreenAction:              "enforce",
		ClamavMessageSizeLimit:        "",
		VirusMailsDeleteDelay:         "",
		EnablePostfixVirtualTransport: false,
		PostfixDagent:                 "",
		PostfixMailboxSizeLimit:       "",
		EnableQuotas:                  V(s.Spec.Features.Quotas),
		PostfixMessageSizeLimit:       "",
		PflogsummTrigger:              "",
		PflogsummRecipient:            "",
		PflogsummSender:               "",
		LogwatchInterval:              "",
		LogwatchRecipient:             "",
		LogwatchSender:                "",
		ReportRecipient:               "",
		ReportSender:                  "",
		LogrotateInterval:             "",
		PostfixInetProtocols:          "",
		DovecotInetProtocols:          "",
		EnableSpamassassin:            V(s.Spec.Features.Spamassassin),
		SpamassassinSpamToInbox:       true,
		EnableSpamassassinKam:         V(s.Spec.Features.SpamassassinKam),
		MoveSpamToJunk:                true,
		SATag:                         "2.0",
		SATag2:                        "6.31",
		SAKill:                        "6.31",
		SASpamSubject:                 "***SPAM*****",

		EnableFetchmail: false,

		EnablePostgrey: V(s.Spec.Features.Postgrey),
		PostgreyDelay:  "300",
		PostgreyMaxAge: "35",
		PostgreyText:   "Delayed by postgrey",

		EnableSRS:         false,
		SRSSenderClasses:  "",
		SRSExcludeDomains: "",
		SRSSecret:         "",

		DefaultRelayHost: "",
		RelayHost:        "",
		RelayPort:        0,
		RelayUser:        "",
		RelayPassword:    "",

		EnableLdap:            s.Spec.Features.LDAP.Enabled,
		LdapStartTLS:          s.Spec.Features.LDAP.StartTLS,
		LdapServerHost:        fmt.Sprintf("ldaps://%s", s.Spec.Features.LDAP.Host),
		LdapSearchBase:        s.Spec.Features.LDAP.SearchBase,
		LdapBindDN:            bindDN,
		LdapBindPW:            bindPW,
		LdapQueryFilterUser:   "(mail=%s)",
		LdapQueryFilterGroup:  "(&(objectclass=group)(mail=%s))",
		LdapQueryFilterAlias:  "(&(objectClass=user)(otherMailbox=%s))",
		LdapQueryFilterDomain: "(mail=*@%s)",

		DovecotTLS:           true,
		DovecotHosts:         s.Spec.Features.LDAP.Host,
		DovecotLdapVersion:   "3",
		DovecotUserFilter:    "(mail=%u)",
		DovecotPassFilter:    "(mail=%u)",
		DovecotMailboxFormat: "",
		DovecotAuthBind:      true,
		DovecotScope:         "subtree",
		DovecotUserAttrs:     "=uid=5000,=gid=5000,=user=%{ldap:mail},=mail=maildir:/var/mail/%d/%n/,=home=/var/mail/%d/%n/,",

		EnableSASLAuthd:            false,
		SASLAuthdMechanisms:        "ldap",
		SASLAuthdMechOptions:       "",
		SASLAuthdLdapServer:        fmt.Sprintf("ldaps://%s", s.Spec.Features.LDAP.Host),
		SASLAuthdLdapBindDn:        bindDN,
		SASLAuthdLdapPassword:      bindPW,
		SASLAuthdLdapSearchBase:    s.Spec.Features.LDAP.SearchBase,
		SASLAuthdLdapFilter:        "(mail=%s)",
		SASLAuthdLdapStartTls:      s.Spec.Features.LDAP.StartTLS,
		SASLAuthdLdapTlsCheckPeer:  false,
		SASLAuthdLdapTlsCacertFile: "",
		SASLAuthdLdapTlsCacertDir:  "",
		SASLAuthdLdapPasswordAttr:  "",
		SASLPasswd:                 "",
		SASLAuthdLdapAuthMethod:    "",
		SASLAuthdLdapMech:          "",
	}
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize("config", s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, "config"),
		},
		Data: config.ToMap(),
	}
}
