package event

import "errors"

// Event represents event type that should be serialized by event serializer.
// If you want some struct to become an event - just inject this interface inside struct.
//
// Example:
//
//	type SomeEvent struct {
//	     event.Event
//	     SomeID uuid.UUID `json:"some-id"`
//	     ...
//	}
type Event interface {
	event()
}

// Marshaler provide methods for serializing/deserializing Event structs.
type Marshaler interface {
	Marshal(Event) ([]byte, error)
}

type Unmarshaler interface {
	Unmarshal([]byte, Event) error
}

type MarshalUnmarshaler interface {
	Marshaler
	Unmarshaler
}

var ErrEventCantBeNil = errors.New("event cant be nil")
