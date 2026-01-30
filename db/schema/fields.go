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
