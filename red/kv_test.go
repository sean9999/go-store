package red_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sean9999/GoFunctional/fslice"
	red "github.com/sean9999/go-store/red"
)

// test that two slices contain exactly the same elements regardless of order
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
	//	create a brand new store
	uniqueKey := fmt.Sprintf("test/kv/%x/%x", time.Now().Minute(), time.Now().Second())
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

	t.Run("no kv collections on a brand new store", func(t *testing.T) {
		got := len(s.KeyValueCollections(ctx))
		want := 0
		if got != want {
			t.Errorf("wanted %d but got %d", want, got)
		}
	})

	t.Run("animals exists", func(t *testing.T) {
		//	create or attach to animals (in this case create)
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
		//	create or attach to animals (in this case attach)
		animals, err := s.KeyValueCollection(ctx, "animals")
		if err != nil {
			t.Error(err)
		}
		animals.Set(ctx, "dog", []byte("bark"))
		animals.Set(ctx, "cat", []byte("meow"))
		animals.Set(ctx, "cow", []byte("moo"))
		want := []string{"dog", "cat", "cow"}
		got := animals.Keys(ctx)
		if !hasSameElements(got, want) {
			t.Errorf("elements in list %v are not the same as %v", want, got)
		}
	})

	t.Run("destroy", func(t *testing.T) {

		//	destroy all collections in the store
		//	thereby effectively destroying the store itself
		err := s.Destroy(ctx)
		if err != nil {
			t.Error(err)
		}

	})

}
