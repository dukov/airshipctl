package k8s

import (
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/printers"
	"k8s.io/kubectl/pkg/cmd/apply"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// NewApplyOptions is a helper function that Creates ApplyOptions of kubectl apply module
// Values set here, are default, and do not conflict with each other, can be used if you
// need `kubectl apply` functionality without calling executing command in shell
// To function properly, you may need to specify files from where to read the resources:
// DeleteOptions.Filenames of returned object has to be set for that
func NewApplyOptions(f cmdutil.Factory, streams genericclioptions.IOStreams) (*apply.ApplyOptions, error) {

	o := apply.NewApplyOptions(streams)
	o.ServerSideApply = false
	o.ForceConflicts = false

	// TODO (k.kalynovskyi) add unit test for this function
	o.ToPrinter = func(operation string) (printers.ResourcePrinter, error) {
		o.PrintFlags.NamePrintFlags.Operation = operation
		if o.DryRun {
			o.PrintFlags.Complete("%s (dry run)")
		}
		if o.ServerDryRun {
			o.PrintFlags.Complete("%s (server dry run)")
		}
		return o.PrintFlags.ToPrinter()
	}

	var err error
	o.Recorder, err = o.RecordFlags.ToRecorder()
	if err != nil {
		return nil, err
	}

	o.DiscoveryClient, err = f.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	dynamicClient, err := f.DynamicClient()
	if err != nil {
		return nil, err
	}

	o.DeleteOptions = o.DeleteFlags.ToOptions(dynamicClient, o.IOStreams)
	// This can only fail if ToDiscoverClient() function fails
	o.OpenAPISchema, err = f.OpenAPISchema()
	if err != nil {
		return nil, err
	}

	o.Validator, err = f.Validator(false)
	if err != nil {
		return nil, err
	}

	o.Builder = f.NewBuilder()
	o.Mapper, err = f.ToRESTMapper()
	if err != nil {
		return nil, err
	}

	o.DynamicClient = dynamicClient

	o.Namespace, o.EnforceNamespace, err = f.ToRawKubeConfigLoader().Namespace()
	if err != nil {
		return nil, err
	}
	return o, nil
}
