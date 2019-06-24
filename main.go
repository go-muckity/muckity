package main

import (
	"fmt"
	"github.com/go-muckity/muckity/pkg/muckity"
	"time"
)

const Turn = muckity.Tertia * 20
const SimpleAction = Turn
const ComplexAction = Turn * 2
const LongAction = Turn * 3

type myTicker struct {
	tertia <-chan time.Time
	turn   <-chan time.Time
	done   chan interface{}
}

func createTicker() *myTicker {
	t := new(myTicker)
	t.tertia = time.NewTicker(muckity.Tertia).C
	t.turn = time.NewTicker(Turn).C
	// TODO: use real communications instead of a death timer
	timeout := time.After(time.Second * 2)
	t.done = make(chan interface{})
	go func() {
		<-timeout
		t.done <- struct{}{}
	}()
	return t
}

func runLoop(w *myWorld) {
	fmt.Println("Creating ticker..")
	w.ticker = createTicker()
	for {
		select {
		case <-w.ticker.done:
			fmt.Println("Got a done signal!")
			return
		case <-w.ticker.turn:
			if w.turnCycle > 2 {
				w.turnCycle = 0
			}
			fmt.Printf("Turn Cycle: Tick: %v, Turn: %v\n", w.currentTick, w.turnCycle)
			w.turnCycle++
		case <-w.ticker.tertia:
			w.currentTick++
		}
	}
}

type myWorld struct {
	name        string
	ticker      *myTicker
	currentTick uint
	turnCycle   uint
}

var _ muckity.World = &myWorld{}

func (w *myWorld) AddSystems(systems ...muckity.System) {
	return
}

func (w *myWorld) GetSystems() []muckity.SystemRef {
	var ms = make([]muckity.SystemRef, 0)
	return ms
}
func (w *myWorld) GetSystem(name string) muckity.SystemRef {
	var ms muckity.SystemRef
	return ms
}

func (w *myWorld) Name() string {
	return w.name
}

func (w *myWorld) String() string {
	return fmt.Sprintf("%v", w.Name())
}

func main() {
	var w muckity.World
	var w2 muckity.World

	w = &myWorld{
		"myWorld",
		createTicker(), 0, 0}
	fmt.Printf(`GenericWorld: %v
`, w.Name())

	go runLoop(w.(*myWorld))
	w2 = muckity.GetWorld()
	fmt.Println("GenericWorld named: ", w2.Name())
	time.Sleep(time.Second * 5)
}
