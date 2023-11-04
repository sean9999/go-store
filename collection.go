package store

import (
	"context"
	"fmt"
	"strings"
)

// a Collection is like a table. A discrete and namespaced area for data of a particular Kind
type Collection interface {
	Kind() string
	Name() string
	Store() *Store
	Forget(context.Context)
	Keyspace() string
}

// a GenericCollection is a simple struct that implements Collection
// and provides functionality that more specialized collections types can inherit.
type GenericCollection struct {
	name  string
	kind  string
	store *Store
}

// Name is a namespace, or a table name
func (c *GenericCollection) Name() string {
	return c.name
}

// Kind is always "kv" for this kind of [Collection]
func (c *GenericCollection) Kind() string {
	return c.kind
}

// Store is a reference to the backing store (redis server)
func (c *GenericCollection) Store() *Store {
	return c.store
}

func (c *GenericCollection) Keyspace() string {
	return fmt.Sprintf("%s:%s:%s", c.store.Namespace, c.Kind(), c.Name())
}

// leave the data, but lose the ability to interact with it.
func (c *GenericCollection) Forget(ctx context.Context) {
	store := c.store
	store.RemoveCollection(c)
	store.DeregisterCollection(ctx, c)
}

func LongToShort(long string) string {
	bits := strings.Split(long, ":")
	return bits[len(bits)-1]
}
