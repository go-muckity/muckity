package muckity

import (
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

var _ System = &GenericSystem{}

// GenericSystem is the default implementation of System
type GenericSystem struct {
	id   uuid.UUID
	name string
}

func (s GenericSystem) String() string             { return s.name }
func (s GenericSystem) UUID() uuid.UUID            { return s.id }
func (s *GenericSystem) Run() (int, error)         { panic("implement me") }
func (s *GenericSystem) Next() <-chan SystemFunc   { panic("implement me") }
func (s *GenericSystem) Init(...interface{}) error { panic("implement me") }
func (s *GenericSystem) Update(System) error       { panic("implement me") }
func (s *GenericSystem) Shutdown()                 { panic("implement me") }

var _ WorldSystem = &GenericWorld{}

// GenericWorld is the default implementation of WorldSystem
type GenericWorld struct {
	*GenericSystem
	ticker      Ticker
	currentTick int
	tickMax     int
	waitTicker  sync.WaitGroup
}

func (w *GenericWorld) Join(interface{}) error {
	return nil
}

var TickNotImplemented = fmt.Errorf("tick function not implemented for this world; this can probably be ignored")

func (w *GenericWorld) Tick() error {
	if w.currentTick%10 == 0 {
		fmt.Printf("world update: %s - TICK: %04d\n", w.Name(), w.currentTick)
	}
	if w.currentTick >= w.tickMax {
		w.Shutdown()
		return nil
	}
	w.currentTick++
	return nil
}
func (w *GenericWorld) Init(opts ...interface{}) error {
	var err error
	var haveTicker = false
	w.tickMax = 1000
	w.currentTick = 0
	for _, opt := range opts {
		switch opt.(type) {
		case Ticker:
			if haveTicker {
				return fmt.Errorf("attempted to assign multiple tickers")
			}
			w.ticker = opt.(Ticker)
			err = w.ticker.Init(opts...)
			if err != nil {
				return err
			}
			haveTicker = true
		case TargetedSystem:
			err = opt.(TargetedSystem).Target(w)
		case func(system WorldSystem) error:
			err = opt.(func(system WorldSystem) error)(w)
		case int:
			w.tickMax = opt.(int)
		}
		if err != nil {
			return err
		}
	}
	if !haveTicker {
		w.ticker = new(GenericTicker)
		_ = w.ticker.Target(w)
		err = w.ticker.Init(opts)
		if err != nil {
			return err
		}
	}
	return err
}
func (w *GenericWorld) Run() (int, error) {
	var err error
	fmt.Println("world starting:", w.Name())
	w.waitTicker.Add(1)
	var tRunner = func() {
		defer w.waitTicker.Done()
		_, err = w.ticker.Run()
	}
	go tRunner()
	w.waitTicker.Wait()
	return w.currentTick, err
}
func (w *GenericWorld) Shutdown() {
	w.ticker.Shutdown()
}
func (w GenericWorld) Name() string { return w.name }
func (w GenericWorld) String() string {
	return fmt.Sprintf("%s:%s", w.Name(), w.UUID().String())
}
func GetWorld() WorldSystem {
	var (
		world WorldSystem
	)
	world = &GenericWorld{
		GenericSystem: &GenericSystem{id: uuid.New(), name: "generic-world"},
		ticker:        nil,
	}
	return world
}

const Turn = Tertia * 20

type GenericTicker struct {
	tertia  <-chan time.Time
	targets map[uuid.UUID]TickingSystem
	close   chan interface{}
}

var _ Ticker = new(GenericTicker)

var InvalidGenericTickerTarget = fmt.Errorf("GenericTicker can only target TickingSystem Systems")

func (t *GenericTicker) Target(target interface{}) error {
	if t.targets == nil {
		t.targets = make(map[uuid.UUID]TickingSystem)
	}
	if v, ok := target.(System); ok {
		if w, ok := v.(TickingSystem); ok {
			t.targets[v.UUID()] = w
			return nil
		} else {
			return InvalidGenericTickerTarget
		}
	} else {
		return InvalidGenericTickerTarget
	}
}
func (t *GenericTicker) Targets() []interface{} {
	var targets []interface{}
	targets = make([]interface{}, 0)
	for _, v := range t.targets {
		targets = append(targets, v)
	}
	return targets
}
func (t *GenericTicker) Untarget(target interface{}) error {
	if v, ok := target.(System); ok {
		if _, ok := v.(TickingSystem); ok {
			delete(t.targets, v.UUID())
			return nil
		} else {
			return InvalidGenericTickerTarget
		}
	} else {
		return InvalidGenericTickerTarget
	}
}
func (t *GenericTicker) Init(...interface{}) error {
	var err error
	t.tertia = time.NewTicker(Tertia).C
	t.close = make(chan interface{})
	return err
}

func tickTarget(t TickingSystem) error {
	return t.Tick()
}

func (t *GenericTicker) tickLoop() (int, error) {
	var count = 0
	for id := range t.targets {
		err := tickTarget(t.targets[id])
		if err != nil {
			return -1, err
		}
		count++
	}
	return count, nil
}

func (t *GenericTicker) Run() (int, error) {
	var count = 0
	for {
		select {
		case <-t.close:
			return -2, nil
		case <-t.tertia:
			c, err := t.tickLoop()
			if err != nil {
				return c, err
			}
			count += c
		}
	}
}

func (t *GenericTicker) Shutdown() {
	close(t.close)
}

func (t GenericTicker) String() string { return fmt.Sprintf("ticker:%s", t.String()) }

func (t GenericTicker) Rate() time.Duration {
	return Tertia
}

func (t GenericTicker) UUID() uuid.UUID {
	panic("implement me")
}

func (t GenericTicker) Next() <-chan SystemFunc {
	panic("implement me")
}

func (t GenericTicker) Update(System) error {
	panic("implement me")
}
