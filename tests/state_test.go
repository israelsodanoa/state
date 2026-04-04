package tests

import (
	"testing"

	"github.com/israelsodanoa/state"
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

}

func Test_Should_Compose_With_Sucess(t *testing.T) {

}
