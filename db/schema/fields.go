package schema

import (
	"reflect"
	"slices"
	"strings"

	"github.com/doug-martin/goqu/v9"
)

func GetFields(data interface{}) []interface{} {
	selects := []interface{}{}
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {
		f, _ := t.FieldByName(t.Field(i).Name)
		dbTag := f.Tag.Get("db")
		if dbTag == "-" {
			continue
		}
		if dbTag != "" {
			asTag := f.Tag.Get("as")
			if asTag == "" {
				selects = append(selects, goqu.C(dbTag))
			} else {
				selects = append(selects, goqu.C(dbTag).As(asTag))
			}
			continue
		}
		jsonTag := f.Tag.Get("json")
		if jsonTag == "" {
			continue
		} else {
			selects = append(selects, goqu.C(jsonTag))
			continue
		}
	}
	return selects
}

func GetFieldsExceptFor(data interface{}, except []string) []interface{} {
	selects := []interface{}{}
	t := reflect.TypeOf(data)
	for i := 0; i < t.NumField(); i++ {
		f, _ := t.FieldByName(t.Field(i).Name)
		dbTag := f.Tag.Get("db")
		if dbTag == "-" {
			continue
		}
		if dbTag != "" {
			if slices.Contains(except, dbTag) {
				continue
			} else {
				asTag := f.Tag.Get("as")
				if asTag == "" {
					selects = append(selects, goqu.C(dbTag))
				} else {
					selects = append(selects, goqu.C(dbTag).As(asTag))
				}
			}
		} else {
			jsonTag := f.Tag.Get("json")
			if jsonTag == "" {
				continue
			} else {
				if strings.Contains(jsonTag, ",") {
					jsonTag = strings.Split(jsonTag, ",")[0]
					if slices.Contains(except, jsonTag) {
						continue
					}
					selects = append(selects, goqu.C(jsonTag))
					continue
				} else {
					if slices.Contains(except, jsonTag) {
						continue
					}
					selects = append(selects, goqu.C(jsonTag))
					continue
				}
			}
		}
	}
	return selects
}

func CheckUpdateFieldsExist(input any, allowedFields []string) (bool, []string) {
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

		// 허용되지 않은 필드이면서 값이 기본값이 아니면 invalid
		if !allowed[field] && !IsZeroValue(value) {
			invalid = append(invalid, field)
		}
	}

	return len(invalid) > 0, invalid
}

// 기본값 체크
func IsZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
