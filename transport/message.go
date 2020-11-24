package transport

type Message struct {
	Headers map[string]string
	Data    []byte
}

func (m *Message) SetError(err error) {
	m.Headers[HeaderError] = err.Error()
}
