package document

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/document/render"
	"opendev.org/airship/airshipctl/pkg/environment"
)

var (
	renderExample = `
#Get all documents containing labels "app=helm" and "service=tiller"
airshipctl document render -l app=helm -l service=armada

#Same with raw filter
airshipctl document render -f 'metadata.labels.app == "helm" && metadata.labels.service == "armada"'

#Get documents with particular value in cirtain filed referenced
#by JSON path 'spec.template.spec.image'
airshipctl document render -f 'spec.template.spec.image == "ubuntu:xenial"'

#Get all Secrets or ConfigMaps
airshipctl document render -f 'kind == "Secret" || kind == "ConfigMap"'`
)

// InitFlags add flags for document render sub-command
func initRenderFlags(settings *render.Settings, cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringArrayVarP(&settings.Label, "label", "l", nil, "Filter documents by Labels")
	flags.StringArrayVarP(&settings.Annotation, "annotation", "a", nil, "Filter documents by Annotations")
	flags.StringArrayVarP(&settings.GroupVersion, "apiversion", "g", nil, "Filter documents by API version")
	flags.StringArrayVarP(&settings.Kind, "kind", "k", nil, "Filter documents by Kinds")
	flags.StringVarP(&settings.RawFilter, "filter", "f", "", "Logical expression for document filtering")
}

// NewRenderCommand create a new command for document rendering
func NewRenderCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	renderSettings := &render.Settings{AirshipCTLSettings: rootSettings}
	renderCmd := &cobra.Command{
		Use:     "render",
		Short:   "Render documents from model",
		Example: renderExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			manifest, err := renderSettings.Config().CurrentContextManifest()
			if err != nil {
				return err
			}
			return renderSettings.Render(manifest.TargetPath, cmd.OutOrStdout())
		},
	}

	initRenderFlags(renderSettings, renderCmd)
	return renderCmd
}
