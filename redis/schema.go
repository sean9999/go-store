package redis

import (
	"fmt"
	"strings"

	"github.com/sean9999/go-store/essence"
)

const (
	schemaPrefix = ".store:schema"
)

// a schema-key registers the existence of a Collection in a Store
func ConstructSchemaKey(namespace, kind, name string) string {
	return fmt.Sprintf("%s:%s:%s:%s", schemaPrefix, namespace, kind, name)
}

// the path where a schema can be found
func (s Store) SchemaKey(col essence.Collection) string {
	return ConstructSchemaKey(s.Namespace, col.Kind(), col.Name())
}

// parse out and return the important values in a schema-key
func DeconstructSchemaKey(schemaKey string) (namespace string, kind string, collectionName string) {
	slugs := strings.Split(schemaKey, ":")
	return slugs[2], slugs[3], slugs[4]
}
