package transport

type Message struct {
	Headers map[string]string
	Data    []byte
}

func NewMessage() *Message {
	return &Message{
		Headers: make(map[string]string),
	}
}

func (m *Message) SetError(err error) {
	m.Headers[HeaderError] = err.Error()
}
