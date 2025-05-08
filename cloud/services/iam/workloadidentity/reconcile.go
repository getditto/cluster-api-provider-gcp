package workloadidentity

import (
	"context"
	"fmt"

	"sigs.k8s.io/cluster-api-provider-gcp/cloud/gcperrors"
	"sigs.k8s.io/cluster-api-provider-gcp/feature"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// Reconcile implements cloud.Reconciler.
func (s *Service) Reconcile(ctx context.Context) error {
	log := log.FromContext(ctx)
	if !feature.Gates.Enabled(feature.WorkloadIDFederation) {
		log.V(4).Info("WorkloadIDFederation feature gate is disabled, skipping reconcile")
		return nil
	}
	log.Info("Reconciling workload identity resources")

	// Reconcile Workload Identity Pool
	err := s.reconcilePool(ctx)
	if err != nil {
		return err
	}

	// Reconcile Workload Identity Provider
	err = s.reconcileProvider(ctx)
	if err != nil {
		return err
	}

	return nil
}

// reconcilePool reconciles the workload identity pool.
func (s *Service) reconcilePool(ctx context.Context) error {
	log := log.FromContext(ctx)

	pool := s.scope.WorkloadIdentityPoolSpec()
	if pool == nil {
		log.V(4).Info("workloadIdentityFederation not enabled, skipping reconcile")
		return nil
	}

	_, err := s.GetPool(ctx, pool.Name)
	if err != nil {
		if gcperrors.IsNotFound(err) {
			log.V(2).Info("Pool not found, creating...")

			err := s.InsertPool(ctx, pool)
			if err != nil {
				return fmt.Errorf("failed to create WI pool: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get WI pool: %w", err)
		}
	} else {

		// TODO: check if it in DELETED state

		log.V(2).Info("WI Pool found, skipping create")
	}
	return nil
}

func (s *Service) reconcileProvider(ctx context.Context) error {
	log := log.FromContext(ctx)

	provider := s.scope.WorkloadIdentityProviderSpec()
	if provider == nil {
		log.V(4).Info("workloadIdentityFederation not enabled, skipping reconcile")
		return nil
	}

	_, err := s.widpp.Get(fmt.Sprintf("%s/workloadIdentityPools/%s/providers/%s", s.parent(), s.scope.WorkloadIdentityPoolSpec().Name, provider.Name)).Do()
	if err != nil {
		if gcperrors.IsNotFound(err) {
			log.V(2).Info("Provider not found, creating...", "providerName", provider.Name)

			_, err := s.widpp.Create(
				fmt.Sprintf("%s/workloadIdentityPools/%s", s.parent(), s.scope.WorkloadIdentityPoolSpec().Name),
				provider,
			).WorkloadIdentityPoolProviderId(provider.Name).Do()
			if err != nil {
				return fmt.Errorf("failed to create WI provider: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get WI provider: %w", err)
		}
	} else {
		log.V(2).Info("WI Provider found, skipping create")
	}
	return nil
}

// Delete implements cloud.Reconciler.
func (s *Service) Delete(ctx context.Context) error {
	return nil
}
