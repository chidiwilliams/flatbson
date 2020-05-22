# FlatBSON

[![Build status](https://github.com/chidiwilliams/flatbson/workflows/Build/badge.svg)](https://github.com/chidiwilliams/flatbson/actions?query=workflow%3ABuild) [![API reference](https://img.shields.io/badge/godoc-reference-5272B4)](https://pkg.go.dev/github.com/chidiwilliams/flatbson?tab=doc) [![codecov](https://codecov.io/gh/chidiwilliams/flatbson/branch/master/graph/badge.svg)](https://codecov.io/gh/chidiwilliams/flatbson)

FlatBSON recursively flattens a Go struct using its BSON tags.

It is particularly useful for partially updating embedded Mongo documents.

For example, to update a `User`'s `Address.Visited` field, first call `flatbson.Flatten` with the parent struct:

```go
type User struct {
  ID      bson.ObjectID `bson:"_id,omitempty"`
  Name    string        `bson:"name,omitempty"`
  Address Address       `bson:"address,omitempty"`
}

type Address struct {
  Street    string    `bson:"street,omitempty"`
  City      string    `bson:"city,omitempty"`
  State     string    `bson:"state,omitempty"`
  VisitedAt time.Time `bson:"visitedAt,omitempty"`
}

flatbson.Flatten(User{Address: {VisitedAt: time.Now().UTC()}})

// Result:
// map[string]interface{}{"address.visitedAt": time.Time{...}}
```

Passing the result to `coll.UpdateOne` updates only the `address.VisitedAt` field instead of overwriting the entire `address` embedded document. See this [blog post](https://dev.to/chidiwilliams/partially-updating-an-embedded-mongo-document-in-go-knn) for more information.

The complete documentation is available on [Godoc](https://pkg.go.dev/github.com/chidiwilliams/flatbson).

## How to Install

```shell script
go get https://github.com/chidiwilliams/flatbson
```
