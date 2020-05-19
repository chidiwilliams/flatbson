package flatbson

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

// Flatten flattens a nested BSON struct.
func Flatten(v interface{}) (map[string]interface{}, error) {
	val := reflect.ValueOf(v)

	val, ok := asStruct(val)
	if !ok {
		return nil, errors.New("v must be a struct or a pointer to a struct")
	}

	m := make(map[string]interface{})
	if err := flattenFields(val, m, ""); err != nil {
		return nil, err
	}

	return m, nil
}

func flattenFields(v reflect.Value, m map[string]interface{}, p string) error {
	for i := 0; i < v.NumField(); i++ {
		tags, _ := bsoncodec.DefaultStructTagParser(v.Type().Field(i))

		if tags.Skip {
			continue
		}

		field := v.Field(i)
		if tags.OmitEmpty && field.IsZero() {
			continue
		}

		s, ok := asStruct(field)
		if ok {
			if err := flattenFields(s, m, p+tags.Name+"."); err != nil {
				return err
			}
			continue
		}

		key := p + tags.Name

		if _, ok = m[key]; ok {
			return fmt.Errorf("duplicated key %s", key)
		}

		m[key] = field.Interface()
	}

	return nil
}

func asStruct(v reflect.Value) (reflect.Value, bool) {
	for {
		switch v.Kind() {
		case reflect.Struct:
			return v, true
		case reflect.Ptr:
			v = v.Elem()
		default:
			return reflect.Value{}, false
		}
	}
}
