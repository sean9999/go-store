package red

import (
	"fmt"
	"strings"

	"github.com/sean9999/go-store/essence"
)

func (s *Store) schemaPrefix() string {
	return fmt.Sprintf(".store:%s:schema", s.Namespace)
}

func (s *Store) dataPrefix() string {
	return fmt.Sprintf(".store:%s:data", s.Namespace)
}

func ConstructSchemaKeyForCollection(s *Store, col essence.Collection) string {
	return fmt.Sprintf("%s:%s:%s", s.schemaPrefix(), col.Kind(), col.Name())
}

// // a schema-key registers the existence of a Collection in a Store
// func ConstructSchemaKey(namespace, kind, name string) string {
// 	return fmt.Sprintf("%s:%s:%s:%s", schemaPrefix, namespace, kind, name)
// }

// // the path where a schema can be found
// func (s Store) SchemaKey(col essence.Collection) string {
// 	return ConstructSchemaKey(s.Namespace, col.Kind(), col.Name())
// }

// parse out and return the important values in a schema-key
func DeconstructSchemaKey(schemaKey string) (namespace string, kind string, collectionName string) {
	slugs := strings.Split(schemaKey, ":")

	//slugs[0] // .store
	//slugs[1] // namespace
	//slugs[2] // "schema" or "data"
	//slugs[3] // kind of collection (kv or list)
	//slugs[4] // collection name

	return slugs[1], slugs[3], slugs[4]
}
