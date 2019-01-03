package muckity

type muckityConfig struct {
	mongo mongoConfig
	// TODO: Config the things
}

type MuckityConfig struct {

}

// TODO: Make this return a system and use init()
func GetMuckityConfig() (*muckityConfig, error) {
	gc := new(muckityConfig)
	return gc, nil
}

// Type implements part of MuckitySystem
func (mc MuckityConfig) Type() string {
	return "systems"
}

// Channels implements part of MuckitySystem
func (mc MuckityConfig) Channels() []chan interface{} {
	// TODO: implement context handling / closing
	return make([]chan interface{}, 0)
}
