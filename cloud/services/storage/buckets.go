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
)

type bucketsInterface interface {
	Get(ctx context.Context, key *meta.Key) (*storage.Bucket, error)
	Insert(ctx context.Context, key *meta.Key, obj *storage.Bucket) error
	Delete(ctx context.Context, key *meta.Key) error
}

type Buckets struct {
	scope Scope
	svc   *storage.BucketsService
}

func NewBucketsService(scope Scope) *Buckets { // Get the storage service client
	return &Buckets{
		scope: scope,
		svc:   storage.NewBucketsService(scope.StorageService()),
	}
}

// Get implements bucketsInterface.
func (b *Buckets) Get(ctx context.Context, key *meta.Key) (*storage.Bucket, error) {
	// Use the client to get the bucket
	return b.svc.Get(key.Name).Context(ctx).Do()
}

// Insert implements bucketsInterface.
func (b *Buckets) Insert(ctx context.Context, key *meta.Key, obj *storage.Bucket) error {
	// Use the client to insert the bucket
	_, err := b.svc.Insert(b.scope.Project(), obj).Context(ctx).Do()
	return err
}

// Delete implements bucketsInterface.
func (b *Buckets) Delete(ctx context.Context, key *meta.Key) error {
	return b.svc.Delete(key.Name).Context(ctx).Do()
}

var _ bucketsInterface = &Buckets{}
