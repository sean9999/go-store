package red

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sean9999/go-store/essence"
)

type CollectionMap struct {
	kv   map[string]essence.KeyValueCollection
	list map[string]essence.ListCollection
}

// a Store is a store of persistant data, structured into Collections.
// Store implements essence.Store.
type Store struct {
	Namespace     string
	client        *redis.Client
	collectionMap CollectionMap
}

func Attach(ns string, opts *redis.Options) *Store {

	ctx := context.Background()
	colMap := CollectionMap{
		kv:   map[string]essence.KeyValueCollection{},
		list: map[string]essence.ListCollection{},
	}
	s := Store{
		Namespace:     ns,
		client:        redis.NewClient(opts),
		collectionMap: colMap,
	}

	//	remember collections
	iter := s.client.Scan(ctx, 0, fmt.Sprintf("%s:*", s.schemaPrefix()), 0).Iterator()
	for iter.Next(ctx) {
		_, kind, colName := DeconstructSchemaKey(iter.Val())
		switch kind {
		case "kv":
			s.rememberCollection(ctx, "kv", colName)
		case "list":
			s.rememberCollection(ctx, "list", colName)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return &s
}

// Destroy the entire Store (database namespace) and everything in it
func (s *Store) Destroy(ctx context.Context) error {
	for _, col := range s.KeyValueCollections(ctx) {
		err := col.Destroy(ctx)
		if err != nil {
			return err
		}
	}
	for _, col := range s.ListCollections(ctx) {
		err := col.Destroy(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) rememberCollection(ctx context.Context, kind, name string) error {
	switch kind {
	case "kv":
		_, err := s.KeyValueCollection(ctx, name)
		if err != nil {
			return err
		}
		return nil
	case "list":
		_, err := s.ListCollection(ctx, name)
		if err != nil {
			return err
		}
		return nil
	default:
		return errors.New("unrecognized kind " + kind)
	}
}

// a LongKey is a fully-qualified key, taking into account the Store's namespace and collection's name
// func (s *Store) LongKey(col essence.Collection, lastBit string) string {
// 	return fmt.Sprintf("%s:%s", col.Keyspace(), lastBit)
// }

// ensure that a collection has a schema-key
func (s *Store) registerCollection(ctx context.Context, col essence.Collection) {
	schemaKey := ConstructSchemaKeyForCollection(s, col)
	s.client.Incr(ctx, schemaKey)
}

// remove the schema-key
func (s *Store) deregisterCollection(ctx context.Context, col essence.Collection) {
	schemaKey := ConstructSchemaKeyForCollection(s, col)
	s.client.Del(ctx, schemaKey)
}

// add a collection the Store's Collection map so it can be accessed
func (s *Store) addCollection(ctx context.Context, col essence.Collection) {
	switch col.Kind() {
	case "kv":
		s.collectionMap.kv[col.Name()] = col.(*KeyValueCollection)
	case "list":
		s.collectionMap.list[col.Name()] = col.(*ListCollection)
	}
}

// remove the collection from the Store's Collection map
func (s *Store) RemoveCollection(ctx context.Context, col essence.Collection) {

	switch col.Kind() {
	case "kv":
		//	delete the collection
		delete(s.collectionMap.kv, col.Name())
		//	delete the schema
		s.deregisterCollection(ctx, col)
	case "list":
		delete(s.collectionMap.list, col.Name())
		s.deregisterCollection(ctx, col)
	}

}

func (s *Store) CollectionExists(ctx context.Context, kind, name string) bool {

	switch kind {
	case "kv":
		_, exists := s.collectionMap.kv[name]
		return exists
	case "list":
		_, exists := s.collectionMap.list[name]
		return exists
	}
	return false

}

func (s *Store) ListCollections(ctx context.Context) map[string]essence.ListCollection {
	return s.collectionMap.list
}

func (s *Store) KeyValueCollections(ctx context.Context) map[string]essence.KeyValueCollection {
	return s.collectionMap.kv
}
