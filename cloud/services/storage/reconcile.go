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
	"fmt"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/storage/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/cluster-api-provider-gcp/feature"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Delete implements cloud.Reconciler.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	if !feature.Gates.Enabled(feature.WorkloadIDFederation) {
		log.V(4).Info("WorkloadIDFederation feature gate is disabled, skipping reconcile")
		return nil
	}

	log.V(2).Info("Deleting cloud storage")

	if s.scope.Bucket() == nil || s.scope.Bucket().Name == "" {
		log.V(4).Info("bucket not configured, skipping delete")
		return nil
	}

	log = log.WithValues("bucket", s.scope.Bucket().Name)

	// Get bucket and check it is owned by this cluster
	bucket, err := s.Buckets.Get(ctx, &meta.Key{Name: s.scope.Bucket().Name})
	if err != nil {
		if gcperrors.IsNotFound(err) {
			log.V(2).Info("Bucket not found, skipping delete")
			return nil
		}
		return fmt.Errorf("failed to get cloud storage bucket: %w", err)
	}
	if bucket.Labels == nil {
		log.V(2).Info("Bucket labels not found, skipping delete")
		return nil
	}
	if lbl, ok := bucket.Labels[infrav1.ClusterTagKey(s.scope.Name())]; !ok {
		log.V(2).Info("Bucket labels not found, skipping delete")
		return nil
	} else if lbl != "owned" {
		log.V(2).Info("Bucket labels not owned by this cluster, skipping delete")
		return nil
	}

	// Delete Bucket
	if err := s.Buckets.Delete(ctx, &meta.Key{Name: s.scope.Bucket().Name}); err != nil {
		log.Error(err, "Failed to delete cloud storage bucket")
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
	log.Info("Reconciling cloud storage resources")

	if s.scope.Bucket() == nil || s.scope.Bucket().Name == "" {
		log.V(4).Info("bucket not configured, skipping reconcile")
		return nil
	}
	log = log.WithValues("bucket", s.scope.Bucket().Name)

	// Reconcile Bucket
	// check if bucket exists
	_, err := s.Buckets.Get(ctx, &meta.Key{Name: s.scope.Bucket().Name})
	if err != nil {
		if gcperrors.IsNotFound(err) {
			log.V(2).Info("Bucket not found, creating bucket")

			err := s.Buckets.Insert(ctx, &meta.Key{Name: s.scope.Bucket().Name}, &storage.Bucket{
				Name:     s.scope.Bucket().Name,
				Location: s.scope.Region(),
				Labels: map[string]string{
					infrav1.ClusterTagKey(s.scope.Name()): "owned",
				},
			})
			if err != nil {
				return fmt.Errorf("failed to create bucket: %w", err)
			}

		} else {
			return fmt.Errorf("failed to get cloud storage bucket: %w", err)
		}

	} else {
		log.V(2).Info("Bucket found, skipping create")
	}

	return nil
}
