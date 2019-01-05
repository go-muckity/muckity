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

type MuckitySystem interface {
	// MuckityContext returns the context for the system in question.
	// For worlds, it should return a derived child context from the world's active context.
	// Systems can pass whatever context they want down the chain, as long as it's a child of
	// the world that spawned the system.
	// TODO: implement this better; this is not the best use-case for context.Context
	Context() context.Context
	MuckityType
}

type MuckityPersistent interface {
	// BSON returns bson primitives; should be re-usable to generate JSON, for example
	BSON() interface{}
	GetId() string
	SetId(key string)
	MuckityType

}

// MuckityStorage implements base storage system for later decoupling
type MuckityStorage interface {
	Save(obj MuckityPersistent) error
	MuckitySystem
}


// Tertia is the next division down from a second; it's also a tick
const Tertia = time.Second / 60

