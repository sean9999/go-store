package main

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sean9999/go-store/red"
)

func barfOn(e error, msg string) {
	if e != nil {
		fmt.Println("🤮\t", msg)
		panic(e)
	}
}

func main() {
	ctx := context.Background()
	myStore := red.Attach("v6", &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	animals, _ := myStore.KeyValueCollection(ctx, "animals")
	animals.Set(ctx, "bird", []byte("chirp"))

	duckSays, _ := animals.Get(ctx, "duck")
	fmt.Println("duckSays", duckSays)

	truckSays, err := animals.Get(ctx, "truck")
	if err == nil {
		fmt.Println("truckSays", truckSays)
	} else {
		fmt.Println("truck is not an animal", err)
	}

	colours, err := myStore.ListCollection(ctx, "colours")
	barfOn(err, "colours")
	err = colours.Push(ctx, []byte("green"))
	barfOn(err, "green")
	colours.Push(ctx, []byte("blue"))
	colours.Push(ctx, []byte("red"))
	colours.Push(ctx, []byte("yellow"))

	all := colours.All(ctx)
	fmt.Println(all)

	colours.Destroy(ctx)
	animals.Destroy(ctx)

}
