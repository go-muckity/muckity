package muckity

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
)

type DefaultMuckity struct {
	id              uuid.UUID
	Name            string
	messagesChannel chan Message
	closerFunc      func()
	handlerFunc     HandlerFunc
	rootFuncs       []RootFunc
	systemMap       SystemMap
}

func (d DefaultMuckity) UUID() uuid.UUID {
	return d.id
}
func (d DefaultMuckity) MessageChannel() chan<- Message {
	return d.messagesChannel
}
func (d DefaultMuckity) String() string {
	return d.id.String()
}

type defaultMuckityJSON struct {
	Name string `json:"name"`
	UUID string `json:"uuid,omitempty"`
}

func (d DefaultMuckity) MarshalJSON() ([]byte, error) {
	jsonMuckity := defaultMuckityJSON{
		Name: d.Name,
		UUID: d.id.String(),
	}
	return json.Marshal(jsonMuckity)
}
func (d *DefaultMuckity) Closer() func() {
	if d.closerFunc != nil {
		return d.closerFunc
	}
	return func() {
		d.messagesChannel <- struct {
			Message string
			Close   bool
		}{Message: "closing channels", Close: true}
		defer close(d.messagesChannel)
		return
	}
}
func (d *DefaultMuckity) Systems() SystemMap {
	return d.systemMap
}
func (d *DefaultMuckity) Handler(ctx context.Context, m Message) (Message, error) {
	if d.handlerFunc != nil {
		return d.handlerFunc(ctx, m)
	}
	return m, nil
}
func (d *DefaultMuckity) Init(cfg InitConfig, rootFuncs ...RootFunc) error {
	d.id = uuid.New()
	d.Name = cfg.Name
	d.messagesChannel = cfg.MessageChannel
	if d.messagesChannel == nil {
		d.messagesChannel = make(chan Message)
	}
	d.closerFunc = cfg.CloseFunc
	d.handlerFunc = cfg.HandlerFunc
	d.rootFuncs = make([]RootFunc, len(rootFuncs))
	for i, fn := range rootFuncs {
		d.rootFuncs[i] = fn
		err := fn(d)
		if err != nil {
			return err
		}
	}
	return nil
}
func (d *DefaultMuckity) UnmarshalJSON(data []byte) error {
	var jsonMuckity = new(defaultMuckityJSON)
	err := json.Unmarshal(data, jsonMuckity)
	if err != nil {
		return err
	}
	d.Name = jsonMuckity.Name
	id, err := uuid.Parse(jsonMuckity.UUID)
	if err != nil {
		return err
	}
	d.id = id
	d.messagesChannel = make(chan Message)
	return nil
}

var _ Muckity = new(DefaultMuckity)
