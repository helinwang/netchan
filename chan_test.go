package netchan_test

import (
	"reflect"
	"sync"
	"testing"

	"github.com/helinwang/netchan"
	"github.com/stretchr/testify/require"
)

type mySendRecver struct {
	mu sync.Mutex
	m  map[string]chan []byte
}

func (m *mySendRecver) Send(addr, topic string, body []byte) error {
	m.mu.Lock()
	ch := m.m[topic]
	if ch == nil {
		ch = make(chan []byte, 1000)
		m.m[topic] = ch
	}
	m.mu.Unlock()

	ch <- body
	return nil
}

func (m *mySendRecver) Recv(topic string) []byte {
	m.mu.Lock()
	ch := m.m[topic]
	if ch == nil {
		ch = make(chan []byte, 1000)
		m.m[topic] = ch
	}
	m.mu.Unlock()

	return <-ch
}

func TestSendRecv(t *testing.T) {
	type data struct {
		A int
		B float32
	}

	send := make(chan interface{})
	recv := make(chan interface{})

	s := netchan.NewHandler(&mySendRecver{m: make(map[string]chan []byte)})
	go func() {
		err := s.HandleSend("", "test", send)
		require.Nil(t, err)
	}()

	go func() {
		err := s.HandleRecv("test", recv, reflect.TypeOf(data{}))
		require.Nil(t, err)
	}()

	d := data{A: 1, B: 2}
	send <- d
	r := <-recv
	require.Equal(t, d, r.(data))
}
