package workloadidentity

import (
	"context"
	"fmt"

	"google.golang.org/api/iam/v1"
)

type workloadIdentityPoolsInterface interface {
	GetPool(ctx context.Context, id string) (*iam.WorkloadIdentityPool, error)
	InsertPool(ctx context.Context, spec *iam.WorkloadIdentityPool) error
	DeletePool(ctx context.Context, id string) error
}

// createWorkloadIdentityPool creates a new workload identity pool
func (s *Service) InsertPool(_ context.Context, spec *iam.WorkloadIdentityPool) error {
	_, err := s.widp.Create(s.parent(), spec).WorkloadIdentityPoolId(spec.Name).Do()
	if err != nil {
		return fmt.Errorf("failed to create workload identity pool: %w", err)
	}

	return nil
}

// GetPool implements workloadIdentityPoolsInterface.
func (s *Service) GetPool(ctx context.Context, id string) (*iam.WorkloadIdentityPool, error) {
	pool, err := s.widp.Get(fmt.Sprintf("%s/workloadIdentityPools/%s", s.parent(), id)).Do()
	if err != nil {
		return nil, err
	}

	return pool, nil
}

// DeletePool implements workloadIdentityPoolsInterface.
func (s *Service) DeletePool(ctx context.Context, id string) error {
	_, err := s.widp.Delete(fmt.Sprintf("%s/workloadIdentityPools/%s", s.parent(), id)).Do()
	if err != nil {
		return fmt.Errorf("failed to delete workload identity pool: %w", err)
	}

	return nil
}

var _ workloadIdentityPoolsInterface = &Service{}
