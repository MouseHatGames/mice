package discovery

import "net"

type Discovery interface {
	Find(svc string) (net.IP, error)
}
