package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	red "github.com/sean9999/go-store/redis"
)

func TestListCollection(t *testing.T) {

	ctx := context.Background()
	uniqueKey := fmt.Sprintf("test/list/%x", time.Now().Nanosecond())
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

	t.Run("colours exist", func(t *testing.T) {
		colours, err := s.ListCollection(ctx, "colours")
		if err != nil {
			t.Error(err)
		}
		got := s.ListCollections(ctx)["colours"]
		want := colours
		if s.ListCollections(ctx)["colours"] != colours {
			t.Errorf("wanted %d but got %d", want, got)
		}
	})

	t.Run("three animals", func(t *testing.T) {
		colours, err := s.ListCollection(ctx, "colours")
		if err != nil {
			t.Error(err)
		}
		colours.Push(ctx, "red")
		colours.Push(ctx, "blue")
		colours.Push(ctx, "yellow")
		colours.Push(ctx, "green")
		colours.Push(ctx, "purple")
		colours.Push(ctx, "orange")
		colours.Push(ctx, "brown")
		want := []string{"red", "blue", "yellow", "green", "purple", "orange", "brown"}
		gotStrings := []string{}
		for _, colour := range colours.All(ctx) {
			gotStrings = append(gotStrings, colour.(string))
		}
		if !hasSameElements(gotStrings, want) {
			t.Errorf("elements in list %v are not the same as %v", want, gotStrings)
		}
	})

	t.Run("destroy", func(t *testing.T) {

		err := s.Destroy(ctx)
		if err != nil {
			t.Error(err)
		}

	})

}
