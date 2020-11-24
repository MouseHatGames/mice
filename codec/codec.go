package codec

type Codec interface {
	Marshal(msg interface{}) ([]byte, error)
	Unmarshal(b []byte, out interface{}) error
}
