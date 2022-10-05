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
	"strings"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
)

const (
	KubeMail   = "kube-mailserver"
	Postmaster = "postmaster"

	LabelName      = "app.kubernetes.io/name"
	LabelInstance  = "app.kubernetes.io/instance"
	LabelComponent = "app.kubernetes.io/component"
	LabelPartOf    = "app.kubernetes.io/part-of"
	LabelManagedBy = "app.kubernetes.io/managed-by"
)

func Labels(s *mailv1alpha1.MailServer, component string) map[string]string {
	// TODO(adphi): improve labels to match k8s best practices
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	return map[string]string{
		LabelName:      s.Spec.Domain,
		LabelInstance:  s.Spec.Domain,
		LabelComponent: component,
		LabelPartOf:    s.Spec.Domain,
		LabelManagedBy: KubeMail,
	}
}

func Normalize(ss ...string) string {
	return strings.Replace(strings.Join(ss, "."), ".", "-", -1)
}

func P[T any](v T) *T {
	return &v
}

type PT[T any] interface {
	*T
}

func V[T any, P PT[T]](p P, defaults ...T) (v T) {
	if p != nil {
		return *p
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return v
}
