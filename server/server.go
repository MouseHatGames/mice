package server

type Server interface {
	AddHandler(h interface{})
	Handle(path string, data []byte) error
}

type server struct {
	handlers map[string]*handler
}

func NewServer() Server {
	return &server{
		handlers: make(map[string]*handler),
	}
}

func (s *server) AddHandler(h interface{}) {
	hdl := newHandler(h)
	s.handlers[hdl.Name] = hdl
}

func (s *server) Handle(path string, data []byte) error {

}
