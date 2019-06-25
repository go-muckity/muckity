package muckity

import (
	"fmt"
)

var _ System = &GenericSystem{}

// GenericSystem is the default implementation of System
type GenericSystem struct {
	name string
}

func (s GenericSystem) String() string {
	return s.name
}

var _ World = &GenericWorld{}

// GenericWorld is the default implementation of World
type GenericWorld struct {
	id      interface{}
	name    string
	systems []System
}

func (w *GenericWorld) AddSystems(systems ...System) {
	for _, system := range systems {
		sysRef := new(GenericSystem)
		system = system
		w.systems = append(w.systems, *sysRef)
	}
}
func (w GenericWorld) GetSystem(name string) System {
	var ref System
	for _, ref = range w.systems {
		if w.name == name {
			return ref
		}
	}
	panic("Could not find requested system! Try using GetSystems")
}
func (w GenericWorld) GetSystems() []System {
	return w.systems
}
func (w GenericWorld) String() string {
	return fmt.Sprintf("%v", w.name)
}
func GetWorld() World {
	var (
		world World
	)
	world = &GenericWorld{nil, "generic-world", make([]System, 0)}
	return world
}
