package muckity

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

type MongoStorage struct {
	Client *mongo.Client
}

type Sluggified interface {
	Slug() string
}

func (ms *MongoStorage) IsAvailable() bool {
	return true
}

func NewMongoStorage() *MongoStorage {
	client, err := mongo.NewClient(`mongodb://muckity:muckity@localhost:27017/muckity`)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	mongo := new(MongoStorage)
	mongo.Client = client
	return mongo
}

// TODO: Make this configurable
var storage = NewMongoStorage()
