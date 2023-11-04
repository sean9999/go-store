package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const (
	schemaPrefix = ".store:schema"
)

// a Store is a store of persistant data, structured into Collections
type Store struct {
	Namespace     string
	Client        *redis.Client
	collectionMap map[string]map[string]Collection
}

func NewStore(ns string) Store {

	emptyCollections := map[string]Collection{}
	colMap := map[string]map[string]Collection{
		"kv": emptyCollections,
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})

	s := Store{
		Namespace:     ns,
		Client:        rdb,
		collectionMap: colMap,
	}

	//	remember collections
	ctx := context.Background()
	iter := s.Client.Scan(ctx, 0, fmt.Sprintf("%s:%s:*", schemaPrefix, s.Namespace), 0).Iterator()
	for iter.Next(ctx) {
		_, kind, colName := DeconstructSchemaKey(iter.Val())
		switch kind {
		case "kv":
			col := NewKeyValueCollection(&s, colName)
			s.AddCollection(col)
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}
	return s
}

func (s Store) GetCollection(kind, name string) (Collection, error) {
	m, exists := s.collectionMap[kind]
	if !exists {
		return nil, fmt.Errorf("there are no collections of type %q", kind)
	}
	col, exists := m[name]
	if !exists {
		return nil, fmt.Errorf("no such collection %q", name)
	}
	return col, nil
}

// a LongKey is a fully-qualified key, taking into account the Store's namespace and collection's name
func (s Store) LongKey(col Collection, lastBit string) string {
	return fmt.Sprintf("%s:%s", col.Keyspace(), lastBit)
}

// ensure that a collection has a schema-key
func (s Store) RegisterCollection(ctx context.Context, col Collection) {
	s.Client.Incr(ctx, s.SchemaKey(col))
}

// remove the schema-key
func (s Store) DeregisterCollection(ctx context.Context, col Collection) {
	s.Client.Del(ctx, s.SchemaKey(col))
}

// add a collection the Store's Collection map so it can be accessed
func (s Store) AddCollection(col Collection) {
	s.collectionMap[col.Kind()][col.Name()] = col
}

// remove the collection from the Store's Collection map
func (s Store) RemoveCollection(col Collection) {
	delete(s.collectionMap[col.Kind()], col.Name())
}

func (s Store) CollectionExists(kind, name string) bool {
	schemaKey := ConstructSchemaKey(s.Namespace, kind, name)
	ctx := context.Background()
	_, err := s.Client.Get(ctx, schemaKey).Result()
	if err == redis.Nil {
		return false
	} else {
		return true
	}
}

func (s Store) DeleteCollection(col Collection) error {
	if s.CollectionExists(col.Kind(), col.Name()) {
		ctx := context.Background()
		delete(s.collectionMap[col.Kind()], col.Name())
		s.DeregisterCollection(ctx, col)
		return nil
	}
	return errors.New("collection doesn't exist")
}
