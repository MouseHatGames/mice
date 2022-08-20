package transport

type Message struct {
	MessageHeaders
	Data []byte
}

func NewMessage() *Message {
	return &Message{
		MessageHeaders: make(MessageHeaders),
	}
}
