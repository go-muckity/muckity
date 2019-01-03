package muckity

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
}

type muckityConfig struct {
	mongo mongoConfig
	// TODO: Config the things
}

// TODO: Make this return a system and use init()
func GetMuckityConfig() (*muckityConfig, error) {
	gc := new(muckityConfig)
	return gc, nil
}
