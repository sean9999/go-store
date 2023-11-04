package main

import (
	"context"
	"fmt"

	"github.com/sean9999/go-store"
)

func barfOn(e error, msg string) {
	if e != nil {
		fmt.Println("barf\t", msg)
		panic(e)
	}
}

func main() {
	ctx := context.Background()
	store := store.NewStore("v6")

	animals, err := store.GetKeyValueCollection("animals")
	barfOn(err, "animals")

	duck, err := animals.Get(ctx, "duck")
	barfOn(err, "duck")

	fmt.Println("duck", duck)

	truck, err := animals.Get(ctx, "Truck")
	if err == nil {
		fmt.Println("truck", truck)
	} else {
		fmt.Println("truck is not an animal", err)
	}

}
