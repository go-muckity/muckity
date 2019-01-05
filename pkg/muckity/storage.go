package muckity

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	url2 "net/url"
	"time"
)

type MongoStorage struct {
	id           interface{}
	dbUrl        *url2.URL
	databaseName string
	parentCtx    context.Context
}

var _ MuckitySystem = &MongoStorage{}

// Name implements part of MuckitySystem
func (ms *MongoStorage) Name() string {
	return fmt.Sprintf("mongodb:%v:%v", ms.dbUrl.Host, ms.dbUrl.Path)
}

// Type implements part of MuckitySystem
func (ms *MongoStorage) Type() string {
	return "muckity:storage"
}

func (ms *MongoStorage) Context() context.Context {
	// TODO: utilize context
	return context.TODO()
}

func (w *MongoStorage) String() string {
	return fmt.Sprintf("%v:%v", w.Type(), w.Name())
}

func (ms MongoStorage) Client() (*mongo.Client, error) {
	var (
		ctx    context.Context
		client *mongo.Client
		err    error
	)

	ctx, _ = context.WithTimeout(ms.parentCtx, time.Second*30)
	client, err = mongo.NewClient(ms.dbUrl.String())
	err = client.Connect(ctx)
	return client, err
}

// Save implements storage persistence for compatible objects
func (ms MongoStorage) Save(obj MuckityPersistent) error {
	var (
		client *mongo.Client
		err    error
		coll   *mongo.Collection
	)
	client, err = ms.Client()
	if err != nil {
		// Don't even try
		return err
	}
	// TODO: lookup config object
	if collName, ok := GetConfig().Get("mongodb.collectionRoot").(string); ok {
		coll = client.Database(ms.databaseName).Collection(collName)
	} else {
		coll = client.Database(ms.databaseName).Collection("muckity")
	}
	pd := obj.BSON()
	opt := newUpsert()
	var id interface{}
	id = obj.GetId()
	if id == "" {
		id = primitive.NewObjectID()
	}
	filter := bson.D{{
		"_id",
		id,
	}}
	res, err := coll.UpdateOne(ms.Context(), filter, pd, &opt)
	if err != nil {
		panic(err)
	}
	if uid, ok := res.UpsertedID.(string); ok {
		obj.SetId(uid)
	}
	return err
}

func NewMongoStorage(ctx context.Context) *MongoStorage {
	var (
		ms     MongoStorage
		url    interface{}
		name   interface{}
		dbUrl  *url2.URL
		dbName string
		err    error
	)
	// TODO: lookup config object
	config := GetConfig()
	config.BindEnv("mongodb.url", "MUCKITY_MONGODB_URL")
	config.BindEnv("mongodb.name", "MUCKITY_MONGODB_NAME")
	url = config.Get("mongodb.url")
	name = config.Get("mongodb.name")
	if parse, ok := url.(string); ok {
		dbUrl, err = url2.Parse(parse)
		if err != nil {
			panic(err)
		}
	}
	if name, ok := name.(string); ok {
		dbName = name
	} else {
		panic("mongodb.name or MUCKITY_MONGODB_NAME must be set")
	}
	ms = MongoStorage{nil, dbUrl, dbName, ctx}
	return &ms
}
