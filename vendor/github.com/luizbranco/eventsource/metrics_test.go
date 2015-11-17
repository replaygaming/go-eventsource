package eventsource

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"
)

func TestDefaultMetricsEventDone(t *testing.T) {
	buf := new(bytes.Buffer)
	log.SetOutput(buf)
	m := DefaultMetrics{}
	e := DefaultEvent{}
	var d time.Duration
	c := []time.Duration{1, 2, 3, 4, 5, 6}
	m.EventDone(e, d, c)
	line := buf.String()
	result := line[0 : len(line)-1]
	expecting := "Event completed - clients 6, avg time 3.50"
	if !strings.Contains(result, expecting) {
		t.Errorf("expected:\n%s\nto be contained in:\n%s\n", expecting, result)
	}
}
