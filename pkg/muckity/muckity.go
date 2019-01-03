package muckity

type muckityRoot interface {
	Name() string
	Type() string
}

// MuckityObject is the basic object in the Muckity ECS
type MuckityObject interface {
	muckityRoot
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
