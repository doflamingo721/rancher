package rke2

import (
	"context"
	"fmt"
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
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

type RKE2EtcdSnapshotRestoreTestSuite struct {
	suite.Suite
	session            *session.Session
	client             *rancher.Client
	clusterName        string
	namespace          string
	kubernetesVersions []string
	cnis               []string
	providers          []string
	nodesAndRoles      []machinepools.NodeRoles
}

var phases = []rkev1.ETCDSnapshotPhase{
	rkev1.ETCDSnapshotPhaseStarted,
	rkev1.ETCDSnapshotPhaseShutdown,
	rkev1.ETCDSnapshotPhaseRestore,
	rkev1.ETCDSnapshotPhaseRestartCluster,
	rkev1.ETCDSnapshotPhaseFinished,
}

// func (p *RKE2EtcdSnapshotRestoreTestSuite) TearDownSuite() {
// 	p.session.Cleanup()
// }

func (r *RKE2EtcdSnapshotRestoreTestSuite) SetupSuite() {
	testSession := session.NewSession(r.T())
	r.session = testSession

	clustersConfig := new(provisioning.Config)
	config.LoadConfig(provisioning.ConfigurationFileKey, clustersConfig)

	r.kubernetesVersions = clustersConfig.KubernetesVersions
	r.cnis = clustersConfig.CNIs
	r.providers = clustersConfig.Providers
	r.nodesAndRoles = clustersConfig.NodesAndRoles

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client

	r.clusterName = r.client.RancherConfig.ClusterName
	r.namespace = r.client.RancherConfig.ClusterName
}

func (r *RKE2EtcdSnapshotRestoreTestSuite) TestEtcdSnapshotRestoreFreshCluster(provider Provider, kubeVersion string, cni string, nodesAndRoles []machinepools.NodeRoles, credential *cloudcredentials.CloudCredential) {
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

		cluster := clusters.NewRKE2ClusterConfig(clusterName, namespace, cni, credential.ID, kubeVersion, machinePools)

		clusterResp, err := clusters.CreateRKE2Cluster(testSessionClient, cluster)
		require.NoError(r.T(), err)

		kubeProvisioningClient, err := r.client.GetKubeAPIProvisioningClient()
		require.NoError(r.T(), err)

		result, err := kubeProvisioningClient.Clusters(namespace).Watch(context.TODO(), metav1.ListOptions{
			FieldSelector:  "metadata.name=" + clusterName,
			TimeoutSeconds: &defaults.WatchTimeoutSeconds,
		})
		require.NoError(r.T(), err)

		checkFunc := clusters.IsProvisioningClusterReady
		fmt.Println("CheckFunc ")
		fmt.Println("Before WaitWatch ")
		err = wait.WatchWait(result, checkFunc)
		fmt.Println("After WaitWatch ")
		assert.NoError(r.T(), err)
		assert.Equal(r.T(), clusterName, clusterResp.ObjectMeta.Name)

		clusterToken, err := clusters.CheckServiceAccountTokenSecret(testSessionClient, clusterName)
		require.NoError(r.T(), err)
		assert.NotEmpty(r.T(), clusterToken)

		cluster, err = r.client.Provisioning.Cluster.ByID(clusterResp.ID)
		require.NoError(r.T(), err)
		require.NotNil(r.T(), cluster.Status)

		// require.NoError(r.T(), r.createSnapshot(clusterName, 3))
		// time.Sleep(30 * time.Second)
		require.NoError(r.T(), r.restoreSnapshot(clusterName, "auto--aws-hcadx-on-demand-auto--aws-hcadx-pool0-f737293b-728c74", 3, "all"))
		fmt.Println("After createSnapshot call")

	})
}

func (r *RKE2EtcdSnapshotRestoreTestSuite) createSnapshot(id string, generation int) error {
	fmt.Println("Inside snapshot function")
	kubeProvisioningClient, err := r.client.GetKubeAPIProvisioningClient()
	require.NoError(r.T(), err)

	cluster, err := kubeProvisioningClient.Clusters(namespace).Get(context.TODO(), id, metav1.GetOptions{})
	if err != nil {
		return err
	}

	cluster.Spec.RKEConfig.ETCDSnapshotCreate = &rkev1.ETCDSnapshotCreate{
		Generation: generation,
	}

	fmt.Println("etcdsnapshot:  ", cluster.Spec.RKEConfig.ETCDSnapshotCreate)
	// time.Sleep(100000000000000000)

	fmt.Println("before update")
	cluster, err = kubeProvisioningClient.Clusters(namespace).Update(context.TODO(), cluster, metav1.UpdateOptions{})
	if err != nil {
		return err
	}

	fmt.Println("after update", cluster)

	fmt.Println("Before WaitWatch ")

	result, err := kubeProvisioningClient.Clusters(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + cluster.ObjectMeta.Name,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	require.NoError(r.T(), err)

	checkFunc := clusters.IsProvisioningClusterReady
	fmt.Println("CheckFunc ")

	err = wait.WatchWait(result, checkFunc)
	fmt.Println("After WaitWatch ")
	assert.NoError(r.T(), err)
	// assert.Equal(r.T(), clusterresponse.Status.ClusterName, clusterresponse.ObjectMeta.Name)

	return nil
}

func (r *RKE2EtcdSnapshotRestoreTestSuite) restoreSnapshot(id string, name string, generation int, restoreconfig string) error {
	fmt.Println("Inside snapshotRestore function")
	kubeProvisioningClient, err := r.client.GetKubeAPIProvisioningClient()
	require.NoError(r.T(), err)

	cluster, err := kubeProvisioningClient.Clusters(namespace).Get(context.TODO(), id, metav1.GetOptions{})
	if err != nil {
		return err
	}
	snapshot.
		cluster.Spec.RKEConfig.ETCDSnapshotRestore = &rkev1.ETCDSnapshotRestore{
		Name:             "auto--aws-hcadx-on-demand-auto--aws-hcadx-pool0-f737293b-029706",
		Generation:       generation,
		RestoreRKEConfig: "all",
	}

	fmt.Println("before cluster update")
	cluster, err = kubeProvisioningClient.Clusters(namespace).Update(context.TODO(), cluster, metav1.UpdateOptions{})
	if err != nil {
		return err
	}
	fmt.Println("after update", cluster)

	fmt.Println("restore Before WaitWatch ")
	results, err := kubeProvisioningClient.Clusters(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector:  "metadata.name=" + cluster.ObjectMeta.Name,
		TimeoutSeconds: &defaults.WatchTimeoutSeconds,
	})
	require.NoError(r.T(), err)

	checkFuncs := clusters.IsProvisioningClusterReady
	fmt.Println("CheckFunc ")

	err = wait.WatchWait(results, checkFuncs)
	fmt.Println("restore After WaitWatch ")
	assert.NoError(r.T(), err)

	return nil
}

func (r *RKE2EtcdSnapshotRestoreTestSuite) TestEtcdSnapshotRestore() {
	for _, providerName := range r.providers {
		subSession := r.session.NewSession()

		provider := CreateProvider(providerName)

		client, err := r.client.WithSession(subSession)
		require.NoError(r.T(), err)

		cloudCredential, err := provider.CloudCredFunc(client)
		require.NoError(r.T(), err)

		for _, kubernetesVersion := range r.kubernetesVersions {
			for _, cni := range r.cnis {
				r.TestEtcdSnapshotRestoreFreshCluster(provider, kubernetesVersion, cni, r.nodesAndRoles, cloudCredential)
			}
		}

		subSession.Cleanup()
	}
}

func TestEtcdSnapshotRestore(t *testing.T) {
	suite.Run(t, new(RKE2EtcdSnapshotRestoreTestSuite))
}
