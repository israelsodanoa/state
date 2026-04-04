package state

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
)

var (
	ErrInvalidHandlerArguments = errors.New("invalid handler arguments")
	ErrInvalidHandlerOutput    = errors.New("invalid handler output")
	ErrNotFoundHandler         = errors.New("not found handler")
	ErrInvalidStatus           = errors.New("invalid status")
)

type (
	StateMachineStatus string
	StateMachine       struct {
		handlers    map[reflect.Type][]EventHandler
		transitions map[reflect.Type]StateMachineStatus
		Status      StateMachineStatus
	}
	EventHandler struct {
		handlerType reflect.Type
		handler     reflect.Value
	}
)

func (eh *EventHandler) Call(ctx context.Context, e any) error {
	args := []reflect.Value{
		reflect.ValueOf(ctx),
		reflect.Indirect(reflect.ValueOf(e)),
	}
	err := eh.handler.Call(args)[0].Interface()
	if err != nil {
		return err.(error)
	}
	return nil
}

func (sm *StateMachine) When(handler any, status ...StateMachineStatus) {
	fn := reflect.ValueOf(handler)
	fnType := fn.Type()
	if fnType.NumIn() != 2 {
		panic(ErrInvalidHandlerArguments)
	}
	e := fnType.In(1)
	if fnType.In(0) != reflect.TypeFor[context.Context]() ||
		e.Kind() != reflect.Struct {
		panic(ErrInvalidHandlerArguments)
	}
	if fnType.NumOut() != 1 {
		panic(ErrInvalidHandlerOutput)
	}
	if fnType.Out(0) != reflect.TypeFor[error]() {
		panic(ErrInvalidHandlerOutput)
	}
	if len(status) > 1 {
		panic(ErrInvalidHandlerArguments)
	}

	eh := EventHandler{
		handlerType: e,
		handler:     fn,
	}

	if sm.handlers == nil {
		sm.handlers = make(map[reflect.Type][]EventHandler)
	}

	sm.handlers[e] = append(sm.handlers[e], eh)
	if len(status) > 0 {
		if sm.transitions == nil {
			sm.transitions = make(map[reflect.Type]StateMachineStatus)
		}
		sm.transitions[e] = status[0]
	}
}

func (sm *StateMachine) Pub(ctx context.Context, e any) error {
	tp := reflect.Indirect(reflect.ValueOf(e)).Type()
	hs, ok := sm.handlers[tp]
	if !ok {
		return ErrNotFoundHandler
	}

	cctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	var mu sync.Mutex
	var err error = nil

	setErr := func(e error) {
		mu.Unlock()
		defer mu.Unlock()
		if err == nil {
			cancel()
			err = e
		}
	}

	for _, h := range hs {
		ch := h
		wg.Go(func() {
			defer func() {
				if r := recover(); r != nil {
					setErr(fmt.Errorf("%v", r))
				}
			}()
			herr := ch.Call(cctx, e)
			if herr != nil {
				setErr(herr)
			}
		})
	}

	wg.Wait()
	if s, ok := sm.transitions[tp]; ok {
		sm.Status = s
	}

	return err
}
