# netchan

[![Build Status](https://travis-ci.org/helinwang/netchan.svg?branch=master)](https://travis-ci.org/helinwang/netchan)

Send and receive over the network with the built-in Go channel.

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/helinwang/netchan
```

## Examples

Send and receive over TCP and UNIX:

```Go
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
```
