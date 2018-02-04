package netchan

import (
	"bytes"
	"encoding/gob"
	"reflect"
)

type Handler struct {
	sr SendRecver
}

func NewHandler(sr SendRecver) *Handler {
	return &Handler{sr: sr}
}

func (s *Handler) HandleRecv(topic string, ch chan<- interface{}, t reflect.Type) error {
	for {
		body := s.sr.Recv(topic)
		v := reflect.New(t)
		dec := gob.NewDecoder(bytes.NewReader(body))
		err := dec.DecodeValue(v)
		if err != nil {
			return err
		}

		ch <- v.Elem().Interface()
	}
}

func (s *Handler) HandleSend(addr, topic string, ch <-chan interface{}) error {
	for v := range ch {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(v)
		if err != nil {
			return err
		}

		err = s.sr.Send(addr, topic, buf.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}
