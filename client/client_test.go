package client

import (
	"testing"

	"github.com/MouseHatGames/mice/broker"
	"github.com/stretchr/testify/assert"
)

type dummy struct {
	n int
}

type mockcodec struct {
	n int
}

func (*mockcodec) Marshal(msg interface{}) ([]byte, error) {
	return nil, nil
}

func (c *mockcodec) Unmarshal(b []byte, out interface{}) error {
	d := out.(*dummy)
	d.n = c.n
	return nil
}

func TestCreateCallback(t *testing.T) {
	c := &client{}

	t.Run("not function", func(t *testing.T) {
		_, err := c.createCallback(123)

		assert.Equal(t, ErrMustBeFunc, err)
	})
	t.Run("input number", func(t *testing.T) {
		_, err := c.createCallback(func(a, b int) {})

		assert.Equal(t, ErrInvalidInput, err)
	})
	t.Run("input not pointer", func(t *testing.T) {
		_, err := c.createCallback(func(a dummy) {})

		assert.Equal(t, ErrInputPointer, err)
	})
	t.Run("normal", func(t *testing.T) {
		cod := &mockcodec{n: 123}
		c := &client{
			codec: cod,
		}

		called := false

		fn, err := c.createCallback(func(d *dummy) {
			assert.Equal(t, d.n, cod.n)
			called = true
		})

		assert.Nil(t, err)
		fn(&broker.Message{})

		assert.True(t, called)
	})
}
