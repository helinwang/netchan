package netchan_test

import (
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/helinwang/netchan"
)

func BenchmarkUnixSendRecv(b *testing.B) {
	addr := "tmp-" + strconv.Itoa(b.N)
	const (
		name = "test"
	)

	type data struct {
		A int
		B float32
	}

	send := make(chan interface{})
	recv := make(chan interface{})

	sr := netchan.NewSendRecv()
	go func() {
		err := sr.ListenAndServe("unix", addr)
		if err != nil {
			panic(err)
		}
	}()

	// wait for the server to start
	time.Sleep(30 * time.Millisecond)

	s := netchan.NewHandler(sr)
	go func() {
		err := s.HandleSend("unix", addr, name, send)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := s.HandleRecv(name, recv, reflect.TypeOf(make([]byte, 0)))
		if err != nil {
			panic(err)
		}
	}()

	toSend := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		toSend[i] = make([]byte, 1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		send <- toSend[i]
		<-recv
	}

	err := os.Remove(addr)
	if err != nil {
		panic(err)
	}
}

func BenchmarkTCPSendRecv(b *testing.B) {
	addr := ":" + strconv.Itoa(b.N%16000+8000)
	const (
		name = "test"
	)

	type data struct {
		A int
		B float32
	}

	send := make(chan interface{})
	recv := make(chan interface{})

	sr := netchan.NewSendRecv()
	go func() {
		err := sr.ListenAndServe("tcp", addr)
		if err != nil {
			panic(err)
		}
	}()

	// wait for the server to start
	time.Sleep(30 * time.Millisecond)

	s := netchan.NewHandler(sr)
	go func() {
		err := s.HandleSend("tcp", addr, name, send)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := s.HandleRecv(name, recv, reflect.TypeOf(make([]byte, 0)))
		if err != nil {
			panic(err)
		}
	}()

	toSend := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		toSend[i] = make([]byte, 1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		send <- toSend[i]
		<-recv
	}

}

func BenchmarkChanSendRecv(b *testing.B) {
	const (
		addr = ""
		name = "test"
	)

	type data struct {
		A int
		B float32
	}

	send := make(chan interface{})
	recv := make(chan interface{})

	s := netchan.NewHandler(&mySendRecver{m: make(map[string]chan []byte)})
	go func() {
		err := s.HandleSend("", addr, name, send)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err := s.HandleRecv(name, recv, reflect.TypeOf(make([]byte, 0)))
		if err != nil {
			panic(err)
		}
	}()

	toSend := make([][]byte, b.N)
	for i := 0; i < b.N; i++ {
		toSend[i] = make([]byte, 1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		send <- toSend[i]
		<-recv
	}

}
