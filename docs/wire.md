# Wire

Wire is a code generation tool that automates connecting components using dependency injection.

Wire was primarily inspired by Java’s Dagger 2, and uses code generation rather than reflection or service locators.

- [official document](https://go.dev/blog/wire)

- existing runtime DI tools (use reflection to do runtime dependency injection)
  - Uber's dig
  - Facebook's inject both


## Install
```
$ go install github.com/google/wire/cmd/wire@latest
```

## How does it work?
Wire has two basic concepts: providers and injectors.

1. Provider: `wire.ProivderSet()`
2. Injector: `wire.Build()`

## 1. Provider
> Providers are ordinary Go functions that “provide” values given their dependencies, which are described simply as parameters to the function.

- Given 3 providers
```go
// NewUserStore is the same function we saw above; it is a provider for UserStore,
// with dependencies on *Config and *mysql.DB.
func NewUserStore(cfg *Config, db *mysql.DB) (*UserStore, error) {...}

// NewDefaultConfig is a provider for *Config, with no dependencies.
func NewDefaultConfig() *Config {...}

// NewDB is a provider for *mysql.DB based on some connection info.
func NewDB(info *ConnectionInfo) (*mysql.DB, error) {...}
```
- Group with ProviderSet

```go
var UserStoreSet = wire.ProviderSet(NewUserStore, NewDefaultConfig)
```
## 2. Injector
> Injectors are generated functions that call providers in dependency order.


- Build provider sets
```
// File: wire_gen.go
// Code generated by Wire. DO NOT EDIT.
//go:generate wire
//+build !wireinject

func initUserStore(info ConnectionInfo) (*UserStore, error) {
    defaultConfig := NewDefaultConfig()
    db, err := NewDB(info)
    if err != nil {
        return nil, err
    }
    userStore, err := NewUserStore(defaultConfig, db)
    if err != nil {
        return nil, err
    }
    return userStore, nil
}
```

finally type cmd to generate `DI` interface code.

```
$ go generate
```

## Tutorial
> https://github.com/google/wire/blob/main/_tutorial/README.md

First we create a message, then we create a greeter with that message, and finally we create an event with that greeter. With all the initialization done, we're ready to start our event.

We are using the dependency injection design principle. In practice, that means we pass in whatever each component needs. This style of design lends itself to writing easily tested code and makes it easy to swap out one dependency with another.

```go
// entity.go
package main

import "github.com/minkj1992/fiber-wire/entity"

func beforeWire() {
	message := entity.NewMessage()
	greeter := entity.NewGreeter(message)
	event := entity.NewEvent(greeter)

	event.Start()
}

func afterWire() {

}

func main() {
	beforeWire()
}


// main.go
package main

import "github.com/minkj1992/fiber-wire/entity"

func beforeWire() {
	message := entity.NewMessage()
	greeter := entity.NewGreeter(message)
	event := entity.NewEvent(greeter)

	event.Start()
}

func main() {
	beforeWire()
}

```


## refs
- [guide](https://github.com/google/wire/blob/main/docs/guide.md)
- [best practice](https://github.com/google/wire/blob/main/docs/best-practices.md)