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

import "fmt"

// ErrUnknownPlugin raised for unregistered plugins
type ErrUnknownPlugin struct {
	Kind string
}

func (e ErrUnknownPlugin) Error() string {
	return fmt.Sprintf("Unknown airship plugin with Kind: %s", e.Kind)
}

// ErrRenderLoop is raised in case of template rendering loop
type ErrRenderLoop struct {
	Resource string
}

func (e ErrRenderLoop) Error() string {
	return fmt.Sprintf("Render loop detected around %s", e.Resource)
}