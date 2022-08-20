package json

import (
	"encoding/json"

	"github.com/MouseHatGames/mice/options"
)

type jsonCodec struct{}

func Codec() options.Option {
	return func(o *options.Options) {
		o.Codec = &jsonCodec{}
	}
}

func (*jsonCodec) Marshal(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}

func (*jsonCodec) Unmarshal(b []byte, out interface{}) error {
	return json.Unmarshal(b, out)
}
