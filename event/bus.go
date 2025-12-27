package event

import (
	"reflect"
	"sync"
)

type handler struct {
	callback reflect.Value

	// if true, will the handler be executed once
	once bool
	// add async
}
type BusOnce interface {
	// the handler will be executed only once, can be unsubscribed
	SubscribeOnce(id string, fn any) error
}
type Bus interface {
	Subscribe(id string, fn any) error
	// Returns error if there are no callbacks subscribed to the topic.
	Unsubscribe(id string, fn any) error
	// Publish/Send a new event, where
	// id is a topic or a event name
	Emit(id string, args ...any)
	BusOnce
}

type EventBus struct {
	subs map[string][]handler
	sync.Mutex
}

func New() Bus {
	return &EventBus{
		subs:  map[string][]handler{},
		Mutex: sync.Mutex{},
	}
}

func (b *EventBus) onSubscribe(id string, fn any, handler handler) error {
	if reflect.ValueOf(fn).Kind() != reflect.Func {
		return ErrNotAFunc
	}
	b.Lock()
	defer b.Unlock()
	b.subs[id] = append(b.subs[id], handler)

	return nil
}

func (b *EventBus) argsToValues(args ...any) []reflect.Value {
	values := make([]reflect.Value, len(args))
	for i, a := range args {
		values[i] = reflect.ValueOf(a)
	}

	return values
}

func (b *EventBus) prepareArguments(hand handler, args ...any) []reflect.Value {
	vals := b.argsToValues(args...)
	expectedParams := hand.callback.Type().NumIn()

	if expectedParams > len(args) {
		missingPar := make([]reflect.Value, expectedParams-len(args))
		for j := 0; j < len(missingPar); j++ {
			missingPar[j] = reflect.Zero(hand.callback.Type().In(len(args) + j))
		}
		return append(vals, missingPar...)
	} else if expectedParams < len(args) {
		return vals[0:expectedParams]
	}
	return vals
}

func (b *EventBus) callHandler(hand handler, args ...any) {
	vals := b.prepareArguments(hand, args...)
	hand.callback.Call(vals)
}

func (b *EventBus) removeHandler(id string, index int) {
	b.subs[id] = append(b.subs[id][:index], b.subs[id][index+1:]...)
}

func (b *EventBus) Emit(id string, args ...any) {
	if len(b.subs) == 0 || len(b.subs[id]) == 0 {
		return
	}

	b.Lock()
	handlers := make([]handler, len(b.subs[id]))
	copy(handlers, b.subs[id])
	b.Unlock()

	for i := len(handlers) - 1; i >= 0; i-- {
		hand := handlers[i]
		b.callHandler(hand, args...)

		if hand.once {
			b.Lock()
			b.removeHandler(id, i)
			b.Unlock()
		}
	}
}

func (b *EventBus) Subscribe(id string, fn any) error {
	return b.onSubscribe(id, fn, handler{
		callback: reflect.ValueOf(fn),
		once:     false,
	})
}

func (b *EventBus) Unsubscribe(id string, fn any) error {
	b.Lock()
	defer b.Unlock()

	handlers := b.subs[id]
	if len(handlers) == 0 {
		return ErrNoHandlers
	}
	for i, h := range handlers {
		val := reflect.ValueOf(fn)
		if h.callback.Type() == val.Type() &&
			h.callback.Pointer() == val.Pointer() {
			b.removeHandler(id, i)
		}
	}

	return nil
}

func (b *EventBus) SubscribeOnce(id string, fn any) error {
	return b.onSubscribe(id, fn, handler{
		callback: reflect.ValueOf(fn),
		once:     true,
	})
}
