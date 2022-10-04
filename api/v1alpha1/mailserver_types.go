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

package v1alpha1

import (
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IPv4 is used for validation of an IPv6 address.
// +kubebuilder:validation:Pattern="^((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))$"
type IPv4 string

// MailServerSpec defines the desired state of MailServer
type MailServerSpec struct {
	// Domain is the mail server domain name
	// +kubebuilder:validation:Required
	Domain string `json:"domain,omitempty"`
	// SPF is the optional SPF configuration
	// If empty, the SPF record will be set to "v=spf1 mx ip4:$PUBLIC_IP -all"
	// $PUBLIC_IP will be replaced by the public IP of the mail server using `curl ifconfig.me`
	// +optional
	SPF string `json:"spf,omitempty"`
	// DMARC is the optional DMARC configuration
	// +optional
	// +kubebuilder:default="v=DMARC1; p=reject; rua=mailto:postmaster@{{ .Domain }}; ruf=mailto:postmaster@{{ .Domain }}; fo=0; adkim=r; aspf=r; pct=100; rf=afrf; ri=86400; sp=quarantine"
	DMARC string `json:"dmarc,omitempty"`
	// Image is the docker-mailserver image to use
	// +kubebuilder:validation:Required
	// +kubebuilder:default="docker.io/mailserver/docker-mailserver:9.1.0"
	Image            string `json:"image,omitempty"`
	DeploymentConfig `json:",inline"`
	// AutoConfig is the autoconfig deployment configuration
	// +optional
	AutoConfig AutoConfigDeployment `json:"autoconfig,omitempty"`
	// LoadBalancerIP is the optional IP address to request for the load balancer
	// +optional
	LoadBalancerIP *IPv4 `json:"loadBalancerIP,omitempty"`
	// OverrideIP is the optional IP address to use for the domain A record
	// +optional
	OverrideIP *IPv4 `json:"overrideIP,omitempty"`
	// Replicas is the number of replicas of the mail server
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=1
	Replicas *int32 `json:"replicas,omitempty"`
	// IssuerRef is the reference to the Cert Manager issuer to use for the certificate
	// +kubebuilder:validation:Required
	IssuerRef cmmeta.ObjectReference `json:"issuerRef"`
	// LDAP is the optional LDAP configuration
	// It expects an Active Directory like server
	// LDAP is not supported yet
	// +optional
	LDAP LDAPConfig `json:"ldap,omitempty"`

	// TODO(adphi): add custom config mounts support for the MailServer Deployment
	// Volume is the optional volume configuration
	// +optional
	Volume VolumeConfig `json:"volume,omitempty"`

	// Traefik is the optional Traefik configuration
	// +optional
	Traefik *TraefikConfig `json:"traefik,omitempty"`
}

type AutoConfigDeployment struct {
	DeploymentConfig `json:",inline"`
	// Image is the github.com/linka-cloud/go-autoconfig image to use
	// +kubebuilder:validation:Required
	// +kubebuilder:default="docker.io/linkacloud/autoconfig:latest"
	Image string `json:"image,omitempty"`
}

type DeploymentConfig struct {
	// ServiceAccountName is the name of the service account to use for the deployment
	// +optional
	ServiceAccountName string `json:"serviceAccountName,omitempty"`
	// Annotations is the optional annotations to add to the deployment
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// Labels is the optional labels to add to the deployment
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
	// Affinity is the optional affinity configuration for the deployment
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`
	// Strategy is the deployment strategy to use to replace existing pods with new ones.
	// +optional
	Strategy appsv1.DeploymentStrategy `json:"strategy,omitempty"`
	// SecurityContext is the optional security context for the deployment
	// +optional
	SecurityContext *corev1.PodSecurityContext `json:"securityContext,omitempty"`
	// TopologySpreadConstraints is the optional topology spread constraints configuration for the deployment
	// +optional
	TopologySpreadConstraints []corev1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
	// Tolerations are the optional toleration configurations for the deployment
	// +optional
	Tolerations []corev1.Toleration `json:"toleration,omitempty"`
	// NodeSelector is the optional node selector configuration for the deployment
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Resources is the optional resource configuration for the deployment
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
}

type TraefikConfig struct {
	Entrypoints TraefikEntrypoints `json:"entrypoints,omitempty"`
}

type TraefikEntrypoints struct {
	// HTTP is the optional HTTP entrypoint configuration
	// +optional
	// +kubebuilder:default="web"
	HTTP string `json:"http,omitempty"`
	// HTTPS is the optional HTTPS entrypoint configuration
	// +optional
	// +kubebuilder:default="websecure"
	HTTPS string `json:"https,omitempty"`
}

type VolumeConfig struct {
	// StorageClass is the name of the storage class to use
	// +optional
	// +kubebuilder:default=""
	StorageClass string `json:"storageClass,omitempty"`
	// Size is the size of the volume, defaults to 1Gi
	// +optional
	// +kubebuilder:default="1Gi"
	Size string `json:"size,omitempty"`
	// AccessModes is the access modes to use, defaults to ReadWriteMany
	// +optional
	// +kubebuilder:default="ReadWriteMany"
	AccessMode corev1.PersistentVolumeAccessMode `json:"accessMode,omitempty"`
}

type LDAPConfig struct {
	// Enabled enables the LDAP configuration
	// +optional
	Enabled bool `json:"enabled,omitempty"`
	// Host is the LDAP server host without ldap:// or ldaps://
	// Only TLS and StartTLS are supported
	// +kubebuilder:validation:Required
	Host string `json:"host,omitempty"`
	// Port is the LDAP server port
	// +kubebuilder:validation:Required
	Port int `json:"port,omitempty"`
	// StartTLS makes the connection should use StartTLS
	StartTLS bool `json:"startTLS,omitempty"`
	// Nameserver is the DNS server to use for the LDAP connection
	// +optional
	Nameserver *IPv4 `json:"nameserver,omitempty"`

	// BindSecret is the name of the secret containing the bind DN and password
	// It expects the following keys:
	// - bindDN: the DN of the LDAP lookup account
	// - bindPW: the password of the LDAP lookup account
	// +kubebuilder:validation:Required
	BindSecret string `json:"bindSecret,omitempty"`

	// SearchBase is the LDAP base where to search for users
	// +kubebuilder:validation:Required
	SearchBase string `json:"searchBase,omitempty"`
	// SearchFilter is the LDAP filter to use to search for users
	// +kubebuilder:validation:Required
	UserFilter string `json:"userFilter,omitempty"`
}

// MailServerStatus defines the observed state of MailServer
type MailServerStatus struct {
	Domain         string `json:"domain,omitempty"`
	Replicas       int32  `json:"replicas,omitempty"`
	Selector       string `json:"selector"`
	VolumeSize     string `json:"volumeSize,omitempty"`
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.selector
// +kubebuilder:resource:path=mailservers,shortName=ms;mail
// +kubebuilder:printcolumn:name="Domain",type=string,JSONPath=`.status.domain`
// +kubebuilder:printcolumn:name="Capacity",type=string,JSONPath=`.status.volumeSize`
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:printcolumn:name="IP",type="string",priority=1,JSONPath=".status.loadBalancerIP"
// +kubebuilder:printcolumn:name="Image",type="string",priority=1,JSONPath=".spec.image"

// MailServer is the Schema for the mailservers API
type MailServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MailServerSpec   `json:"spec,omitempty"`
	Status MailServerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MailServerList contains a list of MailServer
type MailServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MailServer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MailServer{}, &MailServerList{})
}
