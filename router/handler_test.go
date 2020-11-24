package router

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummy struct{}
type dummyin struct{}
type dummyout struct{}

func (*dummy) Test(ctx context.Context, data *dummyin) (*dummyout, error) {
	return nil, nil
}

func TestGetEndpoint(t *testing.T) {
	d := &dummy{}
	m := reflect.ValueOf(d).Type().Method(0)

	ep := getEndpoint(m)

	assert.NotNil(t, ep)
	assert.Equal(t, "Test", ep.Name)
	assert.Equal(t, reflect.TypeOf(&dummyin{}), ep.In)
	assert.Equal(t, reflect.TypeOf(&dummyout{}), ep.Out)
}
