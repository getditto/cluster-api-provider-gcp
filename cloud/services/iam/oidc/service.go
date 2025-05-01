package oidc

import (
	"log"

	gcstorage "google.golang.org/api/storage/v1"
	"sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/storage"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Scope is an interfaces that hold used methods.
type Scope interface {
	cloud.Cluster
	Bucket() *v1beta1.Bucket
	IssuerUri() string
	ManagementClient() client.Client
	StorageService() *gcstorage.Service
}

// Service implements workloadidentity reconciler.
type Service struct {
	scope   Scope
	Objects storage.ObjectsInterface
}

var _ cloud.Reconciler = &Service{}

// New returns Service from given scope.
func New(scope Scope) *Service {
	if scope.StorageService() == nil {
		log.Fatalln("StorageService is nil")
	}

	return &Service{
		scope:   scope,
		Objects: storage.NewObjectsService(scope.StorageService()),
	}
}
