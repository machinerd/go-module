package validation

import "reflect"

func ValidateUpdateFields(input any, allowedFields []string) (bool, []string) {
	allowed := make(map[string]bool)
	for _, f := range allowedFields {
		allowed[f] = true
	}

	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	t := v.Type()
	var invalid []string
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i).Name
		value := v.Field(i)

		if !allowed[field] && !IsZeroValue(value) {
			invalid = append(invalid, field)
		}
	}

	return len(invalid) > 0, invalid
}

func IsZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
