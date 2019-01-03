package muckity

import (
	"errors"
	"fmt"
)

type muckityCore struct {
	// TODO: further categorize once that's sane
	config *MuckityConfig
	// TODO: Make a generic Muckity interface
	storage *MuckityStorage
}

type MuckityCore struct {
	done chan interface {}
	systems map[string]interface{}
	muckityCore
}

func (c MuckityCore) Type() string {
	return "systems"
}

func (c MuckityCore) Channels() map[string]chan interface{} {
	slice := make(map[string]chan interface{}, 0)
	slice["done"] = c.done
	return slice
}

func (c MuckityCore) Config() *MuckityConfig {
	config, err := c.GetSystem("config")
	if err != nil {
		panic(err)
	}
	if config, ok := config.(*MuckityConfig); ok {
		return config
	} else {
		return nil
	}
	return nil
}

func (c MuckityCore) Storage() *MuckityStorage {
	return c.storage
}

func (c MuckityCore) GetSystem(name string) (interface{}, error) {
	if val, ok := c.systems[name]; ok {
		return val, nil
	}
	return nil, errors.New(fmt.Sprintf("Unknown system: %v", name))
}

var core *MuckityCore

// initialize core systems
func Init() * MuckityCore {
	core = new(MuckityCore)
	core.done = make(chan interface{})
	core.systems = make(map[string]interface{}, 0)
	core.systems["config"] = configInit(core.done)
	core.systems["storage"] = dbInit(core.done)
	return core
}
