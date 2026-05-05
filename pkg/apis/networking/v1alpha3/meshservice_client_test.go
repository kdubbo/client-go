//
// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha3_test

import (
	"context"
	"testing"

	apinetworking "github.com/kdubbo/api/networking/v1alpha3"
	networking "github.com/kdubbo/client-go/pkg/apis/networking/v1alpha3"
	"github.com/kdubbo/client-go/pkg/clientset/versioned/fake"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMeshServiceClientCRUD(t *testing.T) {
	ctx := context.Background()
	client := fake.NewSimpleClientset()
	meshServices := client.NetworkingV1alpha3().MeshServices("app")

	created, err := meshServices.Create(ctx, &networking.MeshService{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.dubbo.apache.org/v1alpha3",
			Kind:       "MeshService",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "nginx-routing",
			Namespace: "app",
			Labels:    map[string]string{"app": "nginx"},
		},
		Spec: apinetworking.MeshService{
			Hosts: []string{"nginx.app.svc.cluster.local"},
			Routes: []*apinetworking.MeshServiceRoute{
				{
					Service: []*apinetworking.ServiceDestination{
						{Name: "v1", Host: "nginx.app.svc.cluster.local", Weight: 20},
						{Name: "v2", Host: "nginx.app.svc.cluster.local", Weight: 80},
					},
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	if created.Spec.GetRoutes()[0].GetService()[1].GetWeight() != 80 {
		t.Fatalf("created route weight = %d, want 80", created.Spec.GetRoutes()[0].GetService()[1].GetWeight())
	}

	got, err := meshServices.Get(ctx, "nginx-routing", metav1.GetOptions{})
	if err != nil {
		t.Fatal(err)
	}
	got.Spec.Routes[0].Service[0].Weight = 50
	if _, err := meshServices.Update(ctx, got, metav1.UpdateOptions{}); err != nil {
		t.Fatal(err)
	}

	list, err := meshServices.List(ctx, metav1.ListOptions{LabelSelector: "app=nginx"})
	if err != nil {
		t.Fatal(err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("MeshService list length = %d, want 1", len(list.Items))
	}
}
