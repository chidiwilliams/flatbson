# flatbson

flatbson recursively flattens a Go struct by its BSON tags.

See [Godoc]() for full documentation.

## Example

```go
type Parent struct {
	B *Child `bson:"b,omitempty"`
	C Child  `bson:"c"`
}

type Child struct {
	Y string `bson:"y"`
}

flatbson.Flatten(Parent{nil, Child{"hello"}})

// Result:
// map[string]interface{}{"c.y": "hello"}
```

## Installation

```shell script
go get https://github.com/chidiwilliams/flatbson
```
