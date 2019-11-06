package k8s_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/k8s"
	"opendev.org/airship/airshipctl/testutil"
	k8stest "opendev.org/airship/airshipctl/testutil/k8sutils"
)

func TestNewKubectlFromKubeconfigPath(t *testing.T) {
	kubectl := k8s.NewKubectlFromKubeconfigPath(kubeconfigPath)

	assert.NotNil(t, kubectl.Factory)
	assert.NotNil(t, kubectl.FileSystem)
	assert.NotNil(t, kubectl.IOStreams)
}

func TestApply(t *testing.T) {

	f := k8stest.NewFakeFactoryForRC(t, filenameRC)
	defer f.Cleanup()
	kubectl := k8s.NewKubectlFromKubeconfigPath(kubeconfigPath)
	kubectl.Factory = f
	ao, err := kubectl.ApplyOptions()
	require.NoError(t, err, "failed to get documents from bundle")
	ao.DryRun = true

	b := testutil.NewTestBundle(t, fixtureDir)
	docs, err := b.GetByAnnotation("airshipit.org/initinfra")
	require.NoError(t, err, "failed to get documents from bundle")

	tests := []struct {
		name        string
		expectedErr error
		fs          k8s.FileSystem
	}{
		{
			expectedErr: nil,
			fs: MockFileSystem{
				MockRemoveAll: func() error { return nil },
				MockTempFile: func() (k8s.File, error) {
					return TestFile{
						MockName:  func() string { return filenameRC },
						MockWrite: func() (int, error) { return 0, nil },
						MockClose: func() error { return nil },
					}, nil
				},
			},
		},
		{
			expectedErr: writeOutError,
			fs: MockFileSystem{
				MockTempFile: func() (k8s.File, error) { return nil, writeOutError }},
		},
		{
			expectedErr: TempFileError,
			fs: MockFileSystem{
				MockRemoveAll: func() error { return nil },
				MockTempFile: func() (k8s.File, error) {
					return TestFile{
						MockWrite: func() (int, error) { return 0, TempFileError },
						MockName:  func() string { return filenameRC },
						MockClose: func() error { return nil },
					}, nil
				}},
		},
	}
	for _, test := range tests {
		kubectl.FileSystem = test.fs
		assert.Equal(t, kubectl.Apply(docs, ao), test.expectedErr)
	}

}
