# provider

Provider is a simple library to dependency inversion library for Go.

## Installation

```bash
go get github.com/snowmerak/provider
```

## Usage

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/snowmerak/provider"
)

type Person struct {
	Name string
	Age  int
}

func NewPerson(ctx context.Context) (*Person, error) {
	name, ok := ctx.Value("name").(string)
	if !ok {
		return nil, errors.New("name is not found")
	}

	age, ok := ctx.Value("age").(int)
	if !ok {
		return nil, errors.New("age is not found")
	}

	return &Person{
		Name: name,
		Age:  age,
	}, nil
}

type House struct {
	owner *Person
}

func NewHouse(owner *Person) (*House, error) {
	return &House{
		owner: owner,
	}, nil
}

func main() {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "name", "John Doe")
	ctx = context.WithValue(ctx, "age", 25)

	pv := provider.New()

	if err := pv.Register(NewPerson, NewHouse); err != nil {
		panic(err)
	}

	if err := pv.Construct(ctx); err != nil {
		panic(err)
	}

	fmt.Println(provider.Get[*Person](pv))
	fmt.Println(provider.Get[*House](pv))

	fmt.Println(provider.Run[string](pv, func(p *Person) (string, error) {
		return fmt.Sprintf("name: %s, age: %d", p.Name, p.Age), nil
	}))
}
```

```bash
&{John Doe 25} true
&{0xc00008e060} true
name: John Doe, age: 25 <nil>
```
