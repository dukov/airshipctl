/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
