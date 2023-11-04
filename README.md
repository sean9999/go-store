# Store

Store is a simple wrapper around redis. It provides namespaces and the concept of Collections.

## Getting started

First you must have a running instance of redis.

```go
import (
	"context"
	"fmt"

	"github.com/sean9999/go-store"
)

func main(){
	ctx := context.Background()
	store := store.NewStore("myStore")

    animals, _ := NewKeyValueCollection("animals")
    
    animals.Set(ctx, "duck", "quack")
    animals.Set(ctx, "dog", "barf")

    allAnimals := animals.GetAll(ctx)

    fmt.Println(allAnimals)

}
```