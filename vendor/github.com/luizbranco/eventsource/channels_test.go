package eventsource

import (
	"net/http"
	"reflect"
	"testing"
)

var req, _ = http.NewRequest("GET", "/?channels=a,b,c", nil)

func TestNoChannelsParseRequest(t *testing.T) {
	sub := NoChannels{}
	result := sub.ParseRequest(req)

	if len(result) > 0 {
		t.Errorf("expected:\n%q\nto be empty\n", result)
	}
}

func TestQueryStringChannelsParseRequestEmpty(t *testing.T) {
	sub := QueryStringChannels{}
	result := sub.ParseRequest(req)

	if len(result) > 0 {
		t.Errorf("expected:\n%q\nto be empty\n", result)
	}
}

func TestQueryStringChannelsParseRequest(t *testing.T) {
	sub := QueryStringChannels{Name: "channels"}
	result := sub.ParseRequest(req)
	expecting := []string{"a", "b", "c"}

	if !reflect.DeepEqual(expecting, result) {
		t.Errorf("expected:\n%q\nto be equal to:\n%q\n", result, expecting)
	}
}
