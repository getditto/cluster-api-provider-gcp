package oidc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/go-jose/go-jose/v3"
	gcstorage "google.golang.org/api/storage/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/storage"
	"sigs.k8s.io/cluster-api-provider-gcp/feature"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile implements cloud.Reconciler.
func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	if !feature.Gates.Enabled(feature.WorkloadIDFederation) {
		log.V(4).Info("WorkloadIDFederation feature gate is disabled, skipping reconcile")
		return nil
	}
	log.Info("Reconciling oidc resources")

	if s.scope.Bucket() == nil || s.scope.Bucket().Name == "" {
		log.V(4).Info("bucket not configured, skipping reconcile")
		return nil
	}

	if err := s.reconcileBucketContents(ctx); err != nil {
		return fmt.Errorf("failed to reconcile bucket contents: %w", err)
	}

	return nil
}

// Delete implements cloud.Reconciler.
func (s *Service) Delete(ctx context.Context) error {
	log := log.FromContext(ctx)
	if !feature.Gates.Enabled(feature.WorkloadIDFederation) {
		log.V(4).Info("WorkloadIDFederation feature gate is disabled, skipping reconcile")
		return nil
	}
	log.Info("Reconciling oidc resources")

	if s.scope.Bucket() == nil || s.scope.Bucket().Name == "" {
		log.V(4).Info("bucket not configured, skipping reconcile")
		return nil
	}

	// Delete the OIDC discovery document
	if err := s.Objects.Delete(ctx, s.scope.Bucket().Name, s.buildOidcDocObjectKey()); err != nil {
		log.Error(err, "Failed to delete OIDC discovery document")
	}
	// Delete the JWKS document
	if err := s.Objects.Delete(ctx, s.scope.Bucket().Name, s.buildJwksObjectKey()); err != nil {
		log.Error(err, "Failed to delete JWKS document")
	}
	return nil
}

func (s *Service) reconcileBucketContents(ctx context.Context) error {
	// create the OpenID Connect discovery document
	openIDConfig, err := buildDiscoveryJSON(s.buildIssuerURL())
	if err != nil {
		return err
	}

	storageSvc := storage.New(s.scope)

	if err := storageSvc.Objects.Insert(
		ctx,
		s.scope.Bucket().Name,
		&gcstorage.Object{
			Name:        s.buildOidcDocObjectKey(),
			ContentType: "application/json",
			Acl: []*gcstorage.ObjectAccessControl{
				{
					Entity: "allUsers",
					Role:   "READER",
				},
			},
		},
		bytes.NewReader(openIDConfig)); err != nil {
		return fmt.Errorf("failed to create OIDC discovery document object in GCS: %w", err)
	}

	// Read the <cluster>-sa secret that contains Service Account signing key
	secret := &corev1.Secret{}
	if err := s.scope.ManagementClient().Get(ctx, client.ObjectKey{
		Name:      s.scope.Name() + "-sa",
		Namespace: s.scope.Namespace(),
	}, secret); err != nil {
		return fmt.Errorf("failed to get service account signing secret: %w", err)
	}

	// Create jwks document that will be published to GCS
	key, err := createJwksKey(secret.Data["tls.key"])
	if err != nil {
		return fmt.Errorf("failed to create jwks key: %w", err)
	}
	jwksDocument := jwksDocument{Keys: []jose.JSONWebKey{*key}}
	jwksBytes, err := json.MarshalIndent(jwksDocument, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal jwks payload to json: %w", err)
	}

	if err := storageSvc.Objects.Insert(
		ctx,
		s.scope.Bucket().Name,
		&gcstorage.Object{
			Name:        s.buildJwksObjectKey(),
			ContentType: "application/json",
			Acl: []*gcstorage.ObjectAccessControl{
				{
					Entity: "allUsers",
					Role:   "READER",
				},
			},
		},
		bytes.NewReader(jwksBytes)); err != nil {
		return fmt.Errorf("failed to create jwks document object in GCS: %w", err)
	}

	return nil
}

func (s *Service) buildOidcDocObjectKey() string {
	oidcDiscoveryObjectKey := path.Join(s.scope.Name(), opendIDConfigurationKey)
	return oidcDiscoveryObjectKey
}

func (s *Service) buildJwksObjectKey() string {
	oidcDiscoveryObjectKey := path.Join(s.scope.Name(), jwksKey)
	return oidcDiscoveryObjectKey
}
