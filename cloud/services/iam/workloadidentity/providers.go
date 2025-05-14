package workloadidentity

import (
	"context"
	"fmt"

	"google.golang.org/api/iam/v1"
)

type workloadIdentityPoolProvidersInterface interface {
	GetProvider(ctx context.Context, id string) (*iam.WorkloadIdentityPoolProvider, error)
	InsertProvider(ctx context.Context, spec *iam.WorkloadIdentityPoolProvider) error
	DeleteProvider(ctx context.Context, id string) error
}

func (s *Service) InsertProvider(_ context.Context, spec *iam.WorkloadIdentityPoolProvider) error {
	poolId := fmt.Sprintf("%s/workloadIdentityPools/%s", s.parentByProjectNumber(), s.scope.WorkloadIdentityPoolSpec().Name)
	_, err := s.widpp.Create(poolId, spec).WorkloadIdentityPoolProviderId(spec.Name).Do()
	if err != nil {
		return fmt.Errorf("failed to create workload identity pool provider: %w", err)
	}

	return err
}

// GetProvider implements workloadIdentityPoolProvidersInterface.
func (s *Service) GetProvider(ctx context.Context, id string) (*iam.WorkloadIdentityPoolProvider, error) {
	provider, err := s.widpp.Get(fmt.Sprintf("%s/workloadIdentityPools/%s/providers/%s", s.parentByProjectNumber(), s.scope.WorkloadIdentityPoolSpec().Name, id)).Do()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// DeleteProvider implements workloadIdentityPoolProvidersInterface.
func (s *Service) DeleteProvider(ctx context.Context, id string) error {
	_, err := s.widpp.Delete(fmt.Sprintf("%s/workloadIdentityPools/%s/providers/%s", s.parentByProjectNumber(), s.scope.WorkloadIdentityPoolSpec().Name, id)).Do()
	if err != nil {
		return fmt.Errorf("failed to delete workload identity pool provider: %w", err)
	}

	return nil
}

var _ workloadIdentityPoolProvidersInterface = &Service{}
