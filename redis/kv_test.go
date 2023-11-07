package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sean9999/GoFunctional/fslice"
	red "github.com/sean9999/go-store/redis"
)

func hasSameElements(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	f2 := fslice.From(slice2)
	for _, v := range slice1 {
		if !f2.Includes(v) {
			return false
		}
	}
	return true
}

func TestKeyValueCollection(t *testing.T) {

	ctx := context.Background()
	uniqueKey := fmt.Sprintf("test/kv/%x", time.Now().Nanosecond())
	//var animals essence.KeyValueCollection
	s := red.Attach(uniqueKey, &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	t.Run("Store Namespace", func(t *testing.T) {
		got := s.Namespace
		want := uniqueKey
		if got != want {
			t.Errorf("wanted %q but got %q", want, got)
		}
	})

	t.Run("no list collections on a brand new store", func(t *testing.T) {
		got := len(s.ListCollections(ctx))
		want := 0
		if got != want {
			t.Errorf("wanted %d but got %d", want, got)
		}
	})

	t.Run("no kv collections on a brand new store", func(t *testing.T) {
		got := len(s.KeyValueCollections(ctx))
		want := 0
		if got != want {
			t.Errorf("wanted %d but got %d", want, got)
		}
	})

	t.Run("animals exists", func(t *testing.T) {
		animals, err := s.KeyValueCollection(ctx, "animals")
		if err != nil {
			t.Error(err)
		}
		got := s.KeyValueCollections(ctx)["animals"]
		want := animals
		if s.KeyValueCollections(ctx)["animals"] != animals {
			t.Errorf("wanted %d but got %d", want, got)
		}
	})

	t.Run("three animals", func(t *testing.T) {
		animals, err := s.KeyValueCollection(ctx, "animals")
		if err != nil {
			t.Error(err)
		}
		animals.Set(ctx, "dog", "bark")
		animals.Set(ctx, "cat", "meow")
		animals.Set(ctx, "cow", "moo")
		want := []string{"dog", "cat", "cow"}
		got := animals.Keys(ctx)
		if !hasSameElements(got, want) {
			t.Errorf("elements in list %v are not the same as %v", want, got)
		}
	})

	t.Run("destroy", func(t *testing.T) {

		err := s.Destroy(ctx)
		if err != nil {
			t.Error(err)
		}

	})

}
