package plugins

import (
	"bytes"
	"fmt"
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
	renderPath  map[string]bool
	currentGvkn string
	rp          []string
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
		resource := r.Map()
		rootResource := r.CurId().GvknString()
		renderFunc := func(in string) (string, error) {
			return gt.renderString(in, map[string]bool{rootResource: true})
		}
		_, err := ModifyHashStrings(resource, renderFunc)
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

func (gt *GoTemplater) renderString(in string, renderPath map[string]bool) (string, error) {
	getByPathWrap := func(rGVK, name, path string) (interface{}, error) {
		return gt.getTreeByPath(rGVK, name, path, renderPath)
	}
	tmpl, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"getTreeByPath": getByPathWrap,
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

func (gt *GoTemplater) getTreeByPath(resourceGVK, name, path string, renderPath map[string]bool) (interface{}, error) {
	filterResID := resid.NewResId(gvk.FromString(resourceGVK), name)
	renderPathEntry := fmt.Sprintf("%s|%s", filterResID.GvknString(), path)
	if _, foud := renderPath[renderPathEntry]; foud {
		return nil, RenderLoopError{Resource: filterResID.GvknString()}
	}

	res, err := gt.resourceMap.GetById(filterResID)
	if err != nil {
		return "", err
	}
	tgtField, err := res.GetFieldValue(path)
	if err != nil {
		return "", err
	}
	switch tgtType := tgtField.(type) {
	case map[string]interface{}, []interface{}:
		result, err := yaml.Marshal(tgtType)
		if err != nil {
			return "", err
		}
		renderPath = appendToPath(renderPath, renderPathEntry)
		return gt.renderString(string(result), renderPath)

	case string:
		renderPath = appendToPath(renderPath, renderPathEntry)
		return gt.renderString(tgtType, renderPath)
	}
	return tgtField, nil

}

func appendToPath(src map[string]bool, entry string) (res map[string]bool) {
	res = make(map[string]bool)
	for k, v := range src {
		res[k] = v
	}
	res[entry] = true
	return
}
