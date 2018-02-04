package netchan

type SendRecver interface {
	Send(addr, topic string, body []byte) error
	Recv(topic string) []byte
}

type SendRecv struct {
}

func NewSendRecv() *SendRecv {
	return &SendRecv{}
}

func (s *SendRecv) ListenAndServe(addr string) {
	// TODO
}

func (s *SendRecv) Recv(topic string) []byte {
	// TODO
	return nil
}

func (s *SendRecv) Send(addr, topic string, body []byte) error {
	// TODO
	return nil
}
