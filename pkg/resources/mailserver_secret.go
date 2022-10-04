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
		SpoofProtection:               true,
		EnablePOP3:                    true,
		EnableClamav:                  false,
		EnableAmavis:                  true,
		EnableFail2ban:                true,
		EnableManageSieve:             true,
		PostscreenAction:              "enforce",
		ClamavMessageSizeLimit:        "",
		VirusMailsDeleteDelay:         "",
		EnablePostfixVirtualTransport: false,
		PostfixDagent:                 "",
		PostfixMailboxSizeLimit:       "",
		EnableQuotas:                  false,
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
		EnableSpamassassin:            true,
		SpamassassinSpamToInbox:       "",
		EnableSpamassassinKam:         false,
		MoveSpamToJunk:                false,
		SATag:                         "2.0",
		SATag2:                        "6.31",
		SAKill:                        "6.31",
		SASpamSubject:                 "***SPAM*****",

		EnableFetchmail: false,

		EnablePostgrey: false,
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

		EnableLdap:            s.Spec.LDAP.Enabled,
		LdapStartTLS:          s.Spec.LDAP.StartTLS,
		LdapServerHost:        fmt.Sprintf("ldaps://%s", s.Spec.LDAP.Host),
		LdapSearchBase:        s.Spec.LDAP.SearchBase,
		LdapBindDN:            bindDN,
		LdapBindPW:            bindPW,
		LdapQueryFilterUser:   "(mail=%s)",
		LdapQueryFilterGroup:  "(&(objectclass=group)(mail=%s))",
		LdapQueryFilterAlias:  "(&(objectClass=user)(otherMailbox=%s))",
		LdapQueryFilterDomain: "(mail=*@%s)",

		DovecotTLS:           true,
		DovecotHosts:         s.Spec.LDAP.Host,
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
		SASLAuthdLdapServer:        fmt.Sprintf("ldaps://%s", s.Spec.LDAP.Host),
		SASLAuthdLdapBindDn:        bindDN,
		SASLAuthdLdapPassword:      bindPW,
		SASLAuthdLdapSearchBase:    s.Spec.LDAP.SearchBase,
		SASLAuthdLdapFilter:        "(mail=%s)",
		SASLAuthdLdapStartTls:      s.Spec.LDAP.StartTLS,
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
