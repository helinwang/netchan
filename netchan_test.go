package netchan_test

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/helinwang/netchan"
	"github.com/stretchr/testify/require"
)

type mySendRecver struct {
	mu sync.Mutex
	m  map[string]chan []byte
}

func (m *mySendRecver) Send(addr, name string, body []byte) error {
	m.mu.Lock()
	ch := m.m[name]
	if ch == nil {
		ch = make(chan []byte, 1000)
		m.m[name] = ch
	}
	m.mu.Unlock()

	ch <- body
	return nil
}

func (m *mySendRecver) Recv(name string) []byte {
	m.mu.Lock()
	ch := m.m[name]
	if ch == nil {
		ch = make(chan []byte, 1000)
		m.m[name] = ch
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

func ExampleTCP() {
	const (
		addr = ":8003"
		name = "test"
	)

	type data struct {
		A int
		B float32
	}

	send := make(chan interface{})
	recv := make(chan interface{})

	sr := netchan.NewSendRecv("tcp")
	go func() {
		err := sr.ListenAndServe(addr)
		if err != nil {
			panic(err)
		}
	}()

	// wait for the server to start
	time.Sleep(30 * time.Millisecond)

	s := netchan.NewHandler(sr)
	go func() {
		err := s.HandleSend(addr, name, send)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := s.HandleRecv(name, recv, reflect.TypeOf(data{}))
		if err != nil {
			panic(err)
		}
	}()

	d := data{A: 1, B: 2}
	send <- d
	r := <-recv
	fmt.Println(r)
	// Output:
	// {1 2}
}

func ExampleUnix() {
	const (
		addr = "tmp"
		name = "test"
	)

	type data struct {
		A int
		B float32
	}

	send := make(chan interface{})
	recv := make(chan interface{})

	sr := netchan.NewSendRecv("unix")
	go func() {
		err := sr.ListenAndServe(addr)
		if err != nil {
			panic(err)
		}
	}()

	// wait for the server to start
	time.Sleep(30 * time.Millisecond)

	s := netchan.NewHandler(sr)
	go func() {
		err := s.HandleSend(addr, name, send)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := s.HandleRecv(name, recv, reflect.TypeOf(data{}))
		if err != nil {
			panic(err)
		}
	}()

	d := data{A: 1, B: 2}
	send <- d
	r := <-recv
	fmt.Println(r)
	// Output:
	// {1 2}

	err := os.Remove(addr)
	if err != nil {
		panic(err)
	}
}
