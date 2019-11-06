package k8s

import (
	"os"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"sigs.k8s.io/kustomize/v3/pkg/fs"

	"opendev.org/airship/airshipctl/pkg/document"
	utilyaml "opendev.org/airship/airshipctl/pkg/util/yaml"
)

// Kubectl container holds Factory, Streams and FileSystem to
// interact with upstream kubectl objects and serves as abstraction to kubectl project
type Kubectl struct {
	cmdutil.Factory
	genericclioptions.IOStreams
	FileSystem
}

// NewKubectlFromKubeconfigPath builds an instance
// of Kubectl struct from Path to kubeconfig file
func NewKubectlFromKubeconfigPath(kp string) *Kubectl {
	s := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	kf := genericclioptions.NewConfigFlags(false)
	kf.KubeConfig = &kp

	return &Kubectl{
		Factory:    cmdutil.NewFactory(kf),
		IOStreams:  s,
		FileSystem: Buffer{FileSystem: fs.MakeRealFS()},
	}
}

// Apply is abstraction to kubectl apply command
func (kubectl *Kubectl) Apply(docs []document.Document, ao *apply.ApplyOptions) error {

	tf, err := kubectl.FileSystem.TempFile("initinfra")
	if err != nil {
		return err
	}
	defer kubectl.FileSystem.RemoveAll(tf.Name())
	defer tf.Close()
	for _, doc := range docs {

		// Write out documents to temporary file
		err = utilyaml.WriteOut(tf, doc)
		if err != nil {
			return err
		}
	}
	ao.DeleteOptions.Filenames = []string{tf.Name()}
	return ao.Run()
}

// ApplyOptions is a wrapper over kubectl ApplyOptions, used to build
// new options from the factory and iostreams defined in Kubectl container
func (kubectl *Kubectl) ApplyOptions() (*apply.ApplyOptions, error) {
	return NewApplyOptions(kubectl.Factory, kubectl.IOStreams)
}
