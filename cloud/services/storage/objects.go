package storage

import (
	"context"
	"io"

	"google.golang.org/api/storage/v1"
)

type objectsInterface interface {
	Get(ctx context.Context, bucket, key string) (*storage.Object, error)
	Insert(ctx context.Context, bucket, key string, obj *storage.Object, buf io.Reader) error
	Delete(ctx context.Context, bucket, key string) error
}

type Objects struct {
	scope Scope
	svc   *storage.ObjectsService
}

func NewObjectsService(scope Scope) *Objects { // Get the storage service client
	return &Objects{
		scope: scope,
		svc:   storage.NewObjectsService(scope.StorageService()),
	}
}

// Get implements ObjectsInterface.
func (b *Objects) Get(ctx context.Context, bucket, key string) (*storage.Object, error) {
	// Use the client to get the Object
	return b.svc.Get(bucket, key).Context(ctx).Do()
}

// Insert implements ObjectsInterface.
func (b *Objects) Insert(ctx context.Context, bucket, key string, obj *storage.Object, buf io.Reader) error {
	// Use the client to insert the Object
	_, err := b.svc.Insert(bucket, obj).
		Context(ctx).
		Media(buf).
		Do()
	return err
}

// Delete implements ObjectsInterface.
func (b *Objects) Delete(ctx context.Context, bucket, key string) error {
	return b.svc.Delete(bucket, key).Context(ctx).Do()
}

var _ objectsInterface = &Objects{}
