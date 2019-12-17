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

// ModifyHashStrings goes down the map recursively and tries
// to apply function to each string
func ModifyHashStrings(
	doc interface{},
	fn func(string) (string, error),
) (interface{}, error) {
	var err error
	switch typeV := doc.(type) {
	case map[string]interface{}:
		for k, v := range typeV {
			typeV[k], err = ModifyHashStrings(v, fn)
			if err != nil {
				return nil, err
			}
		}
	case []interface{}:
		for i, v := range typeV {
			typeV[i], err = ModifyHashStrings(v, fn)
			if err != nil {
				return nil, err
			}
		}
	case string:
		return fn(typeV)
	}
	return doc, nil
}
