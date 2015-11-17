package eventsource

import (
	"bytes"
	"reflect"
	"testing"
)

var message = []byte("{id: 1}")
var deflated = "eJyqzkyxUjCsBQQAAP//CfUCUQ=="

func TestDefaultEventBytesWithID(t *testing.T) {
	expecting := []byte("id: 1\ndata: {id: 1}\n\n")
	e := DefaultEvent{
		ID:      1,
		Message: message,
	}
	result := e.Bytes()
	if !bytes.Equal(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestDefaultEventBytesWithName(t *testing.T) {
	expecting := []byte("event: test\ndata: {id: 1}\n\n")
	e := DefaultEvent{
		Name:    "test",
		Message: message,
	}
	result := e.Bytes()
	if !bytes.Equal(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestDefaultEventBytesWithoutName(t *testing.T) {
	expecting := []byte("data: {id: 1}\n\n")
	e := DefaultEvent{
		Message: message,
	}
	result := e.Bytes()
	if !bytes.Equal(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestDefaultEventBytesWithCompression(t *testing.T) {
	expecting := []byte("data: " + deflated + "\n\n")
	e := DefaultEvent{
		Message:  message,
		Compress: true,
	}
	result := e.Bytes()
	if !bytes.Equal(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestDefaultEventDeflate(t *testing.T) {
	expecting := deflated
	e := DefaultEvent{Message: message}
	result := e.deflate()
	if expecting != result {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestDefaultEventClientsWithNoChannel(t *testing.T) {
	client1 := client{channels: []string{"a", "b"}}
	client2 := client{channels: []string{"c", "d"}}
	e := DefaultEvent{}

	expecting := []client{client1, client2}
	result := e.Clients([]client{client1, client2})

	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("expected:\n%v\nto be equal to:\n%v\n", expecting, result)
	}
}

func TestDefaultEventClientsWithChannels(t *testing.T) {
	client1 := client{channels: []string{"a", "b"}}
	client2 := client{channels: []string{"c", "d"}}
	e := DefaultEvent{Channels: []string{"b", "e"}}

	expected := []client{client1}
	result := e.Clients([]client{client1, client2})

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("expected:\n%v\nto be equal to:\n%v\n", expected, result)
	}
}

func TestPingBytes(t *testing.T) {
	expecting := []byte(":ping\n\n")
	result := ping{}.Bytes()
	if !bytes.Equal(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestPingClients(t *testing.T) {
	clients := []client{client{}}
	expecting := clients
	result := ping{}.Clients(clients)
	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("expected:\n%v\ngot:\n%v\n", expecting, result)
	}
}
