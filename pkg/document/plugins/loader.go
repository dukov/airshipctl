package plugins

import (
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/yaml"
)

// PluginRegistry map of plugin kinds to plugin instance
var PluginRegistry = make(map[string]resmap.TransformerPlugin)

// TransformerLoader airship document plugin loader. Loads
type TransformerLoader struct {
	resid.ResId
}

// Config reads plugin configuration structure
func (l *TransformerLoader) Config(
	ldr ifc.Loader, rf *resmap.Factory, c []byte) error {
	if err := yaml.Unmarshal(c, l); err != nil {
		return err
	}
	airshipPlugin, found := PluginRegistry[l.Kind]
	if !found {
		return UnknownPluginError{Kind: l.Kind}
	}
	cfg := airshipPlugin.(resmap.Configurable)
	return cfg.Config(ldr, rf, c)
}

// Transform run plugin's Transorm method
func (l *TransformerLoader) Transform(m resmap.ResMap) error {
	airshipPlugin, found := PluginRegistry[l.Kind]
	if !found {
		return UnknownPluginError{Kind: l.Kind}
	}
	transformer := airshipPlugin.(resmap.Transformer)
	return transformer.Transform(m)
}

// NewTransformerLoader returns aitship document plugin transformer
func NewTransformerLoader() resmap.TransformerPlugin {
	return &TransformerLoader{}
}
