package muckity

import (
	"time"
)

type Namer interface {
	Name() string
}

// System is used for contextual information discovery
type System interface {
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

// Tertia is the next division down from a second; it's also a tick
const Tertia = time.Second / 60
