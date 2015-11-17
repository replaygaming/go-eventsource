package eventsource

import (
	"bufio"
	"bytes"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

type hijackerFail struct {
	httptest.ResponseRecorder
}

func newHijackerFail() *hijackerFail {
	return &hijackerFail{httptest.ResponseRecorder{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}}
}

func (w *hijackerFail) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, errors.New("not supported")
}

type connClosed struct {
	httptest.ResponseRecorder
	conn net.Conn
}

type dummyAddr string

func (a dummyAddr) Network() string { return string(a) }
func (a dummyAddr) String() string  { return string(a) }

type noopConn struct{}

func (noopConn) Read(b []byte) (int, error)         { return 0, errors.New("closed") }
func (noopConn) Write(b []byte) (int, error)        { return 0, errors.New("closed") }
func (noopConn) Close() error                       { return nil }
func (noopConn) LocalAddr() net.Addr                { return dummyAddr("local-addr") }
func (noopConn) RemoteAddr() net.Addr               { return dummyAddr("remote-addr") }
func (noopConn) SetDeadline(t time.Time) error      { return nil }
func (noopConn) SetReadDeadline(t time.Time) error  { return nil }
func (noopConn) SetWriteDeadline(t time.Time) error { return nil }

func newConnClosed() *connClosed {
	return &connClosed{httptest.ResponseRecorder{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}, noopConn{}}
}

func (w *connClosed) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.conn, nil, nil
}

func stubTCPClient() client {
	net.Listen("tcp4", "127.0.0.1:4000")
	conn, _ := net.Dial("tcp4", "127.0.0.1:4000")
	c := client{events: make(chan payload), conn: conn, done: make(chan bool)}
	return c
}
