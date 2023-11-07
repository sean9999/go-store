# Store

Store is a simple wrapper around redis. It provides namespaces and the concept of Collections.

Other back-ends beside redis are planned.

## Getting started

First you must have a running instance of redis.

```sh
$ sudo systemctl start redis-server
```

A simple example main.go:

```go
import (
	"context"
	"fmt"

	"github.com/sean9999/go-store"
)

func main(){

    //  using the redis back-end
	ctx := context.Background()
	s := red.NewStore("myStore")

    animals, _ := s.KeyValueCollection("animals")
    
    animals.Set(ctx, "duck", "quack")
    animals.Set(ctx, "dog", "barf")

    allAnimals := animals.GetAll(ctx)

    fmt.Println(allAnimals)

}
```