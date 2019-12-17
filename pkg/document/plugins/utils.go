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
		// TODO (dukov) Make this more generic to replace string with map or slice
		// instead of just go-template rendering
		return fn(typeV)
	}
	return doc, nil
}
