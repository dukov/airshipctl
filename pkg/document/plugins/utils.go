package plugins

import (
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
)

// ModifyHashStrings goes down the map recursively and tries
// to apply function to each string
func ModifyHashStrings(
	doc interface{},
	rm resmap.ResMap,
	fn func(string, resmap.ResMap) (string, error),
) (interface{}, error) {
	var err error
	switch typeV := doc.(type) {
	case map[string]interface{}:
		for k, v := range typeV {
			typeV[k], err = ModifyHashStrings(v, rm, fn)
			if err != nil {
				return nil, err
			}
		}
	case []interface{}:
		for i, v := range typeV {
			typeV[i], err = ModifyHashStrings(v, rm, fn)
			if err != nil {
				return nil, err
			}
		}
	case string:
		// TODO (dukov) Make this more generic to replace string with map or slice
		// instead of just go-template rendering
		return fn(typeV, rm)
	}
	return doc, nil
}
