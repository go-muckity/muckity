# Go Muckity [![Build Status](https://travis-ci.org/go-muckity/muckity.svg?branch=master)](https://travis-ci.org/go-muckity/muckity)

## Usage

Take a look at `ecs/world.go` for a sample implementation of `MuckityWorld`, first. Then, take a look at `main.go` for the game engine developed to test the entity component system (ecs).

As of this writing (v0.0.1), there are only basic entities and components (refererred to as _systems_ in muckity) provided in the framework.

First things first, you will need to import the ecs into your code.  Note, in this example, I'm importing `ecs` as `muckity` for aesthetics.

```go
import muckity "github.com/go-muckity/muckity/ecs"
```

Next, you will need to implement the interface `ecs.MuckityWorld`:

```go
type MuckityWorld interface {
	AddSystems(systems ...MuckitySystem) // Add one ore more systems under the control of this World
	GetSystem(name string) MuckitySystemRef // Should return a system reference
	GetSystems() []MuckitySystemRef
    // These are part of MuckitySystem and MuckityType, inherited in the declaration
	Context() context.Context // TODO() 
	Name() string // A string representing the name of the world
	Type() string // A string representing a metadata "type" of the world, ie "muckity:world"
	// These are part of MuckityPersistent, also part of the declaration, more on that in the api docs
	BSON() interface{} // returns bson primitives from github.com/mongodb/mongo-go-driver/bson
	GetId() string // returns a string representing what SHOULD be a unique identifier for the world (see below)
	SetId(key string) // Set the id; does not have to be unique, if your implementation allows for it
}
```

Note that `ecs.MuckitySystem` small value wrapping a pointer interface; `ecs.MuckitySystemRef` is a container of a single `ecs.MuckitySystem`, which is a pointer value under the hood.
 
It's exported to allow developers of other implementations to implement helpers for special cases:

```go
func (msr MuckitySystemRef) GetSlug() string {
	var name string = msr.GetSystem().Name() // Calls .Name() on MuckitySystem interface
	// make `name` a slug so our custom storage system doesn't choke on the name
	return name
}
```

Essentially, this allows a world to be a parent system or sibling system or even a child sibling to some other interface, and it can still continue to manage other systems.

```go
import "github.com/go-muckity/muckity/ecs"

var _ ecs.MuckityWorld = &myWorld{} // Useful testing at compile-time if you missed something.

type myWorld struct {
	id string // this is all we need
}

func (w *myWorld) AddSystems(s ...ecs.MuckitySystem) { return // who cares! }
func (w *myWorld) GetSystem(name string) ecs.MuckitySystemRef { return ecs.MuckitySystemRef{} }
func (w *myWorld) GetSystems() []ecs.MuckitySystemRef { return make([]ecs.MuckitySystemRef, 0) }
func (w *myWorld) Context() context.Context { return context.TODO() }
func (w *myWorld) Name() string { return w.GetId() }
func (w *myWorld) Type() string { return "myWorld" }
func (w *myWorld) BSON() interface{} { 
	var w interface{}
	return w
}
func (w *myWorld) GetId() string { return w.id }
func (w *myWorld) SetId(key string) string { w.id = key }
```

This is the bare-minimum implementation of `ecs.MuckityWorld`; it will work with `ecs.MongoStorage`, which is an implementation of `ecs.MuckityStorage`, the next thing to implement if you don't want to use a mongodb backend.

It will also should work with the framework implementation of `ecs.MuckityConfig`, `ecs.muckityConfig` - the default object you get from ecs.GetConfig().

You may note there is not a framework implementation similar to `ecs.muckityConfig` - this is because it hasn't been designed yet, but it will very likely be informed by the mongodb implementation, so that's staying for now.

The long-term intent is to have the framework implementation be flat-file storage, likely YAML or JSON, since both should be compatible with native BSON objects, allowing for custom systems to store binary data with the framework implementation.

## Project Goals

### v0.1.0

- [X] Implement base and basic systems
- [X] Implement agnostic storage and configuration interface examples
- [ ] Create doc.go in ecs, and keep it updated
- [ ] Implement context-handling

### v0.2.0

- [ ] Implement concurrency in ecs example systems
- [ ] Implement example networking system in test game (`main.go`)

### And further...

- [ ] Design and implement MuckityNetwork
- [ ] Implement MuckityClient interface

## Modules

### `/`

Here we have the "game" created to test the development of the module (and later, packages other than `ecs`).  Long-term it will be used to flesh out design ideas for new interfaces in the `ecs` package and others.

### Packages

#### `/ecs`

The ECS, storage and other systems code, used by the game code. Examples are above, but it isn't a lot of code to read through.

