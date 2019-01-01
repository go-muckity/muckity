package muckity

type MetaCollection interface {
	Name() string
	Type() string
	Metadata() interface{}
}

type Persistent interface {
	DBId() string
	PeristentData() map[string]interface{}
	Save() string
}

type MuckityObject interface {
	Name() string
	Type() string
}
