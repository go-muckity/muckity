package ecs

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	url2 "net/url"
	"time"
)

type muckityStorage struct {
	name string
	path string
}

func (ms *muckityStorage) Name() string {
	return fmt.Sprintf( "%v", ms.name)
}

func (ms *muckityStorage) Type() string {
	return fmt.Sprintf("muckity:muckityStorage")
}

func (ms *muckityStorage) Context() MuckityContext {
	return TODO()
}

func (ms muckityStorage) Save(obj MuckityPersistent) error {
	return nil
}

var _ MuckityStorage = &muckityStorage{}

type MongoStorage struct {
	id           interface{}
	dbUrl        *url2.URL
	databaseName string
	parentCtx    MuckityContext // TODO pass this to Run()
}

var _ MuckitySystem = &MongoStorage{}
var _ MuckityStorage = &MongoStorage{}

// Name implements part of MuckitySystem
func (ms *MongoStorage) Name() string {
	return fmt.Sprintf("mongodb:%v:%v%v", ms.dbUrl.Host, ms.dbUrl.Path, ms.databaseName)
}

// Type implements part of MuckitySystem
func (ms *MongoStorage) Type() string {
	return "muckity:storage"
}

func (ms *MongoStorage) Context() MuckityContext {
	// TODO: utilize context
	return TODO()
}

func (ms *MongoStorage) String() string {
	return fmt.Sprintf("%v:%v", ms.Type(), ms.Name())
}

func (ms MongoStorage) Client() (*mongo.Client, error) {
	var (
		ctx    MuckityContext
		client *mongo.Client
		err    error
	)

	ctx, _ = WithTimeout(ms.parentCtx, time.Second*30)
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

func NewMongoStorage(ctx MuckityContext) *MongoStorage {
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

// GetStorage creates a configured and contextualized storage object. It takes any number of parameters, but currently
// only accepts 0 or 1.  With 0 it will return the default storage object (until 0.1.0, MongoStorage).
// With 1, it expects either a MuckityWorld (a parent that has a system named "storage" on it), or a MuckityContext.
// TODO: Implement passing the system name with the World object, to allow for another storage name.
// TODO: Implement passing a MuckityStorage object
// TODO: Implement bool for Once()
func GetStorage(ctx ...interface{}) MuckityStorage {
	var ms MuckityStorage
	if len(ctx) == 0 {
		return NewMongoStorage(TODO())
	}
	if len(ctx) == 1 {
		switch v := ctx[0].(type) {
		case MuckityWorld:
			sRef := v.GetSystem("storage")
			system := sRef.GetSystem()
			if system, ok := system.(MuckityStorage); ok {
				ms = system
			}
		case MuckityContext:
			ms = NewMongoStorage(v)
		default:
			panic(fmt.Sprintf("Unimplemented context for GetStorage(): %v", v))
		}
		return ms
	}
	// TODO: Implement system struct ptr + ctx; needs SetContext() added to MuckitySystem first
	// Should be something like a filter chain:
	// cast to MuckitySystem
	// 		cast to MuckityStorage
	// 			SetContext()
	// 			return
	panic("Too many arguments passed to GetStorage()")
}
