package redis

import (
	"fmt"
	"strings"

	"github.com/sean9999/go-store/essence"
)

// collection implements essence.collection.
// It provides functionality that more specialized collection types can inherit.
type collection struct {
	kind  string
	name  string
	store *Store
}

// a fully-qualified namespace for the collection
func (c *collection) Keyspace() string {
	return fmt.Sprintf("%s:%s:%s", c.store.Namespace, c.Kind(), c.Name())
}

func LongToShort(long string) string {
	bits := strings.Split(long, ":")
	return bits[len(bits)-1]
}

// Name is a namespace, or a table name
func (c *collection) Name() string {
	return c.name
}

// Kind is always "kv" for this kind of [Collection]
func (c *collection) Kind() string {
	return c.kind
}

// Store is a reference to the backing store (redis server)
func (c *collection) Store() essence.Store {
	return c.store
}
