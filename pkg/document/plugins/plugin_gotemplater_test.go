package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/testutil"
)

func TestGotemplater(t *testing.T) {
	t.Run("Transform data", func(t *testing.T) {
		bundle, err := document.NewBundle(testutil.SetupTestFs(t, "testdata/gotemplater"), "/", "/")
		require.NoError(t, err)
		var resourceToRender document.Document
		var dataSource document.Document
		var renderString string
		var srcMap map[string]interface{}
		var src []byte
		resourceToRender, err = bundle.GetByName("resourceToRender")
		require.NoError(t, err)
		dataSource, err = bundle.GetByName("dataSource")
		require.NoError(t, err)
		renderString, err = resourceToRender.GetString("data")
		require.NoError(t, err)
		srcMap, err = dataSource.GetMap("spec.someData")
		require.NoError(t, err)
		src, err = yaml.Marshal(srcMap)
		require.NoError(t, err)

		assert.Equal(t, string(src), renderString)
	})
	t.Run("Bad template", func(t *testing.T) {
		_, err := document.NewBundle(testutil.SetupTestFs(t, "testdata/gotemplaterbadtemplate"), "/", "/")
		assert.Error(t, err)
	})
	t.Run("NonExistent resource link", func(t *testing.T) {
		_, err := document.NewBundle(testutil.SetupTestFs(t, "testdata/gotemplaterbadresource"), "/", "/")
		assert.Error(t, err)
	})
}
