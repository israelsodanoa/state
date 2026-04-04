package tests

import (
	"context"
	"sync"
	"testing"

	"github.com/israelsodanoa/state"
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
	sm := state.NewStateMachine()
	var st DummyState

	state.AddState(sm, &st)
	sm.When(func(ctx context.Context, event DummyEvent) error {
		s := state.GetState[DummyState](sm)
		s.Counter += event.Version

		return nil
	}, DummyStatus)

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			sm.Pub(t.Context(), DummyEvent{Version: i})
		})
	}

	wg.Wait()
	assert.Equal(t, 4950, st.Counter, "they should be equal")
	assert.Equal(t, sm.Status, DummyStatus, "they should be equal")
}

func Test_Should_Compose_With_Sucess(t *testing.T) {
	sm := state.NewStateMachine()
	var st DummyState

	state.AddState(sm, &st)
	sm.When(func(ctx context.Context, event DummyEvent) error {
		s := state.GetState[DummyState](sm)
		s.Counter += event.Version

		return nil
	}, DummyStatus)

	sm.When(func(ctx context.Context, event DummyEvent) error {
		s := state.GetState[DummyState](sm)
		s.Counter += event.Version

		return nil
	})

	sm.When(func(ctx context.Context, event DummyEvent) error {
		s := state.GetState[DummyState](sm)
		s.Counter += event.Version

		return nil
	})

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Go(func() {
			sm.Pub(t.Context(), DummyEvent{Version: i})
		})
	}

	wg.Wait()
	assert.Equal(t, 4950*3, st.Counter, "they should be equal")
	assert.Equal(t, sm.Status, DummyStatus, "they should be equal")
}
