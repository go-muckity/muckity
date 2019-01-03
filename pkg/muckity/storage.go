package muckity

import (
	"context"
	"errors"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"
)

// MuckityStorage implements base storage system for later decoupling
type MuckityStorage struct {
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

func (mc mongoConfig) asURI() string {
	return fmt.Sprintf("mongodb://%v:%v@%v:%v/%v",
		mc.dbUser,
		mc.dbPwd,
		mc.dbHost,
		mc.dbPort,
		mc.dbName)
}

// Name implements part of MuckitySystem
func (ms MuckityStorage) Name() string {
	return fmt.Sprintf("system:mongodb:%v", ms.databaseName)
}

// Type implements part of MuckitySystem
func (ms MuckityStorage) Type() string {
	return "systems"
}

// Channels implements part of MuckitySystem
func (ms MuckityStorage) Channels() map[string]chan interface{} {
	// TODO: implement context handling / closing
	return make(map[string]chan interface{}, 0)
}

// Save implements storage persistence for compatible objects
func (ms *MuckityStorage) Save(obj MuckityObject) (interface{}, error) {
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

func dbInit(done chan interface{}) *MuckityStorage {
	muckCfg := core.Config()
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
	storage := new(MuckityStorage)
	storage.Client = client
	storage.databaseName = muckCfg.mongo.dbName
	return storage
}