/*
Copyright 2025 The Kubernetes Authors.

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

package workloadidentity

import (
	"log"

	"google.golang.org/api/iam/v1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
)

const (
	Location = "global"
)

// Scope is an interfaces that hold used methods.
type Scope interface {
	cloud.Cluster
	IAMService() *iam.Service
	Project() string
	WorkloadIdentityPoolSpec() *iam.WorkloadIdentityPool
	WorkloadIdentityProviderSpec() *iam.WorkloadIdentityPoolProvider
}

// Service implements workloadidentity reconciler.
type Service struct {
	scope Scope
	// workloadIdentityPools workloadIdentityPoolsInterface

	widp  *iam.ProjectsLocationsWorkloadIdentityPoolsService
	widpp *iam.ProjectsLocationsWorkloadIdentityPoolsProvidersService
}

var _ cloud.Reconciler = &Service{}

func New(scope Scope) *Service {
	if scope.IAMService() == nil {
		log.Fatalln("IAMService is nil")
	}

	return &Service{
		scope: scope,
		widp:  iam.NewProjectsLocationsWorkloadIdentityPoolsService(scope.IAMService()),
		widpp: iam.NewProjectsLocationsWorkloadIdentityPoolsProvidersService(scope.IAMService()),
	}
}

func (s *Service) parent() string {
	return "projects/" + s.scope.Project() + "/locations/" + Location
}

// func (s *Service) buildPoolId(name string) string {
// 	return s.scope.Project() + "/" + Location + "/" + name
// }
