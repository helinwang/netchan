package netchan

type SendRecver interface {
	Send(path string, body []byte) error
	Recv(path string) []byte
}

type SendRecv struct {
}

func NewSendRecv() *SendRecv {
	return &SendRecv{}
}

func (s *SendRecv) ListenAndServe(addr string) {
	// TODO
}

func (s *SendRecv) Recv(path string) []byte {
	// TODO
	return nil
}

func (s *SendRecv) Send(path string, body []byte) error {
	// TODO
	return nil
}
