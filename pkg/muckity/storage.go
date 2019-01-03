package muckity

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

// MongoStorage implements base storage system for later decoupling
type MongoStorage struct {
	Client       *mongo.Client
	databaseName string
}

// Persistent interface describes objects / structs that can be saved in storage
// DBId() should return a string - empty if new; _id key otherwise
// PersistentData() should return a bson.D{{}} struct or marshal data
type Persistent interface {
	DBId() string
	PersistentData() interface{}
}

type mongoConfig struct {
	dbUser string
	dbPwd  string
	dbHost string
	dbPort uint
	dbName string
}

func (gc mongoConfig) asURI() string {
	return fmt.Sprintf("mongodb://%v:%v@%v:%v/%v",
		gc.dbUser,
		gc.dbPwd,
		gc.dbHost,
		gc.dbPort,
		gc.dbName)
}

// Save implements storage persistence for compatible objects
func (ms *MongoStorage) Save(obj MuckityObject) (interface{}, error) {
	if pObj, ok := obj.(Persistent); ok {
		coll := ms.Client.Database(ms.databaseName).Collection(obj.Type())
		pd := pObj.PersistentData()
		opt := newUpsert()
		filter := bson.D{{
			"_id",
			pObj.DBId(),
		}}
		res, err := coll.UpdateOne(context.TODO(), filter, pd, &opt)
		if err != nil {
			panic(err)
		}
		if uid, ok := res.UpsertedID.(string); ok {
			if uid != pObj.DBId() {
				return obj, errors.New(fmt.Sprintf("Integrity error, got bad ID: %v", uid))
			}
		}
	} else {
		return obj, errors.New(fmt.Sprintf("Tried to persist a non-persistent object: %v", obj.Name()))
	}
	return obj, nil
}

var storage *MongoStorage

func init() {
	muckCfg, err := GetMuckityConfig()
	if err != nil {
		panic(err)
	}
	parseConfig(muckCfg)
	client, err := mongo.NewClient(muckCfg.mongo.asURI())
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	storage = new(MongoStorage)
	storage.Client = client
	storage.databaseName = muckCfg.mongo.dbName
}

// GetMuckityStorage is a helper function for retrieving the storage system
func GetMuckityStorage() *MongoStorage {
	return storage
}
