package rke2

import (
	"context"
	"fmt"
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/machinepools"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/rancher/rancher/tests/framework/pkg/wait"
	"github.com/rancher/rancher/tests/integration/pkg/defaults"
	provisioning "github.com/rancher/rancher/tests/v2/validation/provisioning"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CertRotationTestSuite struct {
	suite.Suite
	session            *session.Session
	client             *rancher.Client
	config             *rancher.Config
	clusterName        string
	namespace          string
	kubernetesVersions []string
	cnis               []string
	providers          []string
}

// func (p *CertRotationTestSuite) TearDownSuite() {
// 	p.session.Cleanup()
// }

func (r *CertRotationTestSuite) SetupSuite() {
	testSession := session.NewSession(r.T())
	r.session = testSession

	clustersConfig := new(provisioning.Config)
	config.LoadConfig(provisioning.ConfigurationFileKey, clustersConfig)

	r.kubernetesVersions = clustersConfig.KubernetesVersions
	r.cnis = clustersConfig.CNIs
	r.providers = clustersConfig.Providers

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client

	r.clusterName = r.client.RancherConfig.ClusterName
	r.namespace = r.client.RancherConfig.ClusterName
}

func (r *CertRotationTestSuite) ProvisionRKE2Cluster(provider Provider) {
	// time.Sleep(10000000000000)
	providerName := " Node Provider: " + provider.Name
	nodeRoles0 := []machinepools.NodeRoles{
		{
			ControlPlane: true,
			Etcd:         true,
			Worker:       true,
			Quantity:     1,
		},
	}

	tests := []struct {
		name      string
		nodeRoles []machinepools.NodeRoles
		client    *rancher.Client
	}{
		{"1 Node all roles Admin User", nodeRoles0, r.client},
	}

	var name string
	for _, tt := range tests {
		subSession := r.session.NewSession()
		defer subSession.Cleanup()

		client, err := tt.client.WithSession(subSession)
		require.NoError(r.T(), err)

		cloudCredential, err := provider.CloudCredFunc(client)

		require.NoError(r.T(), err)
		kubeVersion := "1.24.2-rancher1-1"
		name = tt.name + providerName + " Kubernetes version: " + kubeVersion
		cni := "calico"
		name += " cni: " + cni
		r.Run(name, func() {
			testSession := session.NewSession(r.T())
			defer testSession.Cleanup()

			testSessionClient, err := tt.client.WithSession(testSession)
			require.NoError(r.T(), err)

			clusterName := provisioning.AppendRandomString(provider.Name)
			fmt.Println(clusterName)
			generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
			machinePoolConfig := provider.MachinePoolFunc(generatedPoolName, namespace)

			machineConfigResp, err := machinepools.CreateMachineConfig(provider.MachineConfig, machinePoolConfig, testSessionClient)
			require.NoError(r.T(), err)

			machinePools := machinepools.RKEMachinePoolSetup(tt.nodeRoles, machineConfigResp)

			cluster := clusters.NewRKE2ClusterConfig(clusterName, namespace, cni, cloudCredential.ID, kubeVersion, machinePools)

			clusterResp, err := clusters.CreateRKE2Cluster(testSessionClient, cluster)
			require.NoError(r.T(), err)
			// time.Sleep(60000000000000)
			kubeProvisioningClient, err := r.client.GetKubeAPIProvisioningClient()
			require.NoError(r.T(), err)

			result, err := kubeProvisioningClient.Clusters(namespace).Watch(context.TODO(), metav1.ListOptions{
				FieldSelector:  "metadata.name=" + clusterName,
				TimeoutSeconds: &defaults.WatchTimeoutSeconds,
			})
			require.NoError(r.T(), err)

			checkFunc := clusters.IsProvisioningClusterReady

			err = wait.WatchWait(result, checkFunc)
			assert.NoError(r.T(), err)
			assert.Equal(r.T(), clusterName, clusterResp.ObjectMeta.Name)

			clusterToken, err := clusters.CheckServiceAccountTokenSecret(client, clusterName)
			require.NoError(r.T(), err)
			assert.NotEmpty(r.T(), clusterToken)
		})
	}
}

func TestCertRotationSuite(t *testing.T) {
	suite.Run(t, new(CertRotationTestSuite))
}

func (r *CertRotationTestSuite) TestProivioning() {
	// time.Sleep(6000000000000)
	// for _, providerName := range r.providers {
	// 	provider := CreateProvider(providerName)
	// 	r.ProvisioningRKE2Cluster(provider)
	// }
	provider := CreateProvider("aws")
	r.ProvisionRKE2Cluster(provider)
}
