package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
)

const (
	Insert  = "INSERT"
	Update  = "UPDATE"
	Delete  = "DELETE"
	Reorder = "REORDER"
)

type Change struct {
	FieldName  string
	ActionType string
	BeforeData *string
	AfterData  *string
}

type Options struct {
	IgnoredKeys []string
	IDKey       string
	OrderKey    string
	Identity    func(path string, item map[string]any) string
}

type DiffInput struct {
	Before  any
	After   any
	Options Options
}

func normalize(value any) (any, error) {
	data, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var result any
	err = decoder.Decode(&result)
	return result, err
}

func auditEmpty(v any) bool {
	if v == nil {
		return true
	}
	switch v := v.(type) {
	case string:
		return v == ""
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	}
	return false
}

func auditEqual(before, after any) bool {
	return (auditEmpty(before) && auditEmpty(after)) || reflect.DeepEqual(before, after)
}

func auditText(v any) *string {
	if v == nil {
		return nil
	}
	if s, ok := v.(string); ok {
		return &s
	}
	data, _ := json.Marshal(v)
	s := string(data)
	return &s
}

func Diff(input DiffInput) ([]Change, error) {
	before, err := normalize(input.Before)
	if err != nil {
		return nil, fmt.Errorf("audit before: %w", err)
	}
	after, err := normalize(input.After)
	if err != nil {
		return nil, fmt.Errorf("audit after: %w", err)
	}
	ignored := map[string]bool{"ref": true, "createdAt": true, "updatedAt": true, "created_at": true, "updated_at": true}
	for _, key := range input.Options.IgnoredKeys {
		ignored[key] = true
	}
	idKey := input.Options.IDKey
	if idKey == "" {
		idKey = "id"
	}
	orderKey := input.Options.OrderKey
	if orderKey == "" {
		orderKey = "priority"
	}
	var logs []Change
	var diffErr error
	emit := func(path, action string, old, current any) {
		if auditEqual(old, current) {
			return
		}
		logs = append(logs, Change{FieldName: path, ActionType: action, BeforeData: auditText(old), AfterData: auditText(current)})
	}

	var walk func(string, string, any, any)
	walk = func(path, action string, old, current any) {
		if diffErr != nil {
			return
		}
		if auditEqual(old, current) {
			return
		}
		oldMap, oldObject := old.(map[string]any)
		newMap, newObject := current.(map[string]any)
		if oldObject || newObject {
			keys := map[string]bool{}
			for key := range oldMap {
				keys[key] = true
			}
			for key := range newMap {
				keys[key] = true
			}
			names := make([]string, 0, len(keys))
			for key := range keys {
				names = append(names, key)
			}
			sort.Strings(names)
			for _, key := range names {
				if ignored[key] {
					continue
				}
				child := key
				if path != "" {
					child = path + "." + key
				}
				walk(child, action, oldMap[key], newMap[key])
			}
			return
		}
		oldList, oldArray := old.([]any)
		newList, newArray := current.([]any)
		if (oldArray || newArray) && auditObjectList(idKey, oldList, newList) {
			identity := func(item any) string {
				m := item.(map[string]any)
				if input.Options.Identity != nil {
					if key := input.Options.Identity(path, m); key != "" {
						return key
					}
				}
				return fmt.Sprint(m[idKey])
			}

			oldByID, newByID := map[string]any{}, map[string]any{}
			var oldIDs, newIDs []any
			for _, item := range oldList {
				if _, exists := oldByID[identity(item)]; exists {
					diffErr = fmt.Errorf("audit %s: duplicate before identity %q", path, identity(item))
					return
				}
				oldByID[identity(item)] = item
				oldIDs = append(oldIDs, item.(map[string]any)[idKey])
			}
			for _, item := range newList {
				if _, exists := newByID[identity(item)]; exists {
					diffErr = fmt.Errorf("audit %s: duplicate after identity %q", path, identity(item))
					return
				}
				newByID[identity(item)] = item
				newIDs = append(newIDs, item.(map[string]any)[idKey])
			}
			var oldCommon, newCommon []string
			for _, item := range oldList {
				id := identity(item)
				if _, ok := newByID[id]; ok {
					oldCommon = append(oldCommon, id)
				}
			}
			for _, item := range newList {
				id := identity(item)
				if _, ok := oldByID[id]; ok {
					newCommon = append(newCommon, id)
				}
			}
			reordered := !reflect.DeepEqual(oldCommon, newCommon)
			if reordered {
				emit(path, "REORDER", oldIDs, newIDs)
			}
			for i, item := range newList {
				previous, exists := oldByID[identity(item)]
				itemAction := action
				if !exists {
					itemAction = "INSERT"
				}
				if reordered && exists {
					previous = auditWithoutOrder(previous, orderKey)
					item = auditWithoutOrder(item, orderKey)
				}
				walk(indexPath(path, i), itemAction, previous, item)
			}
			for i, item := range oldList {
				if _, exists := newByID[identity(item)]; !exists {
					walk(indexPath(path, i), "DELETE", item, nil)
				}
			}
			return
		}
		emit(path, action, old, current)
	}
	walk("", "UPDATE", before, after)
	if diffErr != nil {
		return nil, diffErr
	}
	return logs, nil
}

func auditObjectList(idKey string, lists ...[]any) bool {
	found := false
	for _, list := range lists {
		for _, item := range list {
			m, ok := item.(map[string]any)
			if !ok || m[idKey] == nil {
				return false
			}
			found = true
		}
	}
	return found
}

func auditWithoutOrder(value any, orderKey string) map[string]any {
	result := map[string]any{}
	for key, value := range value.(map[string]any) {
		if key != orderKey {
			result[key] = value
		}
	}
	return result
}

func indexPath(path string, index int) string {
	if path == "" {
		return fmt.Sprint(index)
	}
	return fmt.Sprintf("%s.%d", path, index)
}
