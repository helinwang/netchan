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

func (s *Handler) HandleRecv(path string, ch chan<- interface{}, t reflect.Type) error {
	for {
		body := s.sr.Recv(path)
		v := reflect.New(t)
		dec := gob.NewDecoder(bytes.NewReader(body))
		err := dec.DecodeValue(v)
		if err != nil {
			return err
		}

		ch <- v.Elem().Interface()
	}
}

func (s *Handler) HandleSend(path string, ch <-chan interface{}) error {
	for v := range ch {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(v)
		if err != nil {
			return err
		}

		err = s.sr.Send(path, buf.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}
