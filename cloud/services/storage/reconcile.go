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
	"errors"
	"fmt"

	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/storage/v1"
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

	// Reconcile Bucket
	// check if bucket exists
	bucket, err := s.Buckets.Get(ctx, &meta.Key{Name: s.scope.Bucket().Name})
	if err != nil {
		if isNotFoundError(err) {
			log.V(2).Info("Bucket not found, creating bucket")

			err := s.Buckets.Insert(ctx, &meta.Key{Name: s.scope.Bucket().Name}, &storage.Bucket{
				Name:     s.scope.Bucket().Name,
				Location: s.scope.Region(),
			})
			if err != nil {
				return fmt.Errorf("failed to create bucket: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get cloud storage bucket: %w", err)
		}

	}
	log.V(2).Info("Bucket found, skipping create: ", bucket.Name)

	return nil
}

func isNotFoundError(err error) bool {
	var e *googleapi.Error
	if ok := errors.As(err, &e); ok {
		// Check if the error is a "not found" error
		if e.Code == 404 {
			return true
		}
	}
	return false
}
