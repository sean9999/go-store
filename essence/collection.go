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
	Get(context.Context, string) (any, error)
	Set(context.Context, string, any) error
	Keys(context.Context) []string
	All(context.Context) map[string]any
	Destroy(context.Context) error
}

// a ListCollection has methods allowing it to be operated on as a list
type ListCollection interface {
	Collection
	Pop(context.Context) (any, error)
	Push(context.Context, any) error
	Shift(context.Context) (any, error)
	Unshift(context.Context, any) error
	All(context.Context) []any
	Size(context.Context) int
	Head(context.Context) (any, error)
	Tail(context.Context) (any, error)
	Destroy(context.Context) error
}
