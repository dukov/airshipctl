package k8s_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/kustomize/v3/pkg/fs"

	"opendev.org/airship/airshipctl/pkg/k8s"
	k8stest "opendev.org/airship/airshipctl/testutil/k8sutils"
)

var (
	codec          = scheme.Codecs.LegacyCodec(scheme.Scheme.PrioritizedVersionsAllGroups()...)
	filenameRC     = "../../testdata/k8s/replicationcontroller.yaml"
	kubeconfigPath = "../../testdata/k8s/kubeconfig.yaml"
	fixtureDir     = "../../testdata/k8s/"

	testStreams        = genericclioptions.NewTestIOStreamsDiscard()
	ToDiscoveryError   = errors.New("ToDiscoveryError")
	DynamicClientError = errors.New("DynamicClientError")
	ValidateError      = errors.New("ValidateError")
	ToRESTMapperError  = errors.New("ToRESTMapperError")
	NamespaceError     = errors.New("NamespaceError")
	writeOutError      = errors.New("writeOutError")
	TempFileError      = errors.New("TempFileError")
)

type MockFileSystem struct {
	MockRemoveAll func() error
	MockTempFile  func() (k8s.File, error)
	fs.FileSystem
}

func (fsys MockFileSystem) RemoveAll(name string) error { return fsys.MockRemoveAll() }
func (fsys MockFileSystem) TempFile(prefix string) (k8s.File, error) {
	return fsys.MockTempFile()
}

type TestFile struct {
	k8s.File
	MockName  func() string
	MockWrite func() (int, error)
	MockClose func() error
}

func (f TestFile) Name() string              { return f.MockName() }
func (f TestFile) Write([]byte) (int, error) { return f.MockWrite() }
func (f TestFile) Close() error              { return f.MockClose() }

func TestApplyOptionsRun(t *testing.T) {
	f := k8stest.NewFakeFactoryForRC(t, filenameRC)
	defer f.Cleanup()

	streams := genericclioptions.NewTestIOStreamsDiscard()

	aa, err := k8s.NewApplyOptions(f, streams)
	require.NoError(t, err, "Could not build ApplyAdapter")
	aa.DryRun = true
	aa.DeleteOptions.Filenames = []string{filenameRC}
	assert.NoError(t, aa.Run())
}

func TestNewApplyOptionsFactoryFailures(t *testing.T) {

	tests := []struct {
		f             cmdutil.Factory
		expectedError error
	}{
		{
			f:             k8stest.NewMockKubectlFactory().WithToDiscoveryClientByError(nil, ToDiscoveryError),
			expectedError: ToDiscoveryError,
		},
		{
			f:             k8stest.NewMockKubectlFactory().WithDynamicClientByError(nil, DynamicClientError),
			expectedError: DynamicClientError,
		},
		{
			f:             k8stest.NewMockKubectlFactory().WithValidatorByError(nil, ValidateError),
			expectedError: ValidateError,
		},
		{
			f:             k8stest.NewMockKubectlFactory().WithToRESTMapperByError(nil, ToRESTMapperError),
			expectedError: ToRESTMapperError,
		},
		{
			f: k8stest.NewMockKubectlFactory().
				WithToRawKubeConfigLoaderByError(k8stest.
					NewMockClientConfig().
					WithNamespace("", false, NamespaceError)),
			expectedError: NamespaceError,
		},
	}
	for _, test := range tests {
		_, err := k8s.NewApplyOptions(test.f, testStreams)
		assert.Equal(t, err, test.expectedError)
	}
}
