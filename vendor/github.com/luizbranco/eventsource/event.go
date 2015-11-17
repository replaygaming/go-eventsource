package eventsource

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"strconv"
)

// Event is an interface that defines what the event payload is and to which
// connections it must be sent
type Event interface {
	// Bytes returns the data to be written on the clients connection
	Bytes() []byte

	// Clients receives a list of clients and return a filtered list.
	Clients([]client) []client
}

// DefaultEvent implements the Event interface
type DefaultEvent struct {
	ID       int
	Name     string
	Message  []byte
	Channels []string
	Compress bool
}

// Bytes returns the text/stream message to be sent to the client.
// If the event has name, it is added first, then the data. Optionally, the data
// can be compressed using zlib.
func (e DefaultEvent) Bytes() []byte {
	var buf bytes.Buffer
	if e.ID > 0 {
		buf.WriteString("id: ")
		buf.WriteString(strconv.Itoa(e.ID))
		buf.WriteString("\n")
	}
	if e.Name != "" {
		buf.WriteString("event: ")
		buf.WriteString(e.Name)
		buf.WriteString("\n")
	}
	buf.WriteString("data: ")
	if e.Compress {
		buf.WriteString(e.deflate())
	} else {
		buf.Write(e.Message)
	}
	buf.WriteString("\n\n")
	return buf.Bytes()
}

// Clients selects clients that have at least one channel in
// common with the event or all clients if the event has no channel.
func (e DefaultEvent) Clients(clients []client) []client {
	if len(e.Channels) == 0 {
		return clients
	}
	var subscribed []client
	for _, client := range clients {
	channels:
		for _, cChans := range client.channels {
			for _, eChans := range e.Channels {
				if cChans == eChans {
					subscribed = append(subscribed, client)
					break channels
				}
			}
		}
	}
	return subscribed
}

// deflate compress the event message using zlib default compression and
// returns a base64 encoded string.
func (e DefaultEvent) deflate() string {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(e.Message)
	w.Close()
	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

type ping struct{}

func (ping) Bytes() []byte {
	return []byte(":ping\n\n")
}

func (ping) Clients(clients []client) []client {
	return clients
}
