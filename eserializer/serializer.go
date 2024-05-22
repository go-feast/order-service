// Package eserializer provide serializers and deserializers for Event.
package eserializer

type SerializeDeserializer interface {
	Serialize(Event) ([]byte, error)
	Deserialize([]byte, Event) error
}
