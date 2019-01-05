package muckity

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
)

type World struct {
	id					interface{}
	name				string
	mType 				string
	parentCtx			context.Context
}

func (w *World) Name() string {
	return w.name
}

func (w *World) Type() string {
	return w.mType
}

func (w *World) Context() context.Context {
	return w.parentCtx
}

func (w *World) String() string {
	return fmt.Sprintf("%v:%v", w.Type(), w.Name())
}

func NewWorld() *World {
	return &World{ nil,"name", "muckity:world", context.Background() }
}

func (w *World) BSON() interface{} {
	p := bson.D{
		{"$set", bson.D{
			{
				"name",
				w.Name(),
			},
			{
				"type",
				w.Type(),
			},
		}},
	}
	return p
}

func (w *World) GetId() string {
	// TODO: needs better checking
	return w.id.(fmt.Stringer).String()
}

func (w *World) SetId(id string) {
	w.id = id
}