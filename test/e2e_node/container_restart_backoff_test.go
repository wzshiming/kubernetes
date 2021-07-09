/*
Copyright 2021 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2enode

import (
	"context"

	"github.com/onsi/ginkgo"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/kubernetes/test/e2e/framework"
)

//var _ = SIGDescribe("ContainerRestartBackoff [Slow]", func() {
var _ = SIGDescribe("ContainerRestartBackoff", func() {
	f := framework.NewDefaultFramework("container-restart-backoff-test")
	ginkgo.Context("when a container restart", func() {

		ginkgo.It("should be restarted the container with an exponential backoff", func() {
			ginkgo.By("create container")
			name := "test-container-restart-backoff"
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: name,
				},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyAlways,
					Containers: []v1.Container{
						{
							Name:  "c",
							Image: busyboxImage,
							Command: []string{
								"sleep",
								"600",
							},
							StartupProbe: &v1.Probe{
								Handler: v1.Handler{
									Exec: &v1.ExecAction{
										Command: []string{
											"false",
										},
									},
								},
								PeriodSeconds:    1,
								FailureThreshold: 1,
							},
						},
					},
				},
			}
			pod = f.PodClient().CreateSync(pod)
			watcher, err := f.PodClient().Watch(context.Background(), metav1.ListOptions{
				FieldSelector: fields.OneTermEqualSelector("metadata.name", pod.Name).String(),
			})
			framework.ExpectNoError(err, "watch pod")
			for result := range watcher.ResultChan() {
				pod := result.Object.(*v1.Pod)
				status := pod.Status.ContainerStatuses[0]
				if terminated := status.LastTerminationState.Terminated; terminated == nil {
					framework.Logf("========DEBUG %s %d", result.Type, status.RestartCount)
				} else {
					framework.Logf("========DEBUG %s %d %s %s", result.Type, status.RestartCount, terminated.StartedAt, terminated.FinishedAt)
				}
			}
			framework.ExpectError(nil, "DEBUG")
		})
	})
})
