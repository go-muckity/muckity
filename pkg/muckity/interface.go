package muckity

import (
	"time"
)

type Namer interface {
	Name() string
}
// System is used for contextual information discovery
type System interface {
	// Context returns the context for the system in question.
	// For worlds, it should return a derived child context from the world's active context.
	// Systems can pass whatever context they want down the chain, as long as it's a child of
	// the world that spawned the system.
	// TODO: implement this better; this is not the best use-case for context.Context
	Context() Context
	Namer
}

type SystemRef struct {
	system System
}

func (msr SystemRef) GetSystem() System {
	return msr.system
}

// World models the implementation of a central management system; a "world" in mu* terms
type World interface {
	AddSystems(systems ...System)
	GetSystem(name string) SystemRef
	GetSystems() []SystemRef
	System
}

// Config is the interface used to access config; based on viper as the model implementation uses viper
type Config interface {
	Get(k string) interface{}
	Set(k string, v interface{})
	BindEnv(input ...string) error
	System
}

// Context implements context.Context w/ additional methods (at some point)
type Context interface {
	// Deadline implemented per context.Deadline
	Deadline() (deadline time.Time, ok bool)
	// Done implemented per context.Done()
	Done() <-chan struct{}
	// Err implemented per context.Err()
	Err() error
	// Value implemented per context.Value()
	Value(key interface{}) interface{}
	// Config() returns Config
	Config() Config
	// GenericWorld() returns World
	World() World
	// System() returns System
	CallingSystem() System
}

// Tertia is the next division down from a second; it's also a tick
const Tertia = time.Second / 60
