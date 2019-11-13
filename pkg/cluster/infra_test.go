package cluster_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"

	"opendev.org/airship/airshipctl/pkg/cluster"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/testutil/k8sutils"
)

const (
	kubeconfigPath    = "../../testdata/k8s/kubeconfig.yaml"
	airshipConfigPath = "../../testdata/k8s/config.yaml"

	filenameRC        = "../../testdata/k8s/replicationcontroller.yaml"
	airshipConfigFile = "../../testdata/k8s/config.yaml"
)

var (
	DynamicClientError = errors.New("DynamicClientError")
)

func TestNewInfra(t *testing.T) {
	rs := makeNewFakeRootSettings(t, kubeconfigPath, airshipConfigFile)
	infra := cluster.NewInfra(rs)

	assert.NotNil(t, infra.AirshipCTLSettings)
}

func TestComplete(t *testing.T) {
	rs := makeNewFakeRootSettings(t, kubeconfigPath, airshipConfigFile)
	infra := cluster.NewInfra(rs)
	infra.Complete()

	assert.NotNil(t, infra.FileSystem)
	assert.NotNil(t, infra.Kubectl)
}

func TestDeploy(t *testing.T) {
	rs := makeNewFakeRootSettings(t, kubeconfigPath, airshipConfigFile)
	tf := k8sutils.NewFakeFactoryForRC(t, filenameRC)
	defer tf.Cleanup()

	infra := cluster.NewInfra(rs)
	infra.Complete()
	infra.DryRun = true

	tests := []struct {
		theInfra      *cluster.Infra
		factory       cmdutil.Factory
		prune         bool
		expectedError error
	}{
		{
			factory:       k8sutils.NewMockKubectlFactory().WithDynamicClientByError(nil, DynamicClientError),
			expectedError: DynamicClientError,
		},
		{
			expectedError: nil,
			prune:         false,
			factory:       tf,
		},
		{
			expectedError: nil,
			prune:         true,
			factory:       tf,
		},
	}

	for _, test := range tests {
		// TODO (kkalynovskyi) develop Prune tests, and check actual result from StdOut
		infra.Prune = test.prune
		infra.Factory = test.factory
		actualErr := infra.Deploy()
		assert.Equal(t, actualErr, test.expectedError)
	}
}

// MakeNewFakeRootSettings takes kubeconfig path and directory path to fixture dir as argument.
func makeNewFakeRootSettings(t *testing.T, kp string, dir string) *environment.AirshipCTLSettings {
	t.Helper()
	rs := &environment.AirshipCTLSettings{}

	akp, err := filepath.Abs(kp)
	require.NoError(t, err)

	adir, err := filepath.Abs(dir)
	require.NoError(t, err)

	rs.SetAirshipConfigPath(adir)
	rs.SetKubeConfigPath(akp)

	rs.InitConfig()
	return rs
}
