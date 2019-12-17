/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugins

import (
	"bytes"
	"fmt"
	"text/template"

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
}

// Config loads configuration from plugin config
func (gt *GoTemplater) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	return yaml.Unmarshal(c, gt)
}

// Transform applies transformation to all resources
func (gt *GoTemplater) Transform(m resmap.ResMap) error {
	for _, r := range m.Resources() {
		resource := r.Map()
		rootResource := r.CurId().GvknString()
		renderFunc := func(in string) (string, error) {
			return gt.renderString(in, map[string]bool{rootResource: true}, m)
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

func (gt *GoTemplater) renderString(in string, renderPath map[string]bool, m resmap.ResMap) (string, error) {
	getByPathWrap := func(rGVK, name, path string) (interface{}, error) {
		return gt.getTreeByPath(rGVK, name, path, renderPath, m)
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

func (gt *GoTemplater) getTreeByPath(
	resourceGVK, name, path string,
	renderPath map[string]bool,
	m resmap.ResMap,
) (interface{}, error) {
	filterResID := resid.NewResId(gvk.FromString(resourceGVK), name)
	renderPathEntry := fmt.Sprintf("%s|%s", filterResID.GvknString(), path)
	if _, foud := renderPath[renderPathEntry]; foud {
		return nil, ErrRenderLoop{Resource: filterResID.GvknString()}
	}

	res, err := m.GetById(filterResID)
	if err != nil {
		return "", err
	}
	tgtField, err := res.GetFieldValue(path)
	if err != nil {
		return "", err
	}
	switch tgtType := tgtField.(type) {
	case map[string]interface{}, []interface{}:
		renderPath = appendToPath(renderPath, renderPathEntry)
		renderF := func(in string) (string, error) {
			return gt.renderString(in, renderPath, m)
		}

		// NOTE Dig deeper recurcively to render array/map strings.
		// This is needed since yaml.Marshal splits long strings (at
		// least in case string contains curly braces) into multiple
		// lines which may cause go template parse errors
		_, err := ModifyHashStrings(tgtType, renderF)
		if err != nil {
			return "", err
		}

		result, err := yaml.Marshal(tgtType)
		if err != nil {
			return "", err
		}
		return gt.renderString(string(result), renderPath, m)

	case string:
		renderPath = appendToPath(renderPath, renderPathEntry)
		return gt.renderString(tgtType, renderPath, m)
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
