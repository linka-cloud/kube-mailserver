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

package kubectl_mailserver

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	client2 "sigs.k8s.io/controller-runtime/pkg/client"

	mailv1alpha1 "go.linka.cloud/kube-mailserver/api/v1alpha1"
	"go.linka.cloud/kube-mailserver/pkg/resources"
)

var (
	configFlags   = genericclioptions.NewConfigFlags(true)
	resourceFlags = genericclioptions.NewResourceBuilderFlags().WithAllNamespaces(true)
	client        client2.Client
	ns            string
	RootCmd       = cobra.Command{
		Use:           "mailserver [instance]",
		Short:         "mailserver setup command",
		Args:          cobra.MinimumNArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(cmd.Context())
			defer cancel()

			name := args[0]
			ns, _, _ = configFlags.ToRawKubeConfigLoader().Namespace()
			if ns == "" {
				ns = "default"
			}
			conf, err := configFlags.ToRESTConfig()
			if err != nil {
				return err
			}
			client, err = client2.New(conf, client2.Options{Scheme: scheme.Scheme})
			if err != nil {
				return err
			}
			var ms mailv1alpha1.MailServer
			if err := client.Get(ctx, client2.ObjectKey{Namespace: ns, Name: name}, &ms); err != nil {
				return err
			}
			var deploy appsv1.Deployment
			if err := client.Get(ctx, client2.ObjectKey{Namespace: ns, Name: resources.Normalize("mail", ms.Spec.Domain)}, &deploy); err != nil {
				return err
			}
			var pods corev1.PodList
			if err := client.List(ctx, &pods, client2.InNamespace(ns), client2.MatchingLabels(deploy.Spec.Template.Labels)); err != nil {
				return err
			}
			if len(pods.Items) == 0 {
				return fmt.Errorf("no pods found")
			}
			var p *corev1.Pod
			for _, v := range pods.Items {
				if v.Status.Phase == corev1.PodRunning {
					p = &v
					break
				}
			}
			if p == nil {
				return fmt.Errorf("no running pods found")
			}
			stderr, stdout := replacer{w: cmd.ErrOrStderr()}, replacer{w: cmd.OutOrStdout()}
			if err := Exec(ctx, conf, p, append([]string{"setup"}, args[1:]...), cmd.InOrStdin(), &stdout, &stderr); err != nil {
				return err
			}
			return nil
		},
	}
)

type replacer struct {
	w io.Writer
}

func (w *replacer) Write(p []byte) (n int, err error) {
	b := bytes.Replace(p, []byte("/usr/local/bin/setup"), []byte("setup"), -1)
	b = bytes.Replace(b, []byte("./setup.sh"), []byte("setup"), -1)
	b = bytes.Replace(b, []byte("./setup "), []byte("setup "), -1)
	_, err = w.w.Write(b)
	return len(p), err
}

func init() {
	utilruntime.Must(scheme.AddToScheme(scheme.Scheme))
	utilruntime.Must(mailv1alpha1.AddToScheme(scheme.Scheme))
	flags := pflag.NewFlagSet("kubectl-mailserver", pflag.ExitOnError)
	pflag.CommandLine = flags
	configFlags.AddFlags(RootCmd.Flags())
	resourceFlags.AddFlags(RootCmd.Flags())
}
