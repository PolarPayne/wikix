package render

import "log"

func dict(args ...interface{}) map[string]interface{} {
	if len(args)%2 != 0 {
		return nil
	}

	out := make(map[string]interface{})
	for i := 0; i < len(args); i += 2 {
		k, v := args[i], args[i+1]

		kString, ok := k.(string)
		if !ok {
			log.Printf("%v is not a valid key for dictionary", k)
			return nil
		}

		out[kString] = v
	}

	return out
}

func list(args ...interface{}) []interface{} {
	out := make([]interface{}, len(args))
	for i, v := range args {
		out[i] = v
	}
	return out
}
