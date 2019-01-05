package muckity

import (
	"context"
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
	Context() context.Context
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
	MuckityType
}

// MuckityStorage implements base storage system, de-coupled from object persistence and configuration
type MuckityStorage interface {
	// Save object to storage provider; returns an error if anything failed
	// TODO: implement this better; move Save() to the MuckityPersistent, which calls a different function, defined here
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

// Tertia is the next division down from a second; it's also a tick
const Tertia = time.Second / 60
