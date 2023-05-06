package transport

import (
	goerrors "errors"
	"strconv"

	"github.com/MouseHatGames/mice/errors"
	"github.com/google/uuid"
)

const (
	HeaderPath            = "path"
	HeaderError           = "error"
	HeaderRequestID       = "reqid"
	HeaderParentRequestID = "parentreq"
	HeaderUserID          = "userid"
)

type MessageHeaders map[string]string

func (h *MessageHeaders) ensure() MessageHeaders {
	if *h == nil {
		*h = make(MessageHeaders)
	}

	return *h
}

func (h MessageHeaders) GetPath() (path string, hasPath bool) {
	path, hasPath = h[HeaderPath]
	return
}

func (h *MessageHeaders) SetPath(path string) {
	h.ensure()[HeaderPath] = path
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

func (h *MessageHeaders) SetError(err error) {
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

	h.ensure()[HeaderError] = value
}

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

func (h *MessageHeaders) SetRequestID(id uuid.UUID) {
	h.ensure()[HeaderRequestID] = id.String()
}

func (h *MessageHeaders) SetRandomRequestID() {
	h.ensure()[HeaderRequestID] = uuid.NewString()
}

func (h MessageHeaders) GetParentRequestID() (id uuid.UUID, hasID bool) {
	uuidStr, ok := h[HeaderParentRequestID]
	if !ok {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}

func (h *MessageHeaders) SetParentRequestID(id uuid.UUID) {
	if id != uuid.Nil {
		h.ensure()[HeaderParentRequestID] = id.String()
	}
}

func (h *MessageHeaders) SetUserID(id uint32) {
	h.ensure()[HeaderUserID] = strconv.FormatUint(uint64(id), 10)
}

func (h MessageHeaders) GetUserID() (id uint32, hasID bool) {
	idStr, ok := h[HeaderUserID]
	if !ok {
		return 0, false
	}

	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, false
	}

	return uint32(id64), true
}
