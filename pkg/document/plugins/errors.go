package plugins

import "fmt"

// UnknownPluginError raised for unkregistered plugin kinds
type UnknownPluginError struct {
	Kind string
}

func (e UnknownPluginError) Error() string {
	return fmt.Sprintf("Unknown airship plugin with Kind: %s", e.Kind)
}

// RenderLoopError reaisd in case of template rendeting loop
type RenderLoopError struct {
	Resource string
}

func (e RenderLoopError) Error() string {
	return fmt.Sprintf("Render loop detected around %s", e.Resource)
}
