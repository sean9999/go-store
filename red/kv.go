package red

import (
	"context"
	"fmt"

	"github.com/sean9999/go-store/essence"
)

// a KeyValueCollection is a collection that can operate on key-value pairs
type KeyValueCollection struct {
	collection
}

func (kv *KeyValueCollection) Destroy(ctx context.Context) error {
	shortKeys := kv.Keys(ctx)
	longKeys := make([]string, len(shortKeys))
	for i, k := range shortKeys {
		longKey := fmt.Sprintf("%s:%s", kv.Keyspace(), k)
		longKeys[i] = longKey
	}
	//	delete the data from redis
	err := kv.store.client.Del(ctx, longKeys...).Err()
	if err == nil {
		//	delete the collection itself
		kv.store.RemoveCollection(ctx, kv)
	}
	return err
}

func (s *Store) KeyValueCollection(ctx context.Context, name string) (essence.KeyValueCollection, error) {
	kv, exists := s.collectionMap.kv[name]
	if !exists {
		c := collection{
			kind:  "kv",
			name:  name,
			store: s,
		}
		kv = &KeyValueCollection{c}
		s.addCollection(ctx, kv)
	}
	s.registerCollection(ctx, kv)
	return kv, nil
}

func (kv *KeyValueCollection) Get(ctx context.Context, shortKey string) ([]byte, error) {
	longKey := fmt.Sprintf("%s:%s", kv.Keyspace(), shortKey)
	val, err := kv.store.client.Get(ctx, longKey).Bytes()
	return val, err
}

func (kv *KeyValueCollection) Set(ctx context.Context, shortKey string, v []byte) error {
	longKey := fmt.Sprintf("%s:%s", kv.Keyspace(), shortKey)
	err := kv.store.client.Set(ctx, longKey, v, 0).Err()
	return err
}

func (kv *KeyValueCollection) Keys(ctx context.Context) []string {
	r := []string{}
	keyspace := fmt.Sprintf("%s:*", kv.Keyspace())
	iter := kv.store.client.Scan(ctx, 0, keyspace, 0).Iterator()
	for iter.Next(ctx) {
		r = append(r, LongToShort(iter.Val()))
	}
	return r
}

func (kv *KeyValueCollection) All(ctx context.Context) map[string][]byte {
	m := map[string][]byte{}
	shortKeys := kv.Keys(ctx)
	longKeys := make([]string, len(shortKeys))
	for i, k := range shortKeys {
		longKey := fmt.Sprintf("%s:%s", kv.Keyspace(), k)
		longKeys[i] = longKey
	}
	objs := kv.store.client.MGet(ctx, longKeys...).Val()
	for i, k := range shortKeys {
		m[k] = objs[i].([]byte)
	}
	return m
}
