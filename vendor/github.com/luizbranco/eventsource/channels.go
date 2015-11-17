package eventsource

import (
	"net/http"
	"strings"
)

// ChannelSubscriber interface is used to determine which channels a client has
// subscribed to. This package has two built-in implementations: NoChannels and
// QueryStringChannels, but you can implement your own.
type ChannelSubscriber interface {
	ParseRequest(*http.Request) []string
}

// NoChannels implements the ChannelSubscriber interface by always returning an
// empty list of channels. This is useful for implementing an eventsource with
// global messages only.
type NoChannels struct{}

// ParseRequest returns an empty list of channels.
func (n NoChannels) ParseRequest(req *http.Request) []string {
	return []string{}
}

// QueryStringChannels implements the ChannelSubscriber interface by parsing
// the request querystring and extracting channels separated by commas. Eg.:
// /?channels=a,b,c
type QueryStringChannels struct {
	Name string
}

// ParseRequest parses the querystring and extracts the Name params, spliting
// it by commas.
func (n QueryStringChannels) ParseRequest(req *http.Request) []string {
	channels := req.URL.Query().Get(n.Name)
	if channels == "" {
		return []string{}
	}
	return strings.Split(channels, ",")
}
