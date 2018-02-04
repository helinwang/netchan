package netchan

import (
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type RecvServer struct {
	mu sync.Mutex
	m  map[string]chan []byte
}

type Request struct {
	Name string
	Data []byte
}

func (s *RecvServer) getOrCreateCh(name string) chan []byte {
	s.mu.Lock()
	ch := s.m[name]
	if ch == nil {
		ch = make(chan []byte, 100)
		s.m[name] = ch
	}
	s.mu.Unlock()

	return ch
}

func (s *RecvServer) Put(req Request, _ *int) error {
	ch := s.getOrCreateCh(req.Name)
	ch <- req.Data
	return nil
}

type SendRecver interface {
	Send(address, name string, body []byte) error
	Recv(name string) []byte
}

type SendRecv struct {
	network string
	server  *RecvServer

	mu      sync.Mutex
	clients map[string]*rpc.Client
}

func NewSendRecv(network string) *SendRecv {
	return &SendRecv{
		network: network,
		server:  &RecvServer{m: make(map[string]chan []byte)},
		clients: make(map[string]*rpc.Client),
	}
}

func (s *SendRecv) ListenAndServe(address string) error {
	rpc.Register(s.server)
	rpc.HandleHTTP()
	l, err := net.Listen(s.network, address)
	if err != nil {
		return err
	}

	return http.Serve(l, nil)
}

func (s *SendRecv) Recv(name string) []byte {
	ch := s.server.getOrCreateCh(name)
	return <-ch
}

func (s *SendRecv) Send(address, name string, body []byte) error {
	s.mu.Lock()
	client := s.clients[address]
	s.mu.Unlock()
	if client == nil {
		var err error
		client, err = rpc.DialHTTP(s.network, address)
		if err != nil {
			return err
		}
	}

	var oldClient *rpc.Client
	s.mu.Lock()
	if s.clients[address] == nil {
		s.clients[address] = client
	} else {
		oldClient = client
		client = s.clients[address]
	}
	s.mu.Unlock()

	if oldClient != nil {
		err := oldClient.Close()
		if err != nil {
			return err
		}
	}

	req := Request{Name: name, Data: body}
	err := client.Call("RecvServer.Put", req, nil)
	if err != nil {
		return err
	}

	return nil
}
