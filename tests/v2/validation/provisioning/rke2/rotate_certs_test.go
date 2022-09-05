package rke2

import (
	"context"
	"fmt"
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/cloudcredentials"
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
	session     *session.Session
	client      *rancher.Client
	config      *rancher.Config
	clusterName string
	namespace   string
}

// func (p *CertRotationTestSuite) TearDownSuite() {
// 	p.session.Cleanup()
// }

func (r *CertRotationTestSuite) SetupSuite() {
	testSession := session.NewSession(r.T())
	r.session = testSession

	r.config = new(rancher.Config)
	config.LoadConfig(rancher.ConfigurationFileKey, r.config)

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client

	r.clusterName = r.client.RancherConfig.ClusterName
	r.namespace = r.client.RancherConfig.ClusterName
}

func (r *CertRotationTestSuite) TestCertRotationFreshCluster(provider Provider, kubeVersion string, nodesAndRoles []machinepools.NodeRoles, credential *cloudcredentials.CloudCredential) {
	name := fmt.Sprintf("Provider_%s/Kubernetes_Version_%s/Nodes_%v", provider.Name, kubeVersion, nodesAndRoles)
	r.Run(name, func() {
		testSession := session.NewSession(r.T())
		defer testSession.Cleanup()

		testSessionClient, err := r.client.WithSession(testSession)
		require.NoError(r.T(), err)

		clusterName := provisioning.AppendRandomString(fmt.Sprintf("%s-%s", r.clusterName, provider.Name))
		generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
		machinePoolConfig := provider.MachinePoolFunc(generatedPoolName, namespace)

		machineConfigResp, err := machinepools.CreateMachineConfig(provider.MachineConfig, machinePoolConfig, testSessionClient)
		require.NoError(r.T(), err)

		machinePools := machinepools.RKEMachinePoolSetup(nodesAndRoles, machineConfigResp)

		cluster := clusters.NewRKE2ClusterConfig(clusterName, namespace, "calico", "cc-2rrgf", "v1.24.2-rancher1-1", machinePools)

		//clusters.CreateRKE2Cluster(testSessionClient, cluster)

		clusterResp, err := clusters.CreateRKE2Cluster(testSessionClient, cluster)
		require.NoError(r.T(), err)

		kubeRKEClient, err := r.client.GetKubeAPIRKEClient()
		require.NoError(r.T(), err)

		result, err := kubeRKEClient.RKEControlPlanes(namespace).Watch(context.TODO(), metav1.ListOptions{

			FieldSelector:  "metadata.name=" + cluster.ID,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		})
		require.NoError(r.T(), err)

		checkFunc := clusters.IsProvisioningClusterReady

		err = wait.WatchWait(result, checkFunc)
		assert.NoError(r.T(), err)
		assert.Equal(r.T(), clusterName, clusterResp.ObjectMeta.Name)

	})
}

func TestCertRotationSuite(t *testing.T) {
	suite.Run(t, new(CertRotationTestSuite))
}
