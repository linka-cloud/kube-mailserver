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

package kube_mailserver

import (
	"context"
	"fmt"

	"k8s.io/client-go/discovery"
	ctrl "sigs.k8s.io/controller-runtime"
)

type GroupVersionKind struct {
	Group    string
	Version  string
	Kind     string
	Required bool
}

func CheckGroupVersionKind(ctx context.Context, c discovery.DiscoveryInterface, gvk GroupVersionKind) error {
	gv := fmt.Sprintf("%s/%s", gvk.Group, gvk.Version)
	res, err := c.ServerResourcesForGroupVersion(gv)
	if err != nil {
		return fmt.Errorf("failed to get server resources for group version %s: %w", gv, err)
	}
	for _, r := range res.APIResources {
		if r.Kind == gvk.Kind {
			return nil
		}
	}
	if gvk.Required {
		return fmt.Errorf("resource %s not found in group version %s", gvk.Kind, gv)
	}
	ctrl.LoggerFrom(ctx).Info("resource not found", "group", gvk.Group, "version", gvk.Version, "kind", gvk.Kind)
	return nil
}

func CheckGroupVersionKinds(ctx context.Context, c discovery.DiscoveryInterface, gvks ...GroupVersionKind) error {
	for _, gvk := range gvks {
		if err := CheckGroupVersionKind(ctx, c, gvk); err != nil {
			return err
		}
	}
	return nil
}
