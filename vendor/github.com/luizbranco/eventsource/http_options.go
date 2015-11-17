package eventsource

import (
	"bytes"
	"fmt"
	"net/http"
)

const (
	// HTTP HEADER sent to browser to upgrade the protocol to event-stream
	HEADER = `HTTP/1.1 200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive`
)

var padding string

func init() {
	var buf bytes.Buffer
	buf.WriteByte(':')
	for i := 0; i < 2048; i++ {
		buf.WriteByte(' ')
	}
	buf.WriteByte('\n')
	padding = buf.String()
}

type HttpOptions interface {
	Bytes(*http.Request) []byte
}

type DefaultHttpOptions struct {
	// Cors sets if the server has Cross-Origin Resource Sharing enabled
	Cors bool

	// Internet Explorer < 10 and Chrome < 13 needs a message padding to successfully establish a
	// text stream connection See
	//http://blogs.msdn.com/b/ieinternals/archive/2010/04/06/comet-streaming-in-internet-explorer-with-xmlhttprequest-and-xdomainrequest.aspx
	OldBrowserSupport bool

	// Retry is the amout of time in milliseconds that the client must retry a
	// reconnection
	Retry int
}

// The Bytes function writes a header and body to the browser to establish a
// text/stream connection with retry option and CORS if enabled.
func (h DefaultHttpOptions) Bytes(req *http.Request) []byte {
	var buf bytes.Buffer
	buf.WriteString(HEADER)
	if origin := req.Header.Get("origin"); h.Cors && origin != "" {
		buf.WriteString("\nAccess-Control-Allow-Credentials: true\n")
		cors := fmt.Sprintf("Access-Control-Allow-Origin: %s", origin)
		buf.WriteString(cors)
	}
	buf.WriteString("\n\n")
	if h.OldBrowserSupport {
		buf.WriteString(padding)
	}
	if h.Retry > 0 {
		retry := fmt.Sprintf("retry: %d\n\n", h.Retry)
		buf.WriteString(retry)
	}
	return buf.Bytes()
}
