package event

import "encoding/json"

var _ MarshalUnmarshaler = JSONMarshaler{}

type JSONMarshaler struct{}

func (j JSONMarshaler) Marshal(event Event) ([]byte, error) {
	if event == nil {
		return nil, ErrEventCantBeNil
	}

	return json.Marshal(event)
}

func (j JSONMarshaler) Unmarshal(bytes []byte, event Event) error {
	if event == nil {
		return ErrEventCantBeNil
	}

	return json.Unmarshal(bytes, event)
}
