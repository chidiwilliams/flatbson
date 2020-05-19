package flatbson

import (
	"errors"
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
	flattenFields(val, m, "")
	return m, nil
}

func flattenFields(v reflect.Value, m map[string]interface{}, p string) {
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
			flattenFields(s, m, p+tags.Name+".")
			continue
		}

		m[p+tags.Name] = field.Interface()
	}
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
