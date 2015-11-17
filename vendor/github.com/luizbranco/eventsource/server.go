package eventsource

import (
	"time"

	"github.com/paulsmith/newrelic-go-agent/newrelic"
)

// A server manages all clients, adding and removing them from the pool and
// receiving incoming events to forward to clients
type server struct {
	add      chan client
	remove   chan client
	events   chan Event
	hearbeat time.Duration
	metrics  Metrics
}

// The listen method is used to receive messages to add, remove and send
// events to clients. Every X seconds it sends a ping message to all clients to
// detect stale connections
func (s server) listen() {
	var clients []client
	tick := time.Tick(s.hearbeat)
	for {
		select {
		case c := <-s.add:
			clients = s.spawn(clients, c)
		case c := <-s.remove:
			clients = s.kill(clients, c)
		case e := <-s.events:
			go func() {
				start := time.Now()
				durations := send(e, clients)
				s.metrics.EventDone(e, time.Since(start), durations)
			}()
		case <-tick:
			go func() {
				durations := send(ping{}, clients)
				s.metrics.ClientCount(len(durations))
			}()
		}
	}
}

// send receives an event and a list of clients and send to them the
// text/stream data to be written on the client's connection. It returns a list
// of time.Duration each client took. 0 duration means that the data wasn't
// sent.
func send(e Event, clients []client) []time.Duration {
	durations := []time.Duration{}
	clients = e.Clients(clients)
	size := len(clients)
	if size == 0 {
		return durations
	}
	txnID := newrelic.BeginTransaction()
	newrelic.SetTransactionName(txnID, "publish_event")
	done := make(chan time.Duration, size)
	p := payload{data: e.Bytes(), done: done}

	for _, c := range clients {
		go func(c client) {
			select {
			case c.events <- p:
			case <-c.done:
				p.done <- 0
			}
		}(c)
	}

	for i := 0; i < size; i++ {
		d := <-done
		durations = append(durations, d)
	}
	newrelic.EndTransaction(txnID)
	return durations
}

// The spawn adds a new client to the clients list and launches a goroutine for
// the client to listen to incoming messages. The client receives the remove
// channel necessary to unsubscribe itself from the server.
func (s server) spawn(clients []client, c client) []client {
	go c.listen(s.remove)
	clients = append(clients, c)
	return clients
}

// The kill removes a client from the client list by comparing their events
// channel. The client is removed by being moved to the end of the list and
// reducing the slice length.
func (s server) kill(clients []client, client client) []client {
	index := -1
	for i, c := range clients {
		if client.events == c.events {
			index = i
			break
		}
	}

	if index == -1 {
		panic("client not found")
	}

	last := len(clients) - 1
	if index < last {
		swap := clients[last]
		clients[index] = swap
	}
	clients = clients[:last]

	return clients
}
