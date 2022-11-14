package flatbson_test

import (
	"errors"
	"reflect"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"

	"github.com/8arcturus/flatbson"
)

type root struct {
	A string `bson:"a"`
	B int    `bson:"b"`
}

type rootWithExportedField struct {
	A string `bson:"a"`
	B int    `bson:"b"`
	c int    `bson:"c"`
}

type nestedRoot struct {
	A nestedLeaf   `bson:"a"`
	B nestedBranch `bson:"b"`
}

type nestedBranch struct {
	C nestedLeaf `bson:"c"`
}

type nestedLeaf struct {
	B int `bson:"b"`
}

type nestedRootPtr struct {
	A *nestedLeaf `bson:"a"`
}

type skipRoot struct {
	A int      `bson:"-"`
	B skipLeaf `bson:"b"`
}

type skipLeaf struct {
	C string `bson:"-"`
	D int    `bson:"d"`
}

type omitemptyRoot struct {
	A int           `bson:"a,omitempty"`
	B int           `bson:"b"`
	C omitemptyLeaf `bson:"c"`
}

type omitemptyLeaf struct {
	A string      `bson:"a,omitempty"`
	B interface{} `bson:"b,omitempty"`
	C []string    `bson:"c"`
}

type duplicateRoot struct {
	A string `bson:"a,omitempty"`
	B int    `bson:"a,omitempty"`
}

type duplicateNestedRoot struct {
	A string              `bson:"a"`
	B duplicateNestedLeaf `bson:""`
}

type duplicateNestedLeaf struct {
	A int `bson:"a"`
	B int `bson:"a"`
}

type inlineRoot struct {
	A inlineBranch `bson:"a,inline"`
	X string       `bson:"x"`
}

type inlineBranch struct {
	B inlineLeaf `bson:"b,inline"`
	Y int        `bson:"y"`
}

type inlineLeaf struct {
	C string   `bson:"c,inline"`
	Z []string `bson:"z"`
}

type unexportedRoot struct {
	A unexportedLeaf `bson:"a"`
}

type mixedExportedUnexportedRoot struct {
	A mixedExportedUnexportedLeaf `bson:"a"`
}
type unexportedLeaf struct {
	b string
}

type mixedExportedUnexportedLeaf struct {
	b string `bson:"b"`
	D string `bson:"d"`
}
type customBSONTypeRoot struct {
	A customBSONTypeLeaf `bson:"a"`
}

type customBSONTypeLeaf struct {
	S string
}

func (c customBSONTypeLeaf) MarshalBSONValue() (bsontype.Type, []byte, error) {
	return bson.MarshalValue(c.S)
}

func TestFlatten(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
		err  error
	}{
		{
			name: "non-struct input",
			args: args{v: 23},
			err:  errors.New("v must be a struct or a pointer to a struct"),
		},
		{
			name: "root fields",
			args: args{root{"az", 5}},
			want: map[string]interface{}{"a": "az", "b": 5},
		},
		{
			name: "duplicate keys",
			args: args{duplicateRoot{"as", 12}},
			err:  errors.New("duplicated key a"),
		},
		{
			name: "duplicate nested keys",
			args: args{duplicateNestedRoot{"as", duplicateNestedLeaf{1, 2}}},
			want: nil,
			err:  errors.New("duplicated key b.a"),
		},
		{
			name: "nested fields",
			args: args{nestedRoot{A: nestedLeaf{B: 5}, B: nestedBranch{C: nestedLeaf{B: 60}}}},
			want: map[string]interface{}{"a.b": 5, "b.c.b": 60},
		},
		{
			name: "nested fields with ptrs",
			args: args{v: nestedRootPtr{A: &nestedLeaf{B: 23}}},
			want: map[string]interface{}{"a.b": 23},
		},
		{
			name: "skip fields",
			args: args{v: skipRoot{1, skipLeaf{"23", 74}}},
			want: map[string]interface{}{"b.d": 74},
		},
		{
			name: "omitempty fields",
			args: args{v: omitemptyRoot{0, 0, omitemptyLeaf{"", nil, []string{}}}},
			want: map[string]interface{}{"b": 0, "c.c": []string{}},
		},
		{
			name: "inline fields",
			args: args{inlineRoot{inlineBranch{inlineLeaf{"abc", []string{"jd"}}, 34}, "rwr"}},
			want: map[string]interface{}{"c": "abc", "z": []string{"jd"}, "y": 34, "x": "rwr"},
		},
		{
			name: "unexported fields",
			args: args{unexportedRoot{unexportedLeaf{"abc"}}},
			want: map[string]interface{}{"a": unexportedLeaf{"abc"}},
		},
		{
			name: "fields that can marshal themselves",
			args: args{customBSONTypeRoot{customBSONTypeLeaf{"abc"}}},
			want: map[string]interface{}{"a": customBSONTypeLeaf{"abc"}},
		},
		{
			name: "with unexported field in root struct",
			args: args{rootWithExportedField{"az", 5, 10}},
			want: map[string]interface{}{"a": "az", "b": 5},
		},
		{
			name: "nested struct has unexported field",
			args: args{mixedExportedUnexportedRoot{mixedExportedUnexportedLeaf{"abc", "bc"}}},
			want: map[string]interface{}{"a.d": "bc"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := flatbson.Flatten(tt.args.v)
			if !reflect.DeepEqual(err, tt.err) {
				t.Errorf("Flatten() error = %v, want %v", err, tt.err)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flatten() got = %v, want %v", got, tt.want)
			}
		})
	}
}
