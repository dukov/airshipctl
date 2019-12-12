package document_test

import (
	"testing"

	"opendev.org/airship/airshipctl/cmd/document"
	"opendev.org/airship/airshipctl/testutil"
)

func TestRender(t *testing.T) {
	tests := []*testutil.CmdTest{
		{
			Name:    "render-with-help",
			CmdLine: "-h",
			Cmd:     document.NewRenderCommand(nil),
		},
		{
			Name:    "render-with-defaults",
			CmdLine: "testdata/render",
			Cmd:     document.NewRenderCommand(nil),
		},
		{
			Name: "render-with-all-flags",
			CmdLine: `testdata/render
						-l app=helm
						-a airshipit.org/clustertype=ephemeral
						-g extensions
						-v v1beta1
						-k Deployment`,
			Cmd: document.NewRenderCommand(nil),
		},
		{
			Name:    "render-with-non-existent-label",
			CmdLine: "testdata/render -l doesnot=exist",
			Cmd:     document.NewRenderCommand(nil),
		},
		{
			Name:    "render-with-kind",
			CmdLine: "testdata/render -k BareMetalHost",
			Cmd:     document.NewRenderCommand(nil),
		},
	}
	for _, tt := range tests {
		testutil.RunTest(t, tt)
	}
}
