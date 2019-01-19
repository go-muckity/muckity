package ecs

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
)

// World is the default implementation of MuckityWorld
type World struct {
	id        interface{}
	name      string
	mType     string
	parentCtx MuckityContext
	systems   []MuckitySystemRef
}

var _ MuckityWorld = &World{}

func (w *World) Name() string {
	return w.name
}

func (w *World) Type() string {
	return w.mType
}

func (w *World) Context() MuckityContext {
	// TODO: utilize context
	return w.parentCtx
}

func (w *World) String() string {
	return fmt.Sprintf("%v:%v", w.Type(), w.Name())
}

func (w *World) AddSystems(systems ...MuckitySystem) {
	for _, system := range systems {
		sysRef := new(MuckitySystemRef)
		sysRef.system = system
		w.systems = append(w.systems, *sysRef)
	}
}

var _ MuckityWorld = &World{}

func (w *World) GetSystem(name string) MuckitySystemRef {
	var ref MuckitySystemRef
	for _, ref = range w.systems {
		if ref.GetSystem().Name() == name {
			return ref
		}
	}
	panic("Could not find requested system! Try using GetSystems")
}

func (w *World) GetSystems() []MuckitySystemRef {
	return w.systems
}

func (w *World) BSON() interface{} {
	p := bson.D{
		{"$set", bson.D{
			{
				"name",
				w.Name(),
			},
			{
				"type",
				w.Type(),
			},
		}},
	}
	return p
}

func (w *World) GetId() string {
	// TODO: needs better checking
	return fmt.Sprintf("%v", w.id)
}

func (w *World) SetId(id string) {
	w.id = id
}

func GetWorld(doOnce bool, ctx ...interface{}) MuckityWorld {
	var (
		wCtx    MuckityContext
		world   MuckityWorld
		storage MuckityStorage
	)
	if len(ctx) < 2 {
		// TODO: utilize context
		wCtx = Background()
		world = &World{nil, "descriptive-world", "world", wCtx, make([]MuckitySystemRef, 0)}
		storage = GetStorage(wCtx)
		if len(ctx) == 0 {
			world.AddSystems(GetConfig(), storage)
		}
		if len(ctx) == 1 {
			switch v := ctx[0].(type) {
			case string:
				world = &World{nil, v, "world", wCtx, make([]MuckitySystemRef, 0)}
			case MuckityConfig:
				// Required config model is directly passed
				world.AddSystems(v)
				world.AddSystems(storage)
			case MuckityWorld:
				world = v
			default:
				panic("Unknown interface passed to GetWorld()")
			}
		}
		world.SetId(fmt.Sprintf("%v:%v", world.Type(), world.Name()))
	}
	storage.Save(world) // nop if storage isn't defined.
	return world
}
