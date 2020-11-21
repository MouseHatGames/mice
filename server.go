package mice

// Server represents an interface to the server-side gRPC
type Server interface {
}

type server struct {
}

var _ Server = (*server)(nil)
