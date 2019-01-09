package ecs

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type key int

const (
	rootKey    key = iota // root object pointer, of MuckityType
	configKey             // MuckityConfig interface
	worldKey              // MuckityWorld interface
	storageKey            // MuckityStorage interface
	systemKey             // MuckitySystem interface
)

type emptyMktyCtx int

var _ MuckityContext = new(emptyMktyCtx)

func (emptyMktyCtx) Deadline() (deadline time.Time, ok bool) { return }
func (emptyMktyCtx) Done() <-chan struct{}                   { return nil }
func (emptyMktyCtx) Err() error                              { return nil }
func (emptyMktyCtx) Value(key interface{}) interface{}       { return nil }
func (emptyMktyCtx) Root() MuckityType                       { return nil }
func (emptyMktyCtx) Config() MuckityConfig                   { return nil }
func (emptyMktyCtx) World() MuckityWorld                     { return nil }
func (emptyMktyCtx) Storage() MuckityStorage                 { return nil }
func (emptyMktyCtx) CallingSystem() MuckitySystem            { return nil }
func (emptyMktyCtx) Name() string                            { return "emptyContext" }
func (emptyMktyCtx) Type() string                            { return "muckity:context" }

var (
	background = new(emptyMktyCtx)
	todo       = new(emptyMktyCtx)
)

func Background() MuckityContext { return background }
func TODO() MuckityContext       { return todo }

type cancelCtx struct {
	MuckityContext
	done chan struct{}
	err  error
	mu   sync.Mutex
}

var _ MuckityContext = &cancelCtx{}

func (ctx *cancelCtx) Done() <-chan struct{} { return ctx.done }
func (ctx *cancelCtx) Err() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.err
}

var Canceled = errors.New("context canceled")

type CancelFunc func()

func WithCancel(parent MuckityContext) (MuckityContext, CancelFunc) {
	ctx := &cancelCtx{
		parent,
		make(chan struct{}),
		nil,
		sync.Mutex{},
	}

	cancel := func() { ctx.cancel(Canceled) }

	go func() {
		select {
		case <-parent.Done():
			ctx.cancel(parent.Err())
		case <-ctx.Done():

		}
	}()

	return ctx, cancel
}

func (ctx *cancelCtx) cancel(err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.err != nil {
		return
	}

	ctx.err = err
	close(ctx.done)
}

type deadlineCtx struct {
	*cancelCtx
	deadline time.Time
}

var _ MuckityContext = &deadlineCtx{}

func (ctx *deadlineCtx) Deadline() (deadline time.Time, ok bool) {
	return ctx.deadline, true
}

var DeadlineExceeded = errors.New("deadline exceeded")

func WithDeadline(parent MuckityContext, deadline time.Time) (MuckityContext, CancelFunc) {
	cctx, cancel := WithCancel(parent)

	ctx := &deadlineCtx{
		cctx.(*cancelCtx),
		deadline,
	}

	t := time.AfterFunc(time.Until(deadline), func() {
		ctx.cancel(DeadlineExceeded)
	})

	stop := func() {
		t.Stop()
		cancel()
	}

	return ctx, stop
}

func WithTimeout(parent MuckityContext, timeout time.Duration) (MuckityContext, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

type valueCtx struct {
	MuckityContext
	value, key interface{}
}

var _ MuckityContext = &valueCtx{}

func (ctx *valueCtx) Value(key interface{}) interface{} {
	if key == ctx.key {
		return ctx.value
	}
	return ctx.MuckityContext.Value(key)
}

func WithValue(parent MuckityContext, key, value interface{}) MuckityContext {
	if key == nil {
		panic("key is nil")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("key is not comparable")
	}
	return &valueCtx{
		parent,
		value,
		key,
	}
}

func (ctx *valueCtx) Root() MuckityType {
	return ctx.Value(rootKey).(MuckityType)
}

func (ctx *valueCtx) Config() MuckityConfig {
	return ctx.Value(configKey).(MuckityConfig)
}

func (ctx *valueCtx) World() MuckityWorld {
	return ctx.Value(worldKey).(MuckityWorld)
}

func (ctx *valueCtx) Storage() MuckityStorage {
	return ctx.Value(storageKey).(MuckityStorage)
}

func (ctx *valueCtx) CallingSystem() MuckitySystem {
	return ctx.Value(systemKey).(MuckitySystem)
}

func WithRoot(parent MuckityContext, value MuckityType) MuckityContext {
	return WithValue(parent, rootKey, value)
}

func WithConfig(parent MuckityContext, value MuckityConfig) MuckityContext {
	return WithValue(parent, configKey, value)
}

func WithWorld(parent MuckityContext, value MuckityWorld) MuckityContext {
	return WithValue(parent, worldKey, value)
}

func WithStorage(parent MuckityContext, value MuckityStorage) MuckityContext {
	return WithValue(parent, storageKey, value)
}

func WithSystem(parent MuckityContext, value MuckitySystem) MuckityContext {
	return WithValue(parent, systemKey, value)
}

var rootContext MuckityContext

func newContext(ctx ...interface{}) MuckityContext {
	var mC MuckityContext
	if len(ctx) == 0 {
		return rootContext
	}
	if len(ctx) == 1 {
		mC = ctx[0].(MuckityContext)
		return WithRoot(TODO(), mC)
	}
	return WithRoot(todo, background)
}

// GetContext returns a singleton context object, or an empty root context pointing at
// a default, unexported Background() context. Pass doOnce: true if you want the singleton
func GetContext(doOnce bool, ctx ...interface{}) MuckityContext {
	if doOnce {
		once.Do(func() {
			rootContext = newContext(ctx...)
		})
	}
	return newContext(ctx...)
}
