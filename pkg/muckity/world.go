package muckity

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
)

// World is the default implementation of MuckityWorld
type World struct {
	id        interface{}
	name      string
	mType     string
	parentCtx context.Context
	systems   []MuckitySystemRef
}

var _ MuckitySystem = &World{}

func (w *World) Name() string {
	return w.name
}

func (w *World) Type() string {
	return w.mType
}

func (w *World) Context() context.Context {
	// TODO: utilize context
	return context.TODO()
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

// NewWorld returns an instance of World, attaching instances of MuckitySystems passed in after config
func NewWorld(config MuckityConfig, systems ...MuckitySystem) *World {
	var w = new(World)
	// TODO: utilize context
	w.parentCtx = context.Background()
	w.mType = "world"
	config.BindEnv("world.name", "MUCKITY_WORLD_NAME")
	if name, ok := config.Get("world.name").(string); ok {
		w.name = name
	} else {
		w.name = "" // Blank name or config tests would be pointless...
	}
	if len(systems) > 0 {
		w.AddSystems(systems...)
	}
	w.id = fmt.Sprintf("%v:%v", w.Type(), w.Name())
	return w
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
