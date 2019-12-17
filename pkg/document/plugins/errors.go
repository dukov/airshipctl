package plugins

import "fmt"

// UnknownPluginError raised for unkregistered plugin kinds
type UnknownPluginError struct {
	Kind string
}

func (e UnknownPluginError) Error() string {
	return fmt.Sprintf("Unknown airship plugin with Kind: %s", e.Kind)
}
