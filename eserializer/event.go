package eserializer

import "errors"

type Event interface {
	ToEvent() Event
}

var ErrEventCantBeNil = errors.New("event cant be nil")
