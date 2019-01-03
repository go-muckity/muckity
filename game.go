package main

import (
	"errors"
	"fmt"
	"github.com/machiel/slugify"
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
	timeout := time.After(time.Second * 20)
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
	name        string
	description string
	zones       []string
	ticker 		*myTicker
	currentTick uint
	turnCycle	uint
}

func (w myWorld) Name() string {
	return w.name
}

func (w myWorld) Type() string {
	return "worlds"
}

func (w myWorld) Aliases() []string {
	s := make([]string, 0)
	s = append(s, "something")
	return s
}

func (w myWorld) DBId() string {
	return fmt.Sprintf("world:%v", slugify.Slugify(w.name))
}

func (w myWorld) Metadata() interface{} {
	return w.PersistentData()
}

func (w myWorld) PersistentData() interface{} {
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
			{
				"aliases",
				w.Aliases(),
			},
		}},
	}
	return p
}

func main() {
	core := muckity.Init()
	w := new(myWorld)
	w.name = "Descriptive, aliased, world"
	w.description = "I am a test world created for integration testing of the muckity package."
	w.zones = make([]string, 0)
	store, err := core.GetSystem("storage")
	if err != nil {
		panic(err)
	}
	if store, ok := store.(*muckity.MuckityStorage); ok {
		_, err := store.Save(w)
		if err != nil {
			panic(err)
		}
	} else {
		panic(errors.New(fmt.Sprintf("Storage error, uknown storage object: %v", store)))
	}
	runLoop(w)
}
