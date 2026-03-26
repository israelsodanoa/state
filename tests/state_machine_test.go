package tests

import (
	"context"
	state "state/src"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	DummyState struct {
		Counter int
	}

	DummyEvent struct {
		Version int
	}
)

const (
	DummyStatus state.StateMachineStatus = "running"
)

func Test_Should_Pub_With_Sucess(t *testing.T) {
	var sm state.StateMachine[DummyState]
	sm.TransitionTo(DummyStatus,
		state.When(&sm, func(ctx context.Context, event DummyEvent) error {
			sm.Mutex.Lock()
			defer sm.Mutex.Unlock()
			sm.State.Counter += event.Version

			return nil
		}))

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			sm.Pub(t.Context(), DummyEvent{Version: i})
		})
	}

	wg.Wait()
	assert.Equal(t, 4950, sm.State.Counter, "they should be equal")
	assert.Equal(t, sm.Status, DummyStatus, "they should be equal")
}

func Test_Should_Compose_With_Sucess(t *testing.T) {
	var sm state.StateMachine[DummyState]
	sm.TransitionTo(DummyStatus,
		state.Compose(
			state.When(&sm, func(ctx context.Context, event DummyEvent) error {
				sm.Mutex.Lock()
				defer sm.Mutex.Unlock()
				sm.State.Counter += event.Version

				return nil
			}),
			state.When(&sm, func(ctx context.Context, event DummyEvent) error {
				sm.Mutex.Lock()
				defer sm.Mutex.Unlock()
				sm.State.Counter += event.Version

				return nil
			}),
			state.When(&sm, func(ctx context.Context, event DummyEvent) error {
				sm.Mutex.Lock()
				defer sm.Mutex.Unlock()
				sm.State.Counter += event.Version

				return nil
			})))

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			sm.Pub(t.Context(), DummyEvent{Version: i})
		})
	}

	wg.Wait()
	assert.Equal(t, 4950*3, sm.State.Counter, "they should be equal")
	assert.Equal(t, sm.Status, DummyStatus, "they should be equal")
}
