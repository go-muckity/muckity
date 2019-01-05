package muckity

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	url2 "net/url"
	"os"
	"strings"
	"time"
)

type MongoStorage struct {
	id				interface{}
	dbUrl			*url2.URL
	databaseName 	string
	parentCtx 	context.Context
}

// Name implements part of MuckitySystem
func (ms *MongoStorage) Name() string {
	return fmt.Sprintf("mongodb:%v:%v", ms.dbUrl.Host, ms.dbUrl.Path)
}

// Type implements part of MuckitySystem
func (ms *MongoStorage) Type() string {
	return "muckity:storage"
}

func (ms *MongoStorage) Context() context.Context {
	return ms.parentCtx
}

func (w *MongoStorage) String() string {
	return fmt.Sprintf("%v:%v", w.Type(), w.Name())
}

func (ms MongoStorage) Client() (*mongo.Client, error)  {
	var (
		ctx		context.Context
		client	*mongo.Client
		err		error
	)

	ctx, _ = context.WithTimeout(ms.parentCtx, time.Second * 30)
	client, err = mongo.NewClient(ms.dbUrl.String())
	err = client.Connect(ctx)
	return client, err
}

// Save implements storage persistence for compatible objects
func (ms MongoStorage) Save(obj MuckityPersistent) error {
	var (
		client	*mongo.Client
		err 	error
	)
	client, err = ms.Client()
	if err != nil {
		// Don't even try
		return err
	}
	coll := client.Database(ms.databaseName).Collection("muckity")
	pd := obj.BSON()
	opt := newUpsert()
	var id interface {}
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
		ms 	MongoStorage
		url	*url2.URL
	)
	url = parseConfig()
	pathSplit := strings.Split(url.Path, "/")
	path := pathSplit[len(pathSplit)-1]
	ms = MongoStorage{nil, url, path, ctx }
	return &ms
}

// TODO: use a real config service
func parseConfig() *url2.URL {
	var (
		dbUser	string
		dbPwd	string
		dbHost	string
		dbName	string
		url 	*url2.URL
	)

	if value, ok := os.LookupEnv("MUCKITY_DB_USERNAME"); ok {
		dbUser = value
	} else {
		dbUser = "muckity"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_PWD"); ok {
		dbPwd = value
	} else {
		dbPwd = "muckity"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_HOST"); ok {
		dbHost = value
	} else {
		dbHost = "localhost"
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_PORT"); ok {
		dbHost = fmt.Sprintf("%v:%v", dbHost, value)
	} else {
		dbHost = fmt.Sprintf("%v:27017", dbHost)
	}
	if value, ok := os.LookupEnv("MUCKITY_DB_NAME"); ok {
		dbName = value
	} else {
		dbName = "muckity"
	}
	url, err := url2.Parse(fmt.Sprintf("mongodb://%v:%v@%v/%v", dbUser, dbPwd, dbHost, dbName))
	if err != nil {
		panic(err)
	}
	return url
}

