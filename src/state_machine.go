package state

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/ledongthuc/goterators"
)

type (
	StateMachineStatus  string
	StateMachine[T any] struct {
		Mutex    sync.Mutex
		State    T
		Status   StateMachineStatus
		Handlers []EventTransition
	}
)

func (sm *StateMachine[T]) Pub(ctx context.Context, data any) error {
	tp := reflect.Indirect(reflect.ValueOf(data)).Type()
	h, _, err := goterators.Find(sm.Handlers, func(eh EventTransition) bool {
		return eh.EventHandler.EventType == tp
	})
	if err != nil {
		return err
	}
	err = h.EventHandler.Call(ctx, data)
	if err != nil {
		return err
	}

	sm.Mutex.Lock()
	defer sm.Mutex.Unlock()
	sm.Status = h.Status

	return nil
}

func (sm *StateMachine[T]) TransitionTo(s StateMachineStatus, eh EventHandler) {
	sm.Handlers = append(sm.Handlers, EventTransition{EventHandler: eh, Status: s})
}

func When[T any, E any](
	sm *StateMachine[T],
	handler HandlerFn[E]) EventHandler {
	eh := EventHandler{
		EventType: reflect.TypeFor[E](),
		Handler:   reflect.ValueOf(handler),
	}

	for _, h := range sm.Handlers {
		if h.EventHandler.EventType == eh.EventType {
			panic(errors.ErrUnsupported)
		}
	}

	return eh
}
