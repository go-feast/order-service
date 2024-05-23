// Package eserializer provide serializers and deserializers for Event.
package eserializer

// EventSerializer provide methods for serializing/deserializing Event structs.
type EventSerializer interface {
	Serialize(Event) ([]byte, error)
	Deserialize([]byte, Event) error
}
