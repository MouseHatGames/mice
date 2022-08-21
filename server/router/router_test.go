package router

import (
	"reflect"
	"testing"

	"github.com/MouseHatGames/mice/options"
	"github.com/stretchr/testify/assert"
)

type mockCodec struct {
	out interface{}
}

func (*mockCodec) Marshal(msg interface{}) ([]byte, error) {
	return nil, nil
}
func (c *mockCodec) Unmarshal(b []byte, out interface{}) error {
	c.out = out
	return nil
}

func TestDecode(t *testing.T) {
	c := &mockCodec{}
	s := &router{
		opts: &options.Options{Codec: c},
	}

	ret, err := s.decode(reflect.TypeOf(&dummy{}), []byte{})

	assert.Nil(t, err)
	assert.NotNil(t, ret)
	assert.IsType(t, &dummy{}, ret)
}
