package eventsource

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestAddChannel(t *testing.T) {
	s := server{add: make(chan client)}
	c := client{}
	go s.listen()
	select {
	case s.add <- c:
	case <-time.Tick(1 * time.Second):
		t.Errorf("expected server to be listening to add channel")
	}
}

func TestRemoveChannel(t *testing.T) {
	s := server{add: make(chan client), remove: make(chan client)}
	c := client{}
	go s.listen()
	s.add <- c
	select {
	case s.remove <- c:
	case <-time.Tick(1 * time.Second):
		t.Errorf("expected server to be listening to remove channel")
	}
}

func TestServerPing(t *testing.T) {
	s := server{hearbeat: 1 * time.Nanosecond, add: make(chan client), metrics: NoopMetrics{}}
	e := ping{}
	c := client{events: make(chan payload), conn: noopConn{}}
	go s.listen()
	s.add <- c
	p := <-c.events

	expecting := e.Bytes()
	result := p.data

	if !bytes.Equal(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestServerSpawn(t *testing.T) {
	s := server{}
	c := client{}
	expecting := []client{c}
	result := s.spawn([]client{}, c)
	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("expected:\n%v\nto be equal to:\n%v\n", expecting, result)
	}
}

func TestServerKill(t *testing.T) {
	s := server{}
	c1 := client{events: make(chan payload)}
	c2 := client{events: make(chan payload)}
	expecting := []client{c2}
	result := s.kill([]client{c1, c2}, c1)
	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("expected:\n%v\nto be equal to:\n%v\n", expecting, result)
	}
}

func TestServerKillPanic(t *testing.T) {
	c1 := client{events: make(chan payload)}
	c2 := client{events: make(chan payload)}
	s := server{}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected function to panic, it did not\n")
		}
	}()
	s.kill([]client{c2}, c1)
}

func TestSendPayload(t *testing.T) {
	e := DefaultEvent{Message: message}
	c := client{events: make(chan payload)}
	go c.listen(make(chan client))
	go send(e, []client{c})
	p := <-c.events

	expecting := e.Bytes()
	result := p.data

	if !bytes.Equal(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}

func TestSendMetrics(t *testing.T) {
	e := DefaultEvent{Message: message}
	c := stubTCPClient()
	go c.listen(make(chan client))
	result := send(e, []client{c})
	if len(result) != 1 {
		t.Errorf("expected:\n1 duration\ngot:\n%v\n", len(result))
	}
}

func TestSendError(t *testing.T) {
	e := DefaultEvent{Message: message}
	c := stubTCPClient()
	close(c.done)
	result := send(e, []client{c})

	expecting := []time.Duration{0}
	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
}
