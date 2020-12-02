package discovery

import "errors"

type Discovery interface {
	Find(svc string) (host string, err error)
}

var ErrServiceNotRegistered = errors.New("service not registered")
