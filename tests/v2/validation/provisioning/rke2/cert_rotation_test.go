package rke2

import (
	"fmt"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/cloudcredentials"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/machinepools"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CertRotationTestSuite struct {
	suite.Suite
	session     *session.Session
	client      *rancher.Client
	config      *rancher.Config
	clusterName string
	namespace   string
}

func (p *CertRotationTestSuite) TearDownSuite() {
	p.session.Cleanup()
}

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

		clusterName := AppendRandomString(fmt.Sprintf("%s-%s", r.clusterName, provider.Name))
		generatedPoolName := fmt.Sprintf("nc-%s-pool1-", clusterName)
		machinePoolConfig := provider.MachinePoolFunc(generatedPoolName, namespace)

		machineConfigResp, err := machinepools.CreateMachineConfig(provider.MachineConfig, machinePoolConfig, testSessionClient)
		require.NoError(r.T(), err)

		machinePools := machinepools.RKEMachinePoolSetup(nodesAndRoles, machineConfigResp)

		cluster := clusters.NewRKE2ClusterConfig(clusterName, namespace, "calico", credential.ID, kubeVersion, machinePools)

		// clusterResp, err := clusters.CreateRKE2Cluster(testSessionClient, cluster)
		// require.NoError(r.T(), err)
		clusters.CreateRKE2Cluster(testSessionClient, cluster)

	})
}
