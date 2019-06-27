# Go Muckity [![Build Status](https://travis-ci.org/go-muckity/muckity.svg?branch=master)](https://travis-ci.org/go-muckity/muckity)

## Usage

Take a look at `muckity/world.go` for a sample implementation of `WorldSystem`, first. Then, take a look at `main.go` for the test script of the interfaces.

As of this writing (v0.0.1), there are only basic entities and components (refererred to as _systems_ in muckity) provided in the framework.

First things first, you will need to import the ecs into your code.  Note, in this example, I'm importing `ecs` as `muckity` for aesthetics.

```go
import muckity "github.com/go-muckity/pkg/muckity"
```

Next, you will want to probably implement the interface `muckity.WorldSystem`:

```go
// System is used for contextual information discovery
type System interface {
	UUID() uuid.UUID
	Run() (int, error)
	Next() <-chan SystemFunc
	Init(...interface{}) error
	Update(System) error
	Shutdown()
	fmt.Stringer
}
// WorldSystem models the implementation of a central management system; a "world" in mu* terms
type WorldSystem interface {
	System
	Name() string
	Join(interface{}) error
	TickingSystem
}
type TickingSystem interface {
	Tick() error
}
```

## Project Goals

### v0.1.0

- [X] Implement base and basic systems
- [X] Implement agnostic storage and configuration interface examples
- [X] Create doc.go in ecs, and keep it updated
- [X] Implement context-handling

### v2.0.0

- [X] Implement concurrency in muckity example systems

### And further...

- [ ] Design and implement MuckityNetwork
- [ ] Implement MuckityClient interface
- [ ] Implement example networking system in test game (`main.go`)

## Modules

### `/`

Here we have the "game" created to test the development of the module (and later, packages other than `muckity`).  Long-term it will be used to flesh out design ideas for new interfaces in the `muckity` package and others.

### Packages

#### `/pkg/muckity`

The muckity package itself; World and Ticker systems code, used by the test game code.

