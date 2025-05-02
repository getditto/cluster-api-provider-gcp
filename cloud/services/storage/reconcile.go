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

// Package storage implements reconciler for cloud storage.
package storage

import (
	"context"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/storage/v1"
	"sigs.k8s.io/cluster-api-provider-gcp/feature"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Delete implements cloud.Reconciler.
func (s *Service) Delete(ctx context.Context) error {
	if !feature.Gates.Enabled(feature.WorkloadIDFederation) {
		return nil
	}

	return nil
}

// Reconcile implements cloud.Reconciler.
func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	if !feature.Gates.Enabled(feature.WorkloadIDFederation) {
		log.V(4).Info("WorkloadIDFederation feature gate is disabled, skipping reconcile")
		return nil
	}
	log.V(2).Info("Reconciling storage")

	// reconcile Bucket
	s.buckets.Insert(ctx, &meta.Key{Name: s.scope.Name()}, &storage.Bucket{
		Name: s.scope.Name(),
	})

	return nil
}
