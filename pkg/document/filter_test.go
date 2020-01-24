package document_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/testutil"

	"sigs.k8s.io/kustomize/v3/pkg/resid"
)

func TestEvaluateFilter(t *testing.T) {
	bundle := testutil.NewTestBundle(t, "testdata")

	tests := []struct {
		filter       string
		expectedDocs []string
	}{
		{
			filter: `kind == "BareMetalHost" && metadata.annotations["airshipit.org/clustertype"] == "ephemeral"`,
			expectedDocs: []string{
				"metal3.io_v1alpha1_BareMetalHost|~X|master-0",
			},
		},
		{
			filter: `kind == "BareMetalHost" && spec.bootMACAddress == "01:3b:8b:0c:ec:8b"`,
			expectedDocs: []string{
				"metal3.io_v1alpha1_BareMetalHost|~X|master-1",
			},
		},
		{
			filter: `metadata.name == "master-1" || metadata.name == "master-0"`,
			expectedDocs: []string{
				"metal3.io_v1alpha1_BareMetalHost|~X|master-0",
				"metal3.io_v1alpha1_BareMetalHost|~X|master-1",
			},
		},
		{
			filter: `spec.names.shortNames[0] == "wf"`,
			expectedDocs: []string{
				"apiextensions.k8s.io_v1beta1_CustomResourceDefinition|~X|workflows.argoproj.io",
			},
		},
	}

	for _, tt := range tests {
		filteredBundle, err := document.EvaluateExpressionFilter(tt.filter, bundle)
		require.NoError(t, err)
		docs, err := filteredBundle.GetAllDocuments()
		require.NoError(t, err)
		assert.Equal(t, len(tt.expectedDocs), len(docs))
		for _, res := range tt.expectedDocs {
			rid := resid.FromString(res)
			_, err := filteredBundle.GetKustomizeResourceMap().GetById(rid)
			assert.NoError(t, err)
		}
	}
}
