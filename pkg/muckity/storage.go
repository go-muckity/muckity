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
	return &MongoStorage{}
}

func (ms *MongoStorage) SetClient(url string) *MongoStorage {
	client, err := mongo.NewClient(url)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	ms.Client = client
	return ms
}

// TODO: Make this configurable and pooled
var storage = NewMongoStorage()

func NewMuckityStorage(url string) *MongoStorage {
	return storage.SetClient(url)
}