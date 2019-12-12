package document

import (
	"io"

	"github.com/spf13/cobra"

	"opendev.org/airship/airshipctl/pkg/document"
	"opendev.org/airship/airshipctl/pkg/environment"
	"opendev.org/airship/airshipctl/pkg/errors"
	utilyaml "opendev.org/airship/airshipctl/pkg/util/yaml"

	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/types"
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

// Render prints out filtered documents
func (s *Settings) Render(path string, out io.Writer) error {
	docBundle, err := document.NewBundle(fs.MakeRealFS(), path, "")
	if err != nil {
		return err
	}
	filter := types.Selector{
		Gvk: gvk.Gvk{
			Group:   s.Group,
			Version: s.Version,
			Kind:    s.Kind,
		},
		AnnotationSelector: s.Annotation,
		LabelSelector:      s.Label,
	}
	var docs []document.Document
	docs, err = docBundle.Select(filter)
	if err != nil {
		return err
	}

	for _, doc := range docs {
		err = utilyaml.WriteOut(out, doc)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewRenderCommand create a new command for document rendering
func NewRenderCommand(rootSettings *environment.AirshipCTLSettings) *cobra.Command {
	renderSettings := &Settings{AirshipCTLSettings: rootSettings}
	renderCmd := &cobra.Command{
		Use:   "render",
		Short: "Render documents from model",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO (dukov) Replace this code with reading repo path from
			// configuration context
			if len(args) == 0 {
				return errors.ErrWrongConfig{}
			}
			return renderSettings.Render(args[0], cmd.OutOrStdout())
		},
	}

	renderSettings.InitFlags(renderCmd)
	return renderCmd
}
