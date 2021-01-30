package errors

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net/http"
)

type Error struct {
	StatusCode int16
	ID, Detail string
}

func NewError(code int16, format string, a ...interface{}) *Error {
	return &Error{
		StatusCode: code,
		Detail:     fmt.Sprintf(format, a...),
	}
}

func NewErrorID(id string, code int16, format string, a ...interface{}) *Error {
	return &Error{
		StatusCode: code,
		ID:         id,
		Detail:     fmt.Sprintf(format, a...),
	}
}

// Encode transforms the error instance into a base64-encoded blob
func (e *Error) Encode() (string, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	if err := enc.Encode(e); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// Decode tries to decode en error from a base64-encoded string
func Decode(str string) (*Error, error) {
	bdec, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	b := bytes.NewReader(bdec)
	dec := gob.NewDecoder(b)

	var merr Error
	if err := dec.Decode(&merr); err != nil {
		return nil, err
	}

	return &merr, nil
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", http.StatusText(int(e.StatusCode)), e.Detail)
}
