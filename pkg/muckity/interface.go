package muckity

import (
	"encoding/json"
	"fmt"
	"github.com/mongodb/mongo-go-driver/x/mongo/driver/uuid"
	"time"
)

type Message interface {
	Value() interface{}
}

type Muckity interface {
	UUID() uuid.UUID
	publisher() <-chan Message
	subscriber() chan<- Message
	messages() chan Message
	Cancel() func()
	Worlds(...int) []World
	GlobalSystems() []System
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
}

// System is used for contextual information discovery
type System interface {
	fmt.Stringer
}

// World models the implementation of a central management system; a "world" in mu* terms
type World interface {
	AddSystems(systems ...System)
	GetSystem(name string) System
	GetSystems() []System
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
