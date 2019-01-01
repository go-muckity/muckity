package muckity

import (
	"context"
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

func NewWorld(name string, options ...func(*world) error) (*world, error) {
	w := &world{name, }
	for _, op := range options {
		err := op(w)
		if err != nil {
			return nil, err
		}
	}
	w, err := w.Save()
	return w, err
}

func (w *world) Name() string {
	return w.name
}

func (w *world) Type() string {
	return "world"
}

func (w *world) Slug() string {
	return slugify.Slugify(w.Name())
}

func (w *world) DBId() string {
	return fmt.Sprintf("world:%v", w.Slug())
}

func (w *world) Metadata() interface{} {
	metadata := bson.D{
		{ "$set", bson.D{
			{
				"name",
				w.Name(),
			},
		}},
	}
	return metadata
}

func (w *world) Save() (*world, error) {
	// TODO: Make configurable
	c := storage.Client.Database("muckity").Collection("worlds")
	f := bson.D{{
		"_id",
		w.DBId(),
	}}
	m := w.Metadata()
	o := newUpsert()
	res, err := c.UpdateOne(context.TODO(), f, m, &o)
	if err != nil {
		panic(err)
	}

	if value, ok := res.UpsertedID.(string); ok {
		if value != w.DBId() {
			return w, errors.New(fmt.Sprintf("Integrity error, got bad ID: %v", value))
		}
	}
	return w, nil
}