package codec

import (
	"encoding/json"
	"errors"
	"reflect"
)

var (
	JSONMarshalService = jsonMarshaller{}
	RawMarshalService  = rawMarshaller{
		elementType: reflect.TypeOf([]byte{}).Name(),
	}
)

type (
	jsonMarshaller struct{}
)

func (jsonMarshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (jsonMarshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type (
	rawMarshaller struct {
		elementType string
	}
)

func (m rawMarshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (m rawMarshaller) Unmarshal(data []byte, v interface{}) error {
	_, isBytes := v.(*[]byte)
	if isBytes {
		*v.(*[]byte) = data
		return nil
	}
	return errors.New("the type of object is not []byte")
}
