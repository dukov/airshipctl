package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/testutil"
)

func TestTreeInserter(t *testing.T) {
	t.Run("Transform data", func(t *testing.T) {
		bundle, err := document.NewBundle(testutil.SetupTestFs(t, "testdata/treeinserter"), "/", "/")
		assert.NoError(t, err)
		var resourceToRender document.Document
		var dataSource document.Document
		var renderString string
		var srcMap map[string]interface{}
		var src []byte
		resourceToRender, err = bundle.GetByName("resourceToRender")
		assert.NoError(t, err)
		dataSource, err = bundle.GetByName("dataSource")
		assert.NoError(t, err)
		renderString, err = resourceToRender.GetString("data")
		assert.NoError(t, err)
		srcMap, err = dataSource.GetMap("spec.someData")
		assert.NoError(t, err)
		src, err = yaml.Marshal(srcMap)
		assert.NoError(t, err)

		assert.Equal(t, string(src), renderString)
	})
}
