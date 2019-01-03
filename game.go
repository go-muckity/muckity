package main

import (
	"fmt"
	"github.com/machiel/slugify"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/tsal/muckity/pkg/muckity"
)

type myWorld struct {
	name        string
	description string
	zones       []string
}

func (w myWorld) Name() string {
	return w.name
}

func (w myWorld) Type() string {
	return "worlds"
}

func (w myWorld) Aliases() []string {
	s := make([]string, 0)
	s = append(s, "something")
	return s
}

func (w myWorld) DBId() string {
	return fmt.Sprintf("world:%v", slugify.Slugify(w.name))
}

func (w myWorld) Metadata() interface{} {
	return w.PersistentData()
}

func (w myWorld) PersistentData() interface{} {
	p := bson.D{
		{"$set", bson.D{
			{
				"name",
				w.Name(),
			},
			{
				"description",
				w.description,
			},
			{
				"zones",
				w.zones,
			},
			{
				"aliases",
				w.Aliases(),
			},
		}},
	}
	return p
}

type myTicker struct {

}

func main() {
	storage := muckity.GetMuckityStorage()
	w := new(myWorld)
	w.name = "Descriptive, aliased, world"
	w.description = `I am a really descriptive world.
I'm using a custom struct that implements the Persistent interface.
`
	w.zones = make([]string, 0)
	_, err := storage.Save(w)
	if err != nil {
		panic(err)
	}
}
