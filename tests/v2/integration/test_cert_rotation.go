package integration

import (
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CertRotationTestSuite struct {
	suite.Suite
	session     *session.Session
	client      *rancher.Client
	config      *Config
	clusterName string
	namespace   string
}

func (p *CertRotationTestSuite) TearDownSuite() {
	p.session.Cleanup()
}

func (r *CertRotationTestSuite) SetupSuite() {
	testSession := session.NewSession(r.T())
	r.session = testSession

	r.config = new(Config)
	config.LoadConfig(ConfigurationFileKey, r.config)

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client

	r.clusterName = r.client.RancherConfig.ClusterName
	r.namespace = r.client.RancherConfig.ClusterName
}
