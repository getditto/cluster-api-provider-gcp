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
	"fmt"
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
	ProjectNumber() (int64, error)
	WorkloadIdentityPoolSpec() *iam.WorkloadIdentityPool
	WorkloadIdentityProviderSpec() *iam.WorkloadIdentityPoolProvider
}

// Service implements workloadidentity reconciler.
type Service struct {
	scope         Scope
	projectNumber int64
	widp          *iam.ProjectsLocationsWorkloadIdentityPoolsService
	widpp         *iam.ProjectsLocationsWorkloadIdentityPoolsProvidersService
}

var _ cloud.Reconciler = &Service{}

func New(scope Scope) *Service {
	if scope.IAMService() == nil {
		log.Fatalln("IAMService is nil")
	}

	projectNumber, err := scope.ProjectNumber()
	if err != nil {
		log.Fatalf("failed to get project number: %v", err)
	}

	return &Service{
		scope:         scope,
		projectNumber: projectNumber,
		widp:          iam.NewProjectsLocationsWorkloadIdentityPoolsService(scope.IAMService()),
		widpp:         iam.NewProjectsLocationsWorkloadIdentityPoolsProvidersService(scope.IAMService()),
	}
}

// parent returns the parent path for the workload identity pool and provider.
// The parent path is in the format:
// projects/{project}/locations/global
// where {project} is the project ID.
func (s *Service) parent() string {
	return "projects/" + s.scope.Project() + "/locations/" + Location
}

// parentByProjectNumber returns the parent path for the workload identity pool and provider.
// The parent path is in the format:
// projects/{projectNumber}/locations/global
// where {projectNumber} is the project number.
func (s *Service) parentByProjectNumber() string {
	return fmt.Sprintf("projects/%d/locations/%s", s.projectNumber, Location)
}
