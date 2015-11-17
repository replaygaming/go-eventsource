package eventsource

import (
	"bytes"
	"io"
	"net"
	"reflect"
	"testing"
	"time"
)

func checkRead(t *testing.T, r io.Reader, data []byte, wantErr error) {
	buf := make([]byte, len(data)+10)
	n, err := r.Read(buf)
	if err != wantErr {
		t.Errorf("read: %v", err)
		return
	}
	if n != len(data) || !bytes.Equal(buf[0:n], data) {
		t.Errorf("expected:\n%s\ngot:\n%s\n", data, buf[0:n])
		return
	}
}

func stubPipeClient() (client, net.Conn, chan client) {
	remove := make(chan client)
	done := make(chan bool)
	events := make(chan payload)
	read, write := net.Pipe()
	c := client{done: done, conn: write, events: events}
	go c.listen(remove)
	return c, read, remove
}

func TestClientListen(t *testing.T) {
	c, read, _ := stubPipeClient()
	expecting := []byte("test")
	go func() {
		c.events <- payload{data: expecting, done: make(chan time.Duration)}
	}()
	checkRead(t, read, expecting, nil)
}

func TestClientListenEventsClosed(t *testing.T) {
	c, read, _ := stubPipeClient()
	close(c.events)
	checkRead(t, read, nil, io.EOF)
}

func TestClientListenConnError(t *testing.T) {
	c, _, remove := stubPipeClient()
	c.conn.Close()
	done := make(chan time.Duration)
	go func() {
		c.events <- payload{data: []byte("test"), done: done}
	}()
	expecting := time.Duration(0)
	result := <-done
	if expecting != result {
		t.Errorf("expected:\n%s\ngot:\n%s\n", expecting, result)
	}
	go func() {
		removed := <-remove
		if !reflect.DeepEqual(c, removed) {
			t.Errorf("expected:\n%v\ngot:\n%v\n", c, removed)
		}
	}()
	_, ok := <-c.done
	if ok {
		t.Errorf("expected channel to be closed")
	}
}
