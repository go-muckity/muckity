package muckity

import "time"

type muckityRoot interface {
	Name() string
	Type() string
}

type muckityObject interface {
	Aliases() []string
}

// MuckityObject is the basic object in the Muckity ECS
type MuckityObject interface {
	muckityRoot
	muckityObject
}

type muckitySystem interface {
	Channels() []chan interface{}
}

type MuckitySystem interface {
	muckityRoot
	muckitySystem
}

// Tertia is the next division down from a second; it's also a tick
const Tertia = time.Second / 60
