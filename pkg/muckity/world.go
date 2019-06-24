package muckity

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
)

// GenericWorld is the default implementation of World
type GenericWorld struct {
	id      interface{}
	name    string
	systems []SystemRef
}

var _ World = &GenericWorld{}

func (w *GenericWorld) Name() string {
	return w.name
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

func GetWorld() World {
	var (
		world World
	)
	world = &GenericWorld{nil, "generic-world", make([]SystemRef, 0)}
	return world
}
