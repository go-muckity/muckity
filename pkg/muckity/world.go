package muckity

import (
	"fmt"
	"github.com/machiel/slugify"
	"github.com/mongodb/mongo-go-driver/bson"
)

type MuckityWorld struct {
	name string
}

func (w MuckityWorld) Name() string {
	return w.name
}

func (w MuckityWorld) Type() string {
	return "worlds"
}

func (w MuckityWorld) Aliases() []string {
	return make([]string, 0)
}

func (w MuckityWorld) DBId() string {
	return fmt.Sprintf("world:%v", slugify.Slugify(w.name))
}

func (w MuckityWorld) Metadata() interface{} {
	return w.PersistentData()
}
func (w MuckityWorld) PersistentData() interface{} {
	p := bson.D{
		{"$set", bson.D{
			{
				"name",
				w.Name(),
			},
		}},
	}
	return p
}
