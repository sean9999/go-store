package red

import (
	"context"

	"github.com/sean9999/go-store/essence"
)

type ListCollection struct {
	collection
}

func (lc *ListCollection) Destroy(ctx context.Context) error {
	//err := lc.store.client.LTrim(ctx, lc.Keyspace(), 1, 0).Err()
	err := lc.store.client.Del(ctx, lc.Keyspace()).Err()
	if err == nil {
		lc.store.RemoveCollection(ctx, lc)
	}
	return err
}

func (s *Store) ListCollection(ctx context.Context, name string) (essence.ListCollection, error) {
	lc, exists := s.collectionMap.list[name]
	if !exists {
		c := collection{
			kind:  "list",
			name:  name,
			store: s,
		}
		lc = &ListCollection{c}
		s.addCollection(ctx, lc)
	}
	s.registerCollection(ctx, lc)
	return lc, nil
}

func (lc *ListCollection) All(ctx context.Context) []any {
	length := lc.store.client.LLen(ctx, lc.Keyspace()).Val()
	r := []any{}
	for i := int64(0); i < length; i++ {
		x := lc.store.client.LIndex(ctx, lc.Keyspace(), i).Val()
		r = append(r, x)
	}
	return r
}

func (lc *ListCollection) Head(ctx context.Context) (any, error) {
	x, err := lc.store.client.LIndex(ctx, lc.Keyspace(), 0).Result()
	return x, err
}

func (lc *ListCollection) Tail(ctx context.Context) (any, error) {
	x, err := lc.store.client.LIndex(ctx, lc.Keyspace(), -1).Result()
	return x, err
}

func (lc *ListCollection) Pop(ctx context.Context) (any, error) {
	x, err := lc.store.client.RPop(ctx, lc.Keyspace()).Result()
	return x, err
}

func (lc *ListCollection) Shift(ctx context.Context) (any, error) {
	x, err := lc.store.client.LPop(ctx, lc.Keyspace()).Result()
	return x, err
}

func (lc *ListCollection) Push(ctx context.Context, val any) error {
	err := lc.store.client.RPush(ctx, lc.Keyspace(), val).Err()
	return err
}

func (lc *ListCollection) Unshift(ctx context.Context, val any) error {
	err := lc.store.client.LPush(ctx, lc.Keyspace(), val).Err()
	return err
}

func (lc *ListCollection) Size(ctx context.Context) int {
	size := lc.store.client.LLen(ctx, lc.Keyspace()).Val()
	return int(size)
}
