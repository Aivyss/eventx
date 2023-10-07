package errors

import (
	"errors"
)

type ErrorID int
type Error struct {
	error
	ErrorID
}

const (
	AlreadyRegistered ErrorID = iota
	NotFoundEventListener
)

var (
	AlreadyRegisteredErr = Error{
		error:   errors.New("AlreadyRegistered"),
		ErrorID: AlreadyRegistered,
	}
	NotFoundEventListenerErr = Error{
		error:   errors.New("NotFoundEventListener"),
		ErrorID: NotFoundEventListener,
	}
)
