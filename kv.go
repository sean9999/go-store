package store

import (
	"context"
	"fmt"
)

// a KeyValueCollection is a GenericCollection with methods suitable to a key-value store
type KeyValueCollection struct {
	GenericCollection
}

func (s Store) GetKeyValueCollection(name string) (*KeyValueCollection, error) {
	col, err := s.GetCollection("kv", name)
	if err != nil {
		return nil, err
	}
	return col.(*KeyValueCollection), nil
}

func NewKeyValueCollection(db *Store, name string) *KeyValueCollection {
	ctx := context.Background()
	c := GenericCollection{
		kind:  "kv",
		name:  name,
		store: db,
	}
	kv := &KeyValueCollection{c}
	db.RegisterCollection(ctx, kv)
	db.AddCollection(kv)
	return kv
}

func (kv *KeyValueCollection) Get(ctx context.Context, shortKey string) (string, error) {
	store := kv.Store()
	longKey := kv.store.LongKey(kv, shortKey)
	val, err := store.Client.Get(ctx, longKey).Result()
	return val, err
}

func (kv *KeyValueCollection) Set(ctx context.Context, shortKey string, val any) error {
	store := kv.Store()
	longKey := kv.store.LongKey(kv, shortKey)
	err := store.Client.Set(ctx, longKey, val, 0).Err()
	return err
}

func (kv *KeyValueCollection) Scan(ctx context.Context) []string {
	r := []string{}
	keyspace := fmt.Sprintf("%s:*", kv.Keyspace())
	iter := kv.store.Client.Scan(ctx, 0, keyspace, 0).Iterator()
	for iter.Next(ctx) {
		r = append(r, LongToShort(iter.Val()))
	}
	return r
}

func (kv *KeyValueCollection) GetAll(ctx context.Context) map[string]any {
	m := map[string]any{}
	store := kv.store
	shortKeys := kv.Scan(ctx)
	longKeys := make([]string, len(shortKeys))
	for i, k := range shortKeys {
		longKeys[i] = store.LongKey(kv, k)
	}
	objs := kv.store.Client.MGet(ctx, longKeys...).Val()
	for i, k := range shortKeys {
		m[k] = objs[i]
	}
	return m
}
