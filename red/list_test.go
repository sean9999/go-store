package red_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sean9999/go-store/essence"
	red "github.com/sean9999/go-store/red"
)

func TestListCollection(t *testing.T) {

	ctx := context.Background()
	uniqueKey := fmt.Sprintf("test/list/%x/%x", time.Now().Minute(), time.Now().Second())
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

	t.Run("genres", func(t *testing.T) {
		genres, err := s.ListCollection(ctx, "genres")
		if err != nil {
			t.Error(err)
		}
		got := s.ListCollections(ctx)["genres"]
		want := genres
		if s.ListCollections(ctx)["genres"] != genres {
			t.Errorf("wanted %d but got %d", want, got)
		}
	})

	t.Run("there are two lists: colours and genres", func(t *testing.T) {
		colours, err := s.ListCollection(ctx, "colours")
		if err != nil {
			t.Error(err)
		}
		genres, err := s.ListCollection(ctx, "genres")
		if err != nil {
			t.Error(err)
		}

		got := s.ListCollections(ctx)
		want := map[string]essence.ListCollection{
			"colours": colours,
			"genres":  genres,
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("wanted %v but got %v", want, got)
		}
	})

	t.Run("add 7 genres", func(t *testing.T) {
		genres, err := s.ListCollection(ctx, "genres")
		if err != nil {
			t.Error(err)
		}
		genres.Push(ctx, "bacon")
		genres.Push(ctx, "country")
		genres.Push(ctx, "rap")
		genres.Push(ctx, "jazz")
		genres.Push(ctx, "classical")
		genres.Push(ctx, "pop")
		genres.Push(ctx, "Matilda")
		got := genres.Size(ctx)
		want := 7
		if got != want {
			t.Errorf("wanted %d but got %d", want, got)
		}
	})

	t.Run("bacon and Matilda are not genres", func(t *testing.T) {
		genres, err := s.ListCollection(ctx, "genres")
		if err != nil {
			t.Error(err)
		}
		shiftedElement, _ := genres.Shift(ctx)
		poppedElement, _ := genres.Pop(ctx)

		if shiftedElement != "bacon" {
			t.Errorf("wanted bacon but got %q", shiftedElement)
		}
		if poppedElement != "Matilda" {
			t.Errorf("wanted Matilda but got %q", poppedElement)
		}
		got := genres.Size(ctx)
		want := 5
		if got != want {
			t.Errorf("new list size should have been %d but we got %d", want, got)
		}
	})

	t.Run("destroy", func(t *testing.T) {

		err := s.Destroy(ctx)
		if err != nil {
			t.Error(err)
		}

	})

}
