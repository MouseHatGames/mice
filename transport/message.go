package transport

type Message struct {
	Headers map[string]string
	Data    []byte
}
