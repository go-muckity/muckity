package muckity

import (
	"time"
)

// MuckityType is used for information discovery
type MuckityType interface {
	// Name should be the short name of the object; not necessarily a slug
	Name() string
	// Type should be the name of the type for storage purposes.
	// It should include two or more parts, separated by colon (:) as such:
	// `muckity:example`; this is to indicate it's defined in the muckity package
	// and is a type of "example".  This is only used as metadata, not runtime.
	Type() string
}

// MuckitySystem is used for contextual information discovery
type MuckitySystem interface {
	// MuckityContext returns the context for the system in question.
	// For worlds, it should return a derived child context from the world's active context.
	// Systems can pass whatever context they want down the chain, as long as it's a child of
	// the world that spawned the system.
	// TODO: implement this better; this is not the best use-case for context.Context
	Context() MuckityContext
	MuckityType
}

type MuckitySystemRef struct {
	system MuckitySystem
}

func (msr MuckitySystemRef) GetSystem() MuckitySystem {
	return msr.system
}

// MuckityPersistent is an interface for persisting items, de-coupled from storage system
type MuckityPersistent interface {
	// BSON returns bson primitives; should be re-usable to generate JSON, for example
	BSON() interface{}
	// GetId returns the string key id value according to the referenced value
	GetId() string
	// SetId accepts a string and applies it to the receiver for retrieval by GetId()
	SetId(key string)
}

// MuckityWorld models the implementation of a central management system; a "world" in mu* terms
type MuckityWorld interface {
	AddSystems(systems ...MuckitySystem)
	GetSystem(name string) MuckitySystemRef
	GetSystems() []MuckitySystemRef
	MuckitySystem
	// MuckityPersistent is included to require some kind of persistence modeling on worlds; can use
	// simple implementations:
	//
	// func (w *worldStruct) BSON() interface{} { var iface interface{}; return iface }
	// func (w *worldStruct) GetId() string { return "my-world-id" }
	// func (w *worldStruct) SetID(key string) { return }
	MuckityPersistent
}

// MuckityStorage implements base storage system, de-coupled from object persistence and configuration
type MuckityStorage interface {
	// Save object to storage provider; returns an error if anything failed
	// TODO: implement this better; move Save() to the MuckityPersistent interface; create GetStorage() on MuckityWorld
	Save(obj MuckityPersistent) error
	MuckitySystem
}

// MuckityConfig is the interface used to access config; based on viper as the model implementation uses viper
type MuckityConfig interface {
	Get(k string) interface{}
	Set(k string, v interface{})
	BindEnv(input ...string) error
	MuckitySystem
}

// MuckityContext implements context.Context w/ additional methods (at some point)
type MuckityContext interface {
	// Deadline implemented per context.Deadline
	Deadline() (deadline time.Time, ok bool)
	// Done implemented per context.Done()
	Done() <-chan struct{}
	// Err implemented per context.Err()
	Err() error
	// Value implemented per context.Value()
	Value(key interface{}) interface{}
	// Root() returns MuckityType
	Root() MuckityType
	// Config() returns MuckityConfig
	Config() MuckityConfig
	// World() returns MuckityWorld
	World() MuckityWorld
	// Storage() returns MuckityStorage
	Storage() MuckityStorage
	// System() returns MuckitySystem
	CallingSystem() MuckitySystem
	MuckityType
}

// Tertia is the next division down from a second; it's also a tick
const Tertia = time.Second / 60
