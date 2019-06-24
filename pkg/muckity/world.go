package muckity

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
)

// GenericWorld is the default implementation of World
type GenericWorld struct {
	id        interface{}
	name      string
	parentCtx Context
	systems   []SystemRef
}

var _ World = &GenericWorld{}

func (w *GenericWorld) Name() string {
	return w.name
}

func (w *GenericWorld) Context() Context {
	// TODO: utilize context
	return w.parentCtx
}

func (w *GenericWorld) String() string {
	return fmt.Sprintf("%v", w.Name())
}

func (w *GenericWorld) AddSystems(systems ...System) {
	for _, system := range systems {
		sysRef := new(SystemRef)
		system = system
		w.systems = append(w.systems, *sysRef)
	}
}

var _ World = &GenericWorld{}

func (w *GenericWorld) GetSystem(name string) SystemRef {
	var ref SystemRef
	for _, ref = range w.systems {
		if w.Name() == name {
			return ref
		}
	}
	panic("Could not find requested system! Try using GetSystems")
}

func (w *GenericWorld) GetSystems() []SystemRef {
	return w.systems
}

func (w *GenericWorld) BSON() interface{} {
	p := bson.D{
		{"$set", bson.D{
			{
				"name",
				w.Name(),
			},
		}},
	}
	return p
}

func GetWorld(doOnce bool, ctx ...interface{}) World {
	var (
		wCtx    Context
		world   World
	)
	if len(ctx) < 2 {
		// TODO: utilize context
		wCtx = Background()
		world = &GenericWorld{nil, "descriptive-world", wCtx, make([]SystemRef, 0)}
		if len(ctx) == 1 {
			switch v := ctx[0].(type) {
			case string:
				world = &GenericWorld{nil, v, wCtx, make([]SystemRef, 0)}
			case Config:
				// Required config model is directly passed
				world.AddSystems(v)
			case World:
				world = v
			default:
				panic("Unknown interface passed to GetWorld()")
			}
		}
	}
	return world
}
