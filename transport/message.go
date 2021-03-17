package transport

import (
	goerrors "errors"

	"github.com/MouseHatGames/mice/errors"
)

type Message struct {
	Headers map[string]string
	Data    []byte
}

func NewMessage() *Message {
	return &Message{
		Headers: make(map[string]string),
	}
}

func (m *Message) SetError(err error) {
	var value string

	if merr, ok := err.(*errors.Error); ok {
		enc, err := merr.Encode()

		if err != nil {
			value = err.Error()
		} else {
			value = enc
		}
	} else {
		value = err.Error()
	}

	m.Headers[HeaderError] = value
}

func (m *Message) GetError() (err error, hasError bool) {
	value, ok := m.Headers[HeaderError]
	if !ok {
		return nil, false
	}

	if merr, ok := errors.Decode(value); ok {
		return merr, true
	}

	return goerrors.New(value), true
}
