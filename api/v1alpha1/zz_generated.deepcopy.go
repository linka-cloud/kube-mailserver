//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoConfig) DeepCopyInto(out *AutoConfig) {
	*out = *in
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
	in.Deployment.DeepCopyInto(&out.Deployment)
	in.Ingress.DeepCopyInto(&out.Ingress)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoConfig.
func (in *AutoConfig) DeepCopy() *AutoConfig {
	if in == nil {
		return nil
	}
	out := new(AutoConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AutoConfigDeployment) DeepCopyInto(out *AutoConfigDeployment) {
	*out = *in
	in.DeploymentConfig.DeepCopyInto(&out.DeploymentConfig)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AutoConfigDeployment.
func (in *AutoConfigDeployment) DeepCopy() *AutoConfigDeployment {
	if in == nil {
		return nil
	}
	out := new(AutoConfigDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentConfig) DeepCopyInto(out *DeploymentConfig) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	in.Strategy.DeepCopyInto(&out.Strategy)
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.TopologySpreadConstraints != nil {
		in, out := &in.TopologySpreadConstraints, &out.TopologySpreadConstraints
		*out = make([]v1.TopologySpreadConstraint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Resources.DeepCopyInto(&out.Resources)
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentConfig.
func (in *DeploymentConfig) DeepCopy() *DeploymentConfig {
	if in == nil {
		return nil
	}
	out := new(DeploymentConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Features) DeepCopyInto(out *Features) {
	*out = *in
	if in.POP3 != nil {
		in, out := &in.POP3, &out.POP3
		*out = new(bool)
		**out = **in
	}
	if in.SpoofProtection != nil {
		in, out := &in.SpoofProtection, &out.SpoofProtection
		*out = new(bool)
		**out = **in
	}
	if in.Clamav != nil {
		in, out := &in.Clamav, &out.Clamav
		*out = new(bool)
		**out = **in
	}
	if in.Amavis != nil {
		in, out := &in.Amavis, &out.Amavis
		*out = new(bool)
		**out = **in
	}
	if in.Fail2ban != nil {
		in, out := &in.Fail2ban, &out.Fail2ban
		*out = new(bool)
		**out = **in
	}
	if in.ManageSieve != nil {
		in, out := &in.ManageSieve, &out.ManageSieve
		*out = new(bool)
		**out = **in
	}
	if in.Quotas != nil {
		in, out := &in.Quotas, &out.Quotas
		*out = new(bool)
		**out = **in
	}
	if in.Spamassassin != nil {
		in, out := &in.Spamassassin, &out.Spamassassin
		*out = new(bool)
		**out = **in
	}
	if in.SpamassassinKam != nil {
		in, out := &in.SpamassassinKam, &out.SpamassassinKam
		*out = new(bool)
		**out = **in
	}
	if in.Postgrey != nil {
		in, out := &in.Postgrey, &out.Postgrey
		*out = new(bool)
		**out = **in
	}
	in.LDAP.DeepCopyInto(&out.LDAP)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Features.
func (in *Features) DeepCopy() *Features {
	if in == nil {
		return nil
	}
	out := new(Features)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressConfig) DeepCopyInto(out *IngressConfig) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressConfig.
func (in *IngressConfig) DeepCopy() *IngressConfig {
	if in == nil {
		return nil
	}
	out := new(IngressConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LDAPConfig) DeepCopyInto(out *LDAPConfig) {
	*out = *in
	if in.Nameserver != nil {
		in, out := &in.Nameserver, &out.Nameserver
		*out = new(IPv4)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LDAPConfig.
func (in *LDAPConfig) DeepCopy() *LDAPConfig {
	if in == nil {
		return nil
	}
	out := new(LDAPConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MailServer) DeepCopyInto(out *MailServer) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MailServer.
func (in *MailServer) DeepCopy() *MailServer {
	if in == nil {
		return nil
	}
	out := new(MailServer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MailServer) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MailServerList) DeepCopyInto(out *MailServerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MailServer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MailServerList.
func (in *MailServerList) DeepCopy() *MailServerList {
	if in == nil {
		return nil
	}
	out := new(MailServerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MailServerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MailServerSpec) DeepCopyInto(out *MailServerSpec) {
	*out = *in
	in.DeploymentConfig.DeepCopyInto(&out.DeploymentConfig)
	in.AutoConfig.DeepCopyInto(&out.AutoConfig)
	if in.LoadBalancerIP != nil {
		in, out := &in.LoadBalancerIP, &out.LoadBalancerIP
		*out = new(IPv4)
		**out = **in
	}
	if in.OverrideIP != nil {
		in, out := &in.OverrideIP, &out.OverrideIP
		*out = new(IPv4)
		**out = **in
	}
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	out.IssuerRef = in.IssuerRef
	in.Features.DeepCopyInto(&out.Features)
	out.Volume = in.Volume
	if in.Traefik != nil {
		in, out := &in.Traefik, &out.Traefik
		*out = new(TraefikConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MailServerSpec.
func (in *MailServerSpec) DeepCopy() *MailServerSpec {
	if in == nil {
		return nil
	}
	out := new(MailServerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MailServerStatus) DeepCopyInto(out *MailServerStatus) {
	*out = *in
	if in.Traefik != nil {
		in, out := &in.Traefik, &out.Traefik
		*out = new(bool)
		**out = **in
	}
	if in.AutoConfig != nil {
		in, out := &in.AutoConfig, &out.AutoConfig
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MailServerStatus.
func (in *MailServerStatus) DeepCopy() *MailServerStatus {
	if in == nil {
		return nil
	}
	out := new(MailServerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TraefikConfig) DeepCopyInto(out *TraefikConfig) {
	*out = *in
	out.Entrypoints = in.Entrypoints
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TraefikConfig.
func (in *TraefikConfig) DeepCopy() *TraefikConfig {
	if in == nil {
		return nil
	}
	out := new(TraefikConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TraefikEntrypoints) DeepCopyInto(out *TraefikEntrypoints) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TraefikEntrypoints.
func (in *TraefikEntrypoints) DeepCopy() *TraefikEntrypoints {
	if in == nil {
		return nil
	}
	out := new(TraefikEntrypoints)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeConfig) DeepCopyInto(out *VolumeConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeConfig.
func (in *VolumeConfig) DeepCopy() *VolumeConfig {
	if in == nil {
		return nil
	}
	out := new(VolumeConfig)
	in.DeepCopyInto(out)
	return out
}
