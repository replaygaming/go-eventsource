package eventsource

import (
	"log"
	"time"
)

// Metrics interface allows basic instrumentation of the server
type Metrics interface {
	ClientCount(int)
	EventDone(Event, time.Duration, []time.Duration)
}

// NoopMetrics implements the Metrics interface and does nothing. Useful
// for disable metrics.
type NoopMetrics struct{}

// ClientCount does nothing.
func (NoopMetrics) ClientCount(int) {}

// EventDone does nothing.
func (NoopMetrics) EventDone(Event, time.Duration, []time.Duration) {}

// DefaultMetrics implements the Metrics interface and logs events to the
// stdout.
type DefaultMetrics struct{}

// ClientCount does nothing.
func (DefaultMetrics) ClientCount(int) {}

// EventDone logs to stdout the avg time an event to be sent to clients.
// Clients with error are ignored.
func (DefaultMetrics) EventDone(e Event, _ time.Duration, durations []time.Duration) {
	var sum float64
	var count float64
	var avg float64
	for _, d := range durations {
		if d > 0 {
			sum += float64(d)
			count++
		}
	}
	if count > 0 {
		avg = sum / count
	}
	log.Printf("Event completed - clients %.f, avg time %.2f\n", count, avg)
}
