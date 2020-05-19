package flatbson

import (
	"errors"
	"reflect"

	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

// Flatten flattens a nested BSON struct.
func Flatten(v interface{}) (map[string]interface{}, error) {
	t := reflect.TypeOf(v)
	elem := reflect.ValueOf(v)

	if t.Kind() != reflect.Struct {
		if t.Kind() != reflect.Ptr {
			return nil, errors.New("v must be a struct or a pointer to a struct")
		}

		elem = elem.Elem()
	}

	m := make(map[string]interface{})

	for i := 0; i < elem.NumField(); i++ {
		ft := t.Field(i)

		tags, err := bsoncodec.DefaultStructTagParser(ft)
		if err != nil {
			return nil, err
		}

		m[tags.Name] = elem.Field(i).Interface()
	}

	return m, nil
}
