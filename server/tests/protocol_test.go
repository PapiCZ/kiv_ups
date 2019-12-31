package main

import (
	"bufio"
	"bytes"
	"io"
	"kiv_ups_server/net/tcp/protocol"
	"sync"
	"testing"
)

type FooMessage struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
}

func (f FooMessage) GetTypeId() protocol.MessageType {
	return 125
}

func (f FooMessage) New() protocol.Message {
	return &FooMessage{}
}

func TestSimpleEncode(t *testing.T) {
	proto := protocol.GameProtocol{}
	reader, writer := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		r := bufio.NewReader(reader)
		s, _, _ := r.ReadLine()

		if !bytes.Equal(s, []byte(`|125|45|{"foo":"verylongtextmessage","bar":654789123}`)) {
			t.Errorf("Bad encoded message: %s", s)
		}

		wg.Done()
	}()

	proto.Encode(FooMessage{
		Foo: "verylongtextmessage",
		Bar: 654789123,
	}, writer)

	wg.Wait()
}

func TestSimpleDecode(t *testing.T) {
	def := protocol.NewDefinition()
	def.Register(FooMessage{})
	proto := protocol.GameProtocol{Def: def}
	reader, writer := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		writer.Write([]byte(`|125|45|{"foo":"verylongtextmessage","bar":654789123}`))
		wg.Done()
	}()

	message, _ := proto.Decode(reader)
	msg := message.(*FooMessage)

	wg.Wait()

	if msg.Foo != "verylongtextmessage" || msg.Bar != 654789123 {
		t.Error("Bad decoded message:", msg)
	}
}

func TestEncodeAndDecode(t *testing.T) {
	def := protocol.NewDefinition()
	def.Register(FooMessage{})
	proto := protocol.GameProtocol{Def: def}
	reader, writer := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		r := bufio.NewReader(reader)
		s, _, _ := r.ReadLine()
		reader2, writer2 := io.Pipe()

		wg2 := sync.WaitGroup{}
		wg2.Add(1)
		go func() {
			writer2.Write(s)
			wg2.Done()
		}()

		msg, _ := proto.Decode(reader2)
		message := msg.(*FooMessage)

		wg2.Wait()

		if message.Foo != "verylongtextmessage" || message.Bar != 654789123 {
			t.Error("Bad decoded message:", message)
		}

		wg.Done()
	}()

	proto.Encode(FooMessage{
		Foo: "verylongtextmessage",
		Bar: 654789123,
	}, writer)

	wg.Wait()
}

func TestDecodeAndEncode(t *testing.T) {
	def := protocol.NewDefinition()
	def.Register(FooMessage{})
	proto := protocol.GameProtocol{Def: def}
	reader, writer := io.Pipe()

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		writer.Write([]byte(`|125|45|{"foo":"verylongtextmessage","bar":654789123}`))
		wg.Done()
	}()

	msg, _ := proto.Decode(reader)
	message := msg.(*FooMessage)
	reader2, writer2 := io.Pipe()

	wg.Add(1)
	go func() {
		r := bufio.NewReader(reader2)
		s, _, _ := r.ReadLine()

		if !bytes.Equal(s, []byte(`|125|45|{"foo":"verylongtextmessage","bar":654789123}`)) {
			t.Errorf("Bad encoded message: %s", s)
		}

		wg.Done()
	}()

	proto.Encode(message, writer2)

	wg.Wait()
}
