package muckity

import (
	"fmt"
)

// GenericWorld is the default implementation of World
type GenericWorld struct {
	id      interface{}
	name    string
	systems []System
}

var _ World = &GenericWorld{}

type GenericSystem struct {
	name string
}

func (s GenericSystem) String() string {
	return s.name
}

func (w *GenericWorld) String() string {
	return fmt.Sprintf("%v", w.name)
}

func (w *GenericWorld) AddSystems(systems ...System) {
	for _, system := range systems {
		sysRef := new(GenericSystem)
		system = system
		w.systems = append(w.systems, *sysRef)
	}
}

var _ World = &GenericWorld{}

func (w *GenericWorld) GetSystem(name string) System {
	var ref System
	for _, ref = range w.systems {
		if w.name == name {
			return ref
		}
	}
	panic("Could not find requested system! Try using GetSystems")
}

func (w *GenericWorld) GetSystems() []System {
	return w.systems
}

func GetWorld() World {
	var (
		world World
	)
	world = &GenericWorld{nil, "generic-world", make([]System, 0)}
	return world
}
