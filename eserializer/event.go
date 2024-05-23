package eserializer

import "errors"

// Event represents event type that should be serialized by event serializer.
// If you want some struct to become an event - just inject this interface inside struct.
//
// Example:
//
//	type SomeEvent struct {
//	     eserializer.Event
//	     SomeID uuid.UUID `json:"some-id"`
//	     ...
//	}
type Event interface {
	event()
}

var ErrEventCantBeNil = errors.New("event cant be nil")
