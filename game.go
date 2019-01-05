package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/tsal/muckity/pkg/muckity"
	"time"
)

const Turn = muckity.Tertia * 20
const SimpleAction = Turn
const ComplexAction = Turn * 2
const LongAction = Turn * 3

type myTicker struct {
	tertia <- chan time.Time
	turn <- chan time.Time
	done chan interface{}
}

func createTicker() *myTicker {
	t := new(myTicker)
	t.tertia = time.NewTicker(muckity.Tertia).C
	t.turn = time.NewTicker(Turn).C
	// TODO: use real communications instead of a death timer
	timeout := time.After(time.Second * 10)
	t.done = make(chan interface{})
	go func() {
		<-timeout
		t.done <- struct{}{}
	}()
	return t
}

func runLoop(w *myWorld) error {
	w.ticker = createTicker()
	for {
		select {
		case <- w.ticker.done:
			fmt.Println("Got a done signal!")
			return nil
		case <- w.ticker.turn:
			if w.turnCycle > 2 {
				w.turnCycle = 0
			}
			fmt.Printf("Turn Cycle: Tick: %v, Turn: %v\n", w.currentTick, w.turnCycle)
			w.turnCycle++
		case <- w.ticker.tertia:
			w.currentTick++
		}
	}
	return nil
}

type myWorld struct {
	id					string
	name				string
	myContext			context.Context
	description			string
	zones       		[]string
	ticker 				*myTicker
	currentTick 		uint
	turnCycle			uint
}

func (w *myWorld) Name() string {
	return w.name
}

func (w *myWorld) Type() string {
	return "game:world"
}

func (w *myWorld) Context() context.Context {
	return w.myContext
}

func (w *myWorld) String() string {
	return fmt.Sprintf("%v:%v", w.Type(), w.Name())
}

func (w *myWorld) BSON() interface{} {
	p := bson.D{
		{"$set", bson.D{
			{
				"name",
				w.Name(),
			},
			{
				"description",
				w.description,
			},
			{
				"zones",
				w.zones,
			},
		}},
	}
	return p
}

func (w *myWorld) GetId() string {
	// TODO: needs better checking
	return w.id
}

func (w *myWorld) SetId(id string) {
	w.id = id
}
func main() {
	var (
		w interface{}
		storage muckity.MuckityStorage
	)

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()
	w = &myWorld{
		"world:descriptive-world",
		"world",
		ctx,
		"I am a test world created for integration testing of the muckity package.",
		make([]string, 0),
		createTicker(),
		0,
		0 }

	fmt.Println("Created World:", w)
	storage = muckity.NewMongoStorage(ctx)
	fmt.Println("Created Storage:", storage)
	var pers muckity.MuckityPersistent
	pers = w.(muckity.MuckityPersistent)
	storage.Save(pers)
	runLoop(w.(*myWorld))
}