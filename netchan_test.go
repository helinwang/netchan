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

func (m *mySendRecver) Send(network, addr, name string, body []byte) error {
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
		err := s.HandleSend("", "", "test", send)
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

func TestSendRecver(t *testing.T) {
	const (
		tcpAddr  = ":8004"
		unixAddr = "tmp"
	)

	type data struct {
		A int
		B float32
	}

	cases := []struct {
		network, addr, name string
	}{
		{"unix", unixAddr, "unixCh"},
		{"tcp", tcpAddr, "tcpCh"},
	}

	sr := netchan.NewSendRecv()
	go func() {
		err := sr.ListenAndServe("tcp", tcpAddr)
		require.Nil(t, err)
	}()

	go func() {
		err := sr.ListenAndServe("unix", unixAddr)
		require.Nil(t, err)
	}()

	// wait for the server to start
	time.Sleep(30 * time.Millisecond)

	for _, c := range cases {
		c := c
		send := make(chan interface{})
		recv := make(chan interface{})

		s := netchan.NewHandler(sr)
		go func() {
			err := s.HandleSend(c.network, c.addr, c.name, send)
			require.Nil(t, err)
		}()

		go func() {
			err := s.HandleRecv(c.name, recv, reflect.TypeOf(data{}))
			require.Nil(t, err)
		}()

		for i := 0; i < 100; i++ {
			d := data{A: i, B: float32(i)}
			send <- d
			r := <-recv
			require.Equal(t, d, r.(data))
		}
	}

	err := os.Remove(unixAddr)
	require.Nil(t, err)
}

func Example() {
	const (
		unixAddr = "tmp"
		tcpAddr  = ":8003"
		name     = "test"
	)

	type data struct {
		A int
		B float32
	}

	tcpSend := make(chan interface{})
	unixSend := make(chan interface{})
	recv := make(chan interface{})

	// this example shows sending and receiving over tcp and unix,
	// but you can use only tcp or unix, just send or just
	// receive.
	sr := netchan.NewSendRecv()
	go func() {
		err := sr.ListenAndServe("tcp", tcpAddr)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := sr.ListenAndServe("unix", unixAddr)
		if err != nil {
			panic(err)
		}
	}()

	// wait for the server to start
	time.Sleep(30 * time.Millisecond)

	s := netchan.NewHandler(sr)
	go func() {
		err := s.HandleSend("tcp", tcpAddr, name, tcpSend)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := s.HandleSend("unix", unixAddr, name, unixSend)
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
	tcpSend <- d
	r := (<-recv).(data)
	fmt.Println(r)

	d = data{A: 3, B: 4}
	unixSend <- d
	r = (<-recv).(data)
	fmt.Println(r)
	// Output:
	// {1 2}
	// {3 4}

	err := os.Remove(unixAddr)
	if err != nil {
		panic(err)
	}
}
