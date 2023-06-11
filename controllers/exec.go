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

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"go.linka.cloud/k8s"
	appsv1 "go.linka.cloud/k8s/apps/v1"
	corev1 "go.linka.cloud/k8s/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *MailServerReconciler) execDeployOut(ctx context.Context, deploy *appsv1.Deployment, command string) (string, bool, error) {
	log := ctrl.LoggerFrom(ctx)
	var pods corev1.PodList
	if err := r.List(ctx, &pods, client.InNamespace(deploy.Namespace), client.MatchingLabels(deploy.Spec.Template.Labels)); err != nil {
		return "", false, err
	}
	if len(pods.Items) == 0 {
		log.V(5).Info("mail server pod not ready yet")
		return "", false, nil
	}

	if k8s.Value(pods.Items[0].Status.Phase) != corev1.PodRunning {
		log.V(5).Info("mail server pod not running yet")
		return "", false, nil
	}
	out, err := r.execOut(client.ObjectKeyFromObject(&pods.Items[0]), command)
	return out, true, err
}

func (r *MailServerReconciler) execOut(key client.ObjectKey, command string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	if err := r.exec(key, command, nil, &stdout, &stderr); err != nil {
		return "", fmt.Errorf("%v: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

// exec execute command on specific pod and wait the command's output.
func (r *MailServerReconciler) exec(key client.ObjectKey, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := []string{
		"sh",
		"-c",
		command,
	}

	req := r.GoClient.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Namespace(key.Namespace).
		Name(key.Name).
		SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}
	if stdin == nil {
		option.Stdin = false
	}
	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(r.RestConfig, http.MethodPost, req.URL())
	if err != nil {
		return err
	}
	if err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}); err != nil {
		return err
	}
	return nil
}
