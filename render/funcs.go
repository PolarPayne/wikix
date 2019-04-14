package render

import (
	"errors"
	"fmt"
)

func dict(args ...interface{}) (map[string]interface{}, error) {
	if len(args)%2 != 0 {
		return nil, errors.New("dict must be called with even number of arguments")
	}

	out := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		k, v := args[i], args[i+1]

		kString, ok := k.(string)
		if !ok {
			return nil, fmt.Errorf("%v is not a valid key for dictionary", k)
		}

		out[kString] = v
	}

	return out, nil
}

func list(args ...interface{}) []interface{} {
	out := make([]interface{}, len(args))
	for i, v := range args {
		out[i] = v
	}
	return out
}
