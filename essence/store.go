package essence

import "context"

// a store is a logical namespace that contains collections, which are themselves logical namespaces
type Store interface {
	KeyValueCollection(context.Context, string) (KeyValueCollection, error)
	ListCollection(context.Context, string) (ListCollection, error)
	RemoveCollection(context.Context, Collection)
	CollectionExists(context.Context, string, string) bool
	KeyValueCollections(context.Context) map[string]KeyValueCollection
	ListCollections(context.Context) map[string]ListCollection
	Destroy(context.Context) error
}
