package eserializer

import "encoding/json"

type JSONSerializer struct{}

func (j JSONSerializer) Serialize(event Event) ([]byte, error) {
	if event == nil {
		return nil, ErrEventCantBeNil
	}

	return json.Marshal(event)
}

func (j JSONSerializer) Deserialize(bytes []byte, event Event) error {
	if event == nil {
		return ErrEventCantBeNil
	}

	return json.Unmarshal(bytes, event)
}
