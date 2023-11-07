package essence

import (
	"context"
)

// a Collection is like a table. A discrete and namespaced area for data of a particular Kind
type Collection interface {
	Kind() string
	Name() string
	Store() Store
	Keyspace() string
}

// a KeyValueCollection has methods allowing it to behave as a key-value store
type KeyValueCollection interface {
	Collection
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, []byte) error
	Keys(context.Context) []string
	All(context.Context) map[string][]byte
	Destroy(context.Context) error
}

// a ListCollection has methods allowing it to be operated on as a list
type ListCollection interface {
	Collection
	Pop(context.Context) ([]byte, error)
	Push(context.Context, []byte) error
	Shift(context.Context) ([]byte, error)
	Unshift(context.Context, []byte) error
	All(context.Context) [][]byte
	Size(context.Context) int
	Head(context.Context) ([]byte, error)
	Tail(context.Context) ([]byte, error)
	Destroy(context.Context) error
}
