package flatbson_test

import (
	"reflect"
	"testing"

	"flatbson"
)

func TestFlatten(t *testing.T) {
	tests := []struct {
		name string
		v    interface{}
		want map[string]interface{}
	}{
		{
			name: "no nested fields",
			v: struct {
				A string `bson:"a,omitempty"`
			}{"az"},
			want: map[string]interface{}{"a": "az"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := flatbson.Flatten(tt.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Flatten() = %v, want %v", got, tt.want)
			}
		})
	}
}
