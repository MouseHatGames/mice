package errors

import "testing"

func TestEncodeDecode(t *testing.T) {
	err := NewError(123, "hello %d", 10)

	str, _ := err.Encode()
	err2, e := Decode(str)
	if e != nil {
		t.Fatalf("failed to decode: %s", e)
	}

	if err2.StatusCode != err.StatusCode {
		t.Fatalf("status codes differ: %d and %d", err2.StatusCode, err.StatusCode)
	}

	if err2.Detail != err.Detail {
		t.Fatalf("details differ: %s and %s", err2.Detail, err.Detail)
	}

	if err2.ID != err.ID {
		t.Fatalf("ids differ: %s and %s", err2.ID, err.ID)
	}
}
