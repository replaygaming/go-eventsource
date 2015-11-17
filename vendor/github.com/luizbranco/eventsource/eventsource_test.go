package eventsource

import (
	"bufio"
	"bytes"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestEventsourceStartNoChannels(t *testing.T) {
	es := Eventsource{}
	es.Start()
	result, ok := es.ChannelSubscriber.(NoChannels)
	if !ok {
		t.Errorf("expected to be NoChannels\ngot:\n%T\n", result)
	}
}

func TestEventsourceStartDefaultHttpOptions(t *testing.T) {
	es := Eventsource{}
	es.Start()
	opts, ok := es.HttpOptions.(DefaultHttpOptions)
	if !ok {
		t.Errorf("expected to be DefaultHttpOptions\ngot:\n%T\n", opts)
	}
	expecting := 2000
	retry := opts.Retry
	if expecting != retry {
		t.Errorf("expected retry to be:\n%d\ngot:\n%d\n", expecting, retry)
	}
	cors := opts.Cors
	if !cors {
		t.Errorf("expected Cors to be:\ntrue\ngot:\n%t\n", cors)
	}
	old := opts.OldBrowserSupport
	if !old {
		t.Errorf("expected OldBrowserSupport to be:\ntrue\ngot:\n%t\n", old)
	}
}

func TestEventsourceStartMetricsJSONLogger(t *testing.T) {
	es := Eventsource{}
	es.Start()
	metrics, ok := es.Metrics.(DefaultMetrics)
	if !ok {
		t.Errorf("expected to be DefaultMetrics\ngot:\n%T\n", metrics)
	}
}

func TestEventsourceStartServerChannels(t *testing.T) {
	es := Eventsource{}
	es.Start()
	s := es.server

	if s.add == nil {
		t.Errorf("expected server add channel to be created")
	}
	if s.remove == nil {
		t.Errorf("expected server remove channel to be created")
	}
	if s.events == nil {
		t.Errorf("expected server events channel to be created")
	}
}

func TestEventsourceStartServerListen(t *testing.T) {
	es := Eventsource{}
	es.Metrics = NoopMetrics{}
	es.Start()
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(1 * time.Second)
		timeout <- true
	}()
	select {
	case es.server.events <- DefaultEvent{}:
	case <-timeout:
		t.Errorf("expected server to be listening to events")
	}
}

func TestEventsourceSend(t *testing.T) {
	es := Eventsource{}
	events := make(chan Event, 1)
	es.server = server{events: events}
	expecting := DefaultEvent{Name: "test"}
	es.Send(expecting)
	result := <-events
	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("expected:\n%v\nto be equal to:\n%v\n", expecting, result)
	}
}

func TestEventsourceServeHTTP(t *testing.T) {
	var wg sync.WaitGroup
	s := server{add: make(chan client, 1)}
	wg.Add(1)
	go func() {
		select {
		case c := <-s.add:
			if c.conn == nil {
				t.Errorf("expecting client connection to be assigned")
			}
			if c.events == nil {
				t.Errorf("expecting client events chan to be open")
			}
			expecting := []string{"a", "b", "c"}
			result := c.channels
			if !reflect.DeepEqual(expecting, result) {
				t.Errorf("expected:\n%s\nto be equal to:\n%s\n", expecting, result)
			}
			if c.done == nil {
				t.Errorf("expecting client done chan to be open")
			}
		case <-time.Tick(1 * time.Second):
			t.Errorf("expecting client to be added")
		}
		wg.Done()
	}()
	es := &Eventsource{}
	es.ChannelSubscriber = QueryStringChannels{Name: "channels"}
	es.Start()
	es.server = s
	server := httptest.NewServer(es)
	res, _ := http.Get(server.URL + "?channels=a,b,c")
	reader := bufio.NewReader(res.Body)
	line, _ := reader.ReadBytes('\n')
	if line[0] != ':' {
		t.Errorf("expected:\n%s\nto be equal to:\n%s\n", ":", line)
	}
	wg.Wait()
}

func TestEventsourceServeHTTPHijackingNotSupported(t *testing.T) {
	es := Eventsource{}
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	es.ServeHTTP(w, r)
	errCode := 500
	code := w.Code
	if errCode != code {
		t.Errorf("expected:\n%d\nto be equal to:\n%d\n", errCode, code)
	}

	errMsg := []byte(HijackingError + "\n")
	msg := w.Body.Bytes()

	if !bytes.Equal(errMsg, msg) {
		t.Errorf("expected:\n%s\nto be equal to:\n%s\n", errMsg, msg)
	}
}

func TestEventsourceServeHTTPHijackingError(t *testing.T) {
	es := Eventsource{}
	w := newHijackerFail()
	r, _ := http.NewRequest("GET", "/", nil)
	es.ServeHTTP(w, r)
	errCode := 500
	code := w.Code
	if errCode != code {
		t.Errorf("expected:\n%d\nto be equal to:\n%d\n", errCode, code)
	}

	errMsg := []byte("not supported" + "\n")
	msg := w.Body.Bytes()

	if !bytes.Equal(errMsg, msg) {
		t.Errorf("expected:\n%s\nto be equal to:\n%s\n", errMsg, msg)
	}
}

func TestEventsourceServeHTTPConnectionClosed(t *testing.T) {
	es := Eventsource{}
	es.Start()
	w := newConnClosed()
	r, _ := http.NewRequest("GET", "/", nil)
	es.ServeHTTP(w, r)
	_, err := w.conn.Read([]byte(""))
	if err == nil {
		t.Errorf("expected connection to be closed")
	}
}
