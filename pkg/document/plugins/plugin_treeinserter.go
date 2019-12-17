package plugins

import (
	"bytes"
	"html/template"

	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/yaml"
)

func init() {
	PluginRegistry["TreeInserter"] = NewTreeInserterPlugin()
}

// TreeInserter plugin to insert map subtree from specified resource
type TreeInserter struct {
	resid.ResId
}

// Config loads configuration from plugin config
func (ti *TreeInserter) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	return yaml.Unmarshal(c, ti)
}

// Transform applies transformation to all resources
func (ti *TreeInserter) Transform(m resmap.ResMap) error {
	for _, r := range m.Resources() {
		rm := r.Map()
		ModifyHashStrings(rm, m, renderString)
	}
	return nil
}

// NewTreeInserterPlugin returns plugin instance
func NewTreeInserterPlugin() resmap.TransformerPlugin {
	return &TreeInserter{}
}

func renderString(in string, rm resmap.ResMap) (string, error) {
	tmpl, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"getTreeByPath": getTreeByPath,
		},
	).Parse(in)
	if err != nil {
		return "", nil
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, rm)
	if err != nil {
		return "", nil
	}
	return buf.String(), nil

}

func getTreeByPath(resourceGVK, name, path string, rm resmap.ResMap) (string, error) {
	filterResID := resid.NewResId(gvk.FromString(resourceGVK), name)
	res, err := rm.GetById(filterResID)
	if err != nil {
		return "", err
	}
	tree, err := res.GetFieldValue(path)
	if err != nil {
		return "", err
	}
	result, err := yaml.Marshal(tree)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
