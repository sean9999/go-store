package redis

import (
	"context"
	"fmt"

	"github.com/sean9999/go-store/essence"
)

// a KeyValueCollection is a GenericCollection with methods suitable to a key-value store
type KeyValueCollection struct {
	collection
}

func (kv *KeyValueCollection) Destroy(ctx context.Context) error {
	shortKeys := kv.Keys(ctx)
	longKeys := make([]string, len(shortKeys))
	for i, k := range shortKeys {
		longKeys[i] = kv.store.LongKey(kv, string(k))
	}
	return kv.store.client.Del(ctx, longKeys...).Err()
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

func (kv *KeyValueCollection) Get(ctx context.Context, shortKey string) (any, error) {
	longKey := kv.store.LongKey(kv, shortKey)
	val, err := kv.store.client.Get(ctx, longKey).Result()
	return val, err
}

func (kv *KeyValueCollection) Set(ctx context.Context, shortKey string, v any) error {
	longKey := kv.store.LongKey(kv, shortKey)
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

func (kv *KeyValueCollection) All(ctx context.Context) map[string]any {
	m := map[string]any{}
	store := kv.store
	shortKeys := kv.Keys(ctx)
	longKeys := make([]string, len(shortKeys))
	for i, k := range shortKeys {
		longKeys[i] = store.LongKey(kv, string(k))
	}
	objs := kv.store.client.MGet(ctx, longKeys...).Val()
	for i, k := range shortKeys {
		m[k] = objs[i]
	}
	return m
}
