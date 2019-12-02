package document

import (
	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/errors"
)

// Settings for document rendering
type Settings struct {
	*environment.AirshipCTLSettings
	// Label filter documents by label string
	Label string
	// Annotation filter documents by annotation string
	Annotation string
	// Group filter documents by API group
	Group string
	// Version filter documents by API version
	Version string
	// Kind filter documents by document kind
	Kind string
}

// InitFlags add flags for document render sub-command
func (s *Settings) InitFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringVarP(&s.Label, "label", "l", "", "Label for document filtering")
	flags.StringVarP(&s.Annotation, "annotation", "a", "", "Annotation for document filtering")
	flags.StringVarP(&s.Group, "group", "g", "", "Group for document filtering")
	flags.StringVarP(&s.Version, "version", "v", "", "Version for document filtering")
	flags.StringVarP(&s.Kind, "kind", "k", "", "Kind for document filtering")
}

// NewRenderCommand create a new command for document rendering
func NewRenderCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	renderSettings := &Settings{AirshipCTLSettings: rootSettings}
	renderCmd := &cobra.Command{
		Use:   "render",
		Short: "Render documents from model",
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.ErrNotImplemented{}
		},
	}

	renderSettings.InitFlags(renderCmd)
	return renderCmd
}
