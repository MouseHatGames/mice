package transport

import (
	goerrors "errors"

	"github.com/MouseHatGames/mice/errors"
	"github.com/google/uuid"
)

const (
	HeaderPath      = "path"
	HeaderError     = "error"
	HeaderRequestID = "reqid"
)

type MessageHeaders map[string]string

func (h MessageHeaders) GetRequestID() (id uuid.UUID, hasID bool) {
	uuidStr, ok := h[HeaderRequestID]
	if !ok {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}

func (h MessageHeaders) MustGetRequestID() uuid.UUID {
	id, _ := h.GetRequestID()
	return id
}

func (h MessageHeaders) SetRequestID(id uuid.UUID) {
	h[HeaderRequestID] = id.String()
}

func (h MessageHeaders) SetRandomRequestID() {
	h[HeaderRequestID] = uuid.NewString()
}

func (h MessageHeaders) GetPath() (path string, hasPath bool) {
	path, hasPath = h[HeaderPath]
	return
}

func (h MessageHeaders) SetPath(path string) {
	h[HeaderPath] = path
}

func (h MessageHeaders) GetError() (err error, hasError bool) {
	value, ok := h[HeaderError]
	if !ok {
		return nil, false
	}

	if merr, ok := errors.Decode(value); ok {
		return merr, true
	}

	return goerrors.New(value), true
}

func (h MessageHeaders) SetError(err error) {
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

	h[HeaderError] = value
}
