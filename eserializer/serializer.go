// Package eserializer provides EventSerializer interface that can be implemented by any serializer.
package eserializer

// EventSerializer provide methods for serializing/deserializing Event structs.
// If you want to create your own serializer - just implement this interface.
type EventSerializer interface {
	Serialize(Event) ([]byte, error)
	Deserialize([]byte, Event) error
}
