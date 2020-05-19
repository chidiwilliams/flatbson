package flatbson_test

import (
	"errors"
	"reflect"
	"testing"

	"flatbson"
)

type root struct {
	A string `bson:"a"`
	B int    `bson:"b"`
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
			want: nil,
			err:  errors.New("v must be a struct or a pointer to a struct"),
		},
		{
			name: "root fields",
			args: args{root{"az", 5}},
			want: map[string]interface{}{"a": "az", "b": 5},
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
