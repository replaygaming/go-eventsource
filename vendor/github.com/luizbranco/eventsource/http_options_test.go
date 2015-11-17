package eventsource

import (
	"bytes"
	"net/http"
	"testing"
)

func TestDefaultHttpOptionsBytesWithoutRetry(t *testing.T) {
	expecting := []byte(
		`HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive

`)

	var req, _ = http.NewRequest("GET", "/", nil)
	options := DefaultHttpOptions{}
	result := options.Bytes(req)

	if !bytes.Equal(result, expecting) {
		t.Errorf("expected:\n%q\nto equal to:\n%q\n", expecting, result)
	}
}

func TestDefaultHttpOptionsBytesWithRetry(t *testing.T) {
	expecting := []byte(
		`HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive

retry: 2000

`)

	var req, _ = http.NewRequest("GET", "/", nil)
	options := DefaultHttpOptions{Retry: 2000}
	result := options.Bytes(req)

	if !bytes.Equal(result, expecting) {
		t.Errorf("expected:\n%q\nto equal to:\n%q\n", expecting, result)
	}
}

func TestDefaultHttpOptionsBytesWithCorsEnabled(t *testing.T) {
	expecting := []byte(
		`HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive
Access-Control-Allow-Credentials: true
Access-Control-Allow-Origin: http://localhost/

`)

	var req, _ = http.NewRequest("GET", "/", nil)
	req.Header.Set("origin", "http://localhost/")
	options := DefaultHttpOptions{Cors: true}
	result := options.Bytes(req)

	if !bytes.Equal(result, expecting) {
		t.Errorf("expected:\n%q\nto equal to:\n%q\n", expecting, result)
	}
}

func TestDefaultHttpOptionsBytesWithOldBrowserSupportEnabled(t *testing.T) {
	expecting := []byte("HTTP/1.1 200 OK\nContent-Type: text/event-stream\nCache-Control: no-cache\nConnection: keep-alive\n\n:                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                \n")

	var req, _ = http.NewRequest("GET", "/", nil)
	options := DefaultHttpOptions{OldBrowserSupport: true}
	result := options.Bytes(req)

	if !bytes.Equal(result, expecting) {
		t.Errorf("expected:\n%q\nto equal to:\n%q\n", expecting, result)
	}
}
