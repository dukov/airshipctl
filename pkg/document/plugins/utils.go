package plugins

import (
	"bytes"
	"text/template"

	"sigs.k8s.io/kustomize/v3/pkg/resmap"
)

func ModifyHashStrings(doc interface{}, rm resmap.ResMap) (interface{}, error) {
	var err error
	switch typeV := doc.(type) {
	case map[string]interface{}:
		for k, v := range typeV {
			typeV[k], err = ModifyHashStrings(v, rm)
			if err != nil {
				return nil, err
			}
		}
	case []interface{}:
		for i, v := range typeV {
			typeV[i], err = ModifyHashStrings(v, rm)
			if err != nil {
				return nil, err
			}
		}
	case string:
		tmpl, err := template.New("tmpl").Funcs(
			template.FuncMap{
				"getTreeByPath": getTreeByPath,
			},
		).Parse(typeV)
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
	return doc, nil
}
