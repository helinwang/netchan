# netchan

[![Build Status](https://travis-ci.org/helinwang/netchan.svg?branch=master)](https://travis-ci.org/helinwang/netchan)

Send and receive over the network with the built-in Go channel.

## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/helinwang/netchan
```

## Examples

Send and receive over TCP:

```Go
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
```

Send and receive over Unix domain socket:

```Go
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
```
