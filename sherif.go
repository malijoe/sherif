package sherif

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
)

type Decoder interface {
	Decode(obj any) error
}

type DecoderFunc func(obj any) error

func (f DecoderFunc) Decode(i any) error {
	return f(i)
}

func JSONDecoderFunc(data []byte) DecoderFunc {
	return func(obj any) error {
		return json.Unmarshal(data, obj)
	}
}

func BSONDecoderFunc(data []byte) DecoderFunc {
	return func(obj any) error {
		return bson.Unmarshal(data, obj)
	}
}

func BSONValueDecoderFunc(t bsontype.Type, data []byte) DecoderFunc {
	vr := bsonrw.NewBSONValueReader(t, data)
	dec, err := bson.NewDecoder(vr)
	return func(obj any) error {
		if err != nil {
			return err
		}
		return dec.Decode(obj)
	}
}

type Unmarshaler interface {
	Unmarshal(unmarshal func(any) error) error
}

// unmarshaler is a wrapper for an Unmarshaler implementation that, itself, implements json.Unmarshaler,
// yaml.Unmarshaler, bson.Unmarshaler, ...
type unmarshaler struct {
	Unmarshaler
}

func newUnmarshaler(u Unmarshaler) *unmarshaler {
	return &unmarshaler{Unmarshaler: u}
}

func (u *unmarshaler) UnmarshalJSON(data []byte) error {
	return u.Unmarshal(JSONDecoderFunc(data).Decode)
}

func (u *unmarshaler) UnmarshalYAML(unmarshal func(any) error) error {
	return u.Unmarshal(unmarshal)
}

func (u *unmarshaler) UnmarshalBSON(data []byte) error {
	return u.Unmarshal(BSONDecoderFunc(data))
}

func Unmarshal(decoder Decoder, u Unmarshaler) (err error) {
	err = decoder.Decode(newUnmarshaler(u))
	return
}

type Marshaller interface {
	Marshal() any
}

type encoder struct {
	Marshaller
}

func newEncoder(m Marshaller) *encoder {
	return &encoder{Marshaller: m}
}

func Marshal(m Marshaller) *encoder {
	return newEncoder(m)
}

func (e *encoder) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Marshal())
}

func (e *encoder) MarshalYAML() (any, error) {
	return e.Marshal(), nil
}

func (e *encoder) MarshalBSON() ([]byte, error) {
	return bson.Marshal(e.Marshal())
}
