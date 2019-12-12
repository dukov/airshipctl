package document_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"opendev.org/airship/airshipctl/cmd/document"
	"opendev.org/airship/airshipctl/pkg/config"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/testutil"
)

func TestRender(t *testing.T) {
	cfg := config.InitConfig(t)
	cfg.Manifests["test"] = &config.Manifest{TargetPath: "testdata/render"}
	cfg.CurrentContext = "def_ephemeral"
	ctx, err := cfg.GetContext("def_ephemeral")
	require.NoError(t, err)
	ctx.Manifest = "test"

	settings := &environment.AirshipCTLSettings{}
	settings.SetConfig(cfg)

	tests := []*testutil.CmdTest{
		{
			Name:    "render-with-help",
			CmdLine: "-h",
			Cmd:     document.NewRenderCommand(nil),
		},
		{
			Name: "render-with-all-flags",
			CmdLine: `-l app=helm
					  -a airshipit.org/clustertype=ephemeral
					  -g extensions/v1beta1
					  -k Deployment`,
			Cmd: document.NewRenderCommand(settings),
		},
		{
			Name:    "render-with-multiple-labels",
			CmdLine: "-l app=helm -l name=tiller",
			Cmd:     document.NewRenderCommand(settings),
		},
		{
			Name:    "render-with-raw-filter",
			CmdLine: `-f kind=="BareMetalHost"&&metadata.name=="master-1"`,
			Cmd:     document.NewRenderCommand(settings),
		},
	}
	for _, tt := range tests {
		testutil.RunTest(t, tt)
	}
}
