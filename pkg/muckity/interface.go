package muckity

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Message interface {
}
type SystemMap map[string]System

// UnknownSystemMapTypeErr implements a comparable error type to allow graceful handling
var StringOrSystemErr = fmt.Errorf("can only delete using system uuid or a System interface")

// ExistingSystemMapSystemErr implements a comparable error type to allow graceful handling
var CannotAddExistingErr = fmt.Errorf("cannot add existing system")

func (m SystemMap) exists(key string) bool {
	_, ok := m[key]
	return ok
}
func (m SystemMap) Get(s string) System {
	if v, ok := m[s]; ok {
		return v
	}
	return nil
}
func (m SystemMap) Add(s System) error {
	if m == nil {
		return fmt.Errorf("nil SystemMap")
	}
	if m.exists(s.String()) {
		return CannotAddExistingErr
	}
	m[s.UUID().String()] = s
	return nil
}
func (m SystemMap) Del(s interface{}) error {
	switch s.(type) {
	case string:
		v, _ := s.(string)
		delete(m, v)
		return nil
	case System:
		v, _ := s.(System)
		delete(m, v.String())
		return nil
	default:
		return StringOrSystemErr
	}
}
func (m SystemMap) Init(systems ...System) error {
	for _, system := range systems {
		err := m.Add(system)
		if err != nil {
			return err
		}
	}
	return nil
}

type InitConfig struct {
	Name           string
	PubChannel     chan Message
	SubChannel     chan Message
	MessageChannel chan Message
	CloseFunc      func()
	Systems        SystemMap
	HandlerFunc    HandlerFunc
}
type Muckity interface {
	UUID() uuid.UUID // hard requirement for a valid uuid
	MessageChannel() chan<- Message
	Closer() func()
	Systems() SystemMap
	Handler(context.Context, Message) (Message, error)
	Init(InitConfig, ...RootFunc) error
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
}

type HandlerFunc func(context.Context, Message) (Message, error)
type RootFunc func(muckity Muckity) error
type SystemFunc func(system System) error

// System is used for contextual information discovery
type System interface {
	UUID() uuid.UUID
	Run() (int, error)
	Next() <-chan SystemFunc
	Init(...interface{}) error
	Update(System) error
	Shutdown()
	fmt.Stringer
}

// WorldSystem models the implementation of a central management system; a "world" in mu* terms
type WorldSystem interface {
	System
	Name() string
	Join(interface{}) error
	TickingSystem
}

type TickingSystem interface {
	Tick() error
}

type Ticker interface {
	Rate() time.Duration
	TargetedSystem
	System
}

type TargetedSystem interface {
	Target(interface{}) error
	Targets() []interface{}
	Untarget(interface{}) error
}

// Tertia is the next division down from a second; it's also a tick
const Tertia = time.Second / 60
