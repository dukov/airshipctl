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
	PluginRegistry["GoTemplater"] = NewGoTemplaterPlugin()
}

// GoTemplater plugin to insert map subtree from specified resource
type GoTemplater struct {
	resid.ResId
	resourceMap resmap.ResMap
}

// Config loads configuration from plugin config
func (gt *GoTemplater) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	return yaml.Unmarshal(c, gt)
}

// Transform applies transformation to all resources
func (gt *GoTemplater) Transform(m resmap.ResMap) error {
	gt.resourceMap = m
	for _, r := range m.Resources() {
		rm := r.Map()
		_, err := ModifyHashStrings(rm, gt.renderString)
		if err != nil {
			return err
		}
	}
	return nil
}

// NewGoTemplaterPlugin returns plugin instance
func NewGoTemplaterPlugin() resmap.TransformerPlugin {
	return &GoTemplater{}
}

func (gt *GoTemplater) renderString(in string) (string, error) {
	tmpl, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"getTreeByPath": gt.getTreeByPath,
		},
	).Parse(in)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, nil)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

func (gt *GoTemplater) getTreeByPath(resourceGVK, name, path string) (string, error) {
	filterResID := resid.NewResId(gvk.FromString(resourceGVK), name)
	res, err := gt.resourceMap.GetById(filterResID)
	if err != nil {
		return "", err
	}
	tree, err := res.GetFieldValue(path)
	if err != nil {
		return "", err
	}
	// TODO(dukov) check if value is go template and render it as well
	result, err := yaml.Marshal(tree)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
