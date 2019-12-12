package render_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/document/render"
)

func TestRender(t *testing.T) {
	fixturePath := "../testdata"
	tests := []struct {
		settings   *render.Settings
		expResFile string
		expErr     error
	}{
		{
			settings:   &render.Settings{},
			expResFile: "noFilter.yaml",
			expErr:     nil,
		},
		{
			settings: &render.Settings{
				Label:        []string{"app=helm"},
				Annotation:   []string{"airshipit.org/clustertype=ephemeral"},
				GroupVersion: []string{"v1"},
				Kind:         []string{"Service"},
			},
			expResFile: "allFilters.yaml",
			expErr:     nil,
		},
		{
			settings: &render.Settings{
				Label: []string{"app=helm", "name=tiller"},
			},
			expResFile: "multiLabels.yaml",
			expErr:     nil,
		},
		{
			settings: &render.Settings{
				RawFilter: `kind=="BareMetalHost"&&metadata.name=="master-1"`,
			},
			expResFile: "rawFilter.yaml",
			expErr:     nil,
		},
		{
			settings: &render.Settings{
				Kind:      []string{"Service"},
				RawFilter: `kind=="BareMetalHost"`,
			},
			expResFile: "",
			expErr:     document.ErrWrongRenderArgs{},
		},
		{
			settings: &render.Settings{
				Label: []string{"app"},
			},
			expResFile: "",
			expErr:     document.ErrBadRenderArgFormat{Arg: "app"},
		},
	}

	for _, tt := range tests {
		var expectedOut []byte
		var err error
		if tt.expResFile != "" {
			expectedOut, err = ioutil.ReadFile(path.Join("testdata", "expected", tt.expResFile))
			require.NoError(t, err)
		}
		out := &bytes.Buffer{}
		err = tt.settings.Render(fixturePath, out)
		assert.Equal(t, tt.expErr, err)
		assert.Equal(t, expectedOut, out.Bytes())
	}
}
