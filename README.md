# flatbson

flatbson recursively flattens a Go struct using its BSON tags.

See [Godoc](https://pkg.go.dev/github.com/chidiwilliams/flatbson) for full documentation.

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
