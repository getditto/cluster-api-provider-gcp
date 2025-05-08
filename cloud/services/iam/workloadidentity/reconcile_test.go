package workloadidentity

import (
	"context"

	"google.golang.org/api/compute/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/utils/ptr"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/scope"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func init() {
	_ = clusterv1.AddToScheme(scheme.Scheme)
	_ = infrav1.AddToScheme(scheme.Scheme)
}

func getBaseClusterScope() (*scope.ClusterScope, error) {
	fakec := fake.NewClientBuilder().
		WithScheme(scheme.Scheme).
		Build()

	fakeCluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cluster",
			Namespace: "default",
		},
		Spec: clusterv1.ClusterSpec{},
	}

	fakeGCPCluster := &infrav1.GCPCluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cluster",
			Namespace: "default",
		},
		Spec: infrav1.GCPClusterSpec{
			Project: "my-proj",
			Region:  "us-central1",
			WorkloadIdentityFederation: &infrav1.WorkloadIdentityFederation{
				Enabled:      ptr.To(true),
				ProviderName: "my-provider",
				PoolName:     "my-pool",
			},
		},
		Status: infrav1.GCPClusterStatus{
			FailureDomains: clusterv1.FailureDomains{
				"us-central1-a": clusterv1.FailureDomainSpec{ControlPlane: true},
			},
		},
	}
	clusterScope, err := scope.NewClusterScope(context.TODO(), scope.ClusterScopeParams{
		Client:     fakec,
		Cluster:    fakeCluster,
		GCPCluster: fakeGCPCluster,
		GCPServices: scope.GCPServices{
			Compute: &compute.Service{},
		},
	})
	if err != nil {
		return nil, err
	}

	return clusterScope, nil
}
