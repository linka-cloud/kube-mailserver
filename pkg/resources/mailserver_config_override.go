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
	corev1 "go.linka.cloud/k8s/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	mv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

const (
	ConfigOverride = "config-override"
	PostfixGroups  = "ldap-groups.cf"
	PostfixMain    = "postfix-main.cf"
	PostfixMaster  = "postfix-master.cf"
	DovecotConf    = "dovecot.cf"
)

type ConfigMapEntry struct {
	Name  string
	Value string
}

var (
	PostfixADGroups = ConfigMapEntry{
		Name: PostfixGroups,
		Value: `bind                     = yes
bind_dn                  = cn=admin,dc=domain,dc=com
bind_pw                  = admin
query_filter             = (&(mailGroupMember=%s)(mailEnabled=TRUE))
search_base              = ou=people,dc=domain,dc=com
server_host              = mail.domain.com
start_tls                = no
version                  = 3

leaf_result_attribute = mail
special_result_attribute = member
`,
	}
)

// https://docker-mailserver.github.io/docker-mailserver/edge/config/advanced/kubernetes/#configure-the-mailserver
var (
	PostfixProxyMain = ConfigMapEntry{
		Name:  PostfixMain,
		Value: `postscreen_upstream_proxy_protocol = haproxy`,
	}
	PostfixProxyMaster = ConfigMapEntry{
		Name: PostfixMaster,
		Value: `smtp/inet/postscreen_upstream_proxy_protocol=haproxy
submission/inet/smtpd_upstream_proxy_protocol=haproxy
smtps/inet/smtpd_upstream_proxy_protocol=haproxy
`,
	}
	DovecotProxy = ConfigMapEntry{
		Name: DovecotConf,
		Value: `
	haproxy_trusted_networks = 127.0.0.0/8 #TODO: Replace with LoadBalancer IP
    service imap-login {
      inet_listener imap {
        haproxy = yes
      }
      inet_listener imaps {
        haproxy = yes
      }
    }
`,
	}
)

func ConfigOverrideConfigMap(s *mv1alpha1.MailServer) *corev1.ConfigMap {
	data := make(map[string]string)
	if s.Spec.Features.LDAP.Enabled {
		data[PostfixADGroups.Name] = PostfixADGroups.Value
	}
	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      Normalize(ConfigOverride, s.Spec.Domain),
			Namespace: s.Namespace,
			Labels:    Labels(s, ConfigOverride),
		},
		Data: data,
	}
}
