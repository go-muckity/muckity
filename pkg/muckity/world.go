package muckity

import (
	"errors"
	"fmt"
	"github.com/machiel/slugify"
	"github.com/mongodb/mongo-go-driver/bson"
)

type World interface {
	MuckityObject
}

type world struct {
	name string
}

func NewWorld(name string, options ...func(*world) error) (world, error) {
	w := world{name}
	for _, op := range options {
		err := op(&w)
		if err != nil {
			return w, err
		}
	}
	res, err := storage.Save(w)
	if world, ok := res.(world); ok {
		return world, nil
	} else {
		panic(errors.New(fmt.Sprintf("Unknown error saving new world: %v vs. %v", w, world)))
	}
	if err != nil {
		panic(err)
	}
	return w, err
}

func (w world) Name() string {
	return w.name
}

func (w world) Type() string {
	return "worlds"
}

func (w world) Slug() string {
	return slugify.Slugify(w.Name())
}

func (w world) DBId() string {
	return fmt.Sprintf("world:%v", w.Slug())
}

func (w world) Metadata() interface{} {
	return w.PersistentData()
}
func (w world) PersistentData() interface{} {
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

func (w world) Save() (world, error) {
	if newWorld, err := storage.Save(w); err != nil {
		return w, err
	} else {
		if world, ok := newWorld.(world); ok {
			return world, nil
		}
	}
	return w, errors.New(fmt.Sprintf("Unknnown error saving world: ", w))
}
