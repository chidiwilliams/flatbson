// Package flatbson provides a function for recursively flattening a Go struct by its BSON tags.

package flatbson

import (
	"errors"
	"fmt"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

// Flatten returns a map with keys and values corresponding to the field name
// and values of struct v and its nested structs according to its BSON tags.
// It iterates over each field recursively and sets fields that are not nil.
//
// The BSON struct tags behave in line with the bsoncodec specification. See:
// https://godoc.org/go.mongodb.org/mongo-driver/bson/bsoncodec#StructTags
// for definitions. The supported tags are name, skip, omitempty, and inline.
//
// Flatten returns an error if v is not a struct or a pointer to a struct, or
// if the tags produce duplicate keys.
//
//     type A struct {
//       B *X `bson:"b,omitempty"`
//       C X  `bson:"c"`
//     }
//
//     type X struct { Y string `bson:"y"` }
//
//     Flatten(A{nil, X{"hello"}}) returns map[string]interface{}{"c.y": "hello"}
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

		if s, ok := asStruct(field); ok {
			fp := p
			if !tags.Inline {
				fp = p + tags.Name + "."
			}
			if err := flattenFields(s, m, fp); err != nil {
				return err
			}
			continue
		}

		key := p + tags.Name
		if _, ok := m[key]; ok {
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
