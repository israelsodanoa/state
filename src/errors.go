package state

import "errors"

var (
	ErrHandlerAlreadyRegistered = errors.New("Handler already registered")
	ErrInvalidStatus            = errors.New("Invalid StateMachineStatus")
)
