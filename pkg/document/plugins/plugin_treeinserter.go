package plugins

import (
	"fmt"

	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/yaml"
)

func init() {
	PluginRegistry["TreeInserter"] = NewTreeInserterPLugin()
}

type TreeInserter struct {
	resid.ResId
}

func (ti *TreeInserter) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	//if err := yaml.Unmarshal(c, ti); err != nil {
	//return err
	//}
	return yaml.Unmarshal(c, ti)
}

func (ti *TreeInserter) Transform(m resmap.ResMap) error {
	fmt.Printf("DDD transform %#v\n", m)
	for _, r := range m.Resources() {
		rm := r.Map()
		//refl := reflect.ValueOf(&rm)
		//walk(refl)
		ModifyHashStrings(rm, m)
		fmt.Printf("DDD result %#v\n", rm)
	}
	return nil
}

func NewTreeInserterPLugin() resmap.TransformerPlugin {
	return &TreeInserter{}
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
