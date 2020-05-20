# flatbson

[![API reference](https://img.shields.io/badge/godoc-reference-5272B4)](https://pkg.go.dev/github.com/chidiwilliams/flatbson?tab=doc)

flatbson recursively flattens a Go struct using its BSON tags.

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

See the [Godoc](https://pkg.go.dev/github.com/chidiwilliams/flatbson) for the complete documentation.

## Installation

```shell script
go get https://github.com/chidiwilliams/flatbson
```
