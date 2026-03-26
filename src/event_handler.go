package state

import (
	"context"
	"reflect"
	"sync"
)

type (
	EventHandler struct {
		EventType   reflect.Type
		Handler     reflect.Value
		ValidStatus []StateMachineStatus
	}
	EventTransition struct {
		EventHandler EventHandler
		Status       StateMachineStatus
	}
	HandlerFn[E any] func(context.Context, E) error
)

func (e *EventHandler) Call(ctx context.Context, data any) error {
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.Indirect(reflect.ValueOf(data)),
	}
	err := e.Handler.Call(args)[0].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}

func Compose(handlers ...EventHandler) EventHandler {
	wrapper := func(ctx context.Context, data any) error {
		var wg sync.WaitGroup
		var err error
		var mutex sync.Mutex
		cctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for _, h := range handlers {
			wg.Go(func() {
				inner_err := h.Call(cctx, data)
				mutex.Lock()
				defer mutex.Unlock()
				if inner_err != nil && err == nil {
					cancel()
					err = inner_err
				}
			})
		}

		wg.Wait()
		return err
	}

	eh := EventHandler{
		EventType: handlers[0].EventType,
		Handler:   reflect.ValueOf(wrapper),
	}

	return eh
}
