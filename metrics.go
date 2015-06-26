package main

import (
	"strconv"
	"time"

	"github.com/luizbranco/eventsource"
	"github.com/peterbourgon/g2s"
)

const sampleRate = 0.25
const roundUpMs = 1e3

type StatdsD struct {
	prefix string
	statsd g2s.Statter
}

func NewMetrics(addr string, prefix string) (StatdsD, error) {
	stats := StatdsD{prefix: prefix + "."}
	statsd, err := g2s.Dial("udp", addr)
	if err != nil {
		return stats, err
	}
	stats.statsd = statsd
	return stats, nil
}

func (s StatdsD) ClientCount(count int) {
	n := strconv.Itoa(count)
	s.statsd.Gauge(1, s.prefix+"connections", n)
}

func (s StatdsD) EventDone(e eventsource.Event, d time.Duration, c []time.Duration) {
	s.statsd.Counter(1, s.prefix+"publish.count", 1)
	s.statsd.Timing(sampleRate, s.prefix+"publish.timing", d*roundUpMs)
	for _, t := range c {
		if t > 0 {
			s.statsd.Timing(sampleRate, s.prefix+"publish.connection_write.timing", t*roundUpMs)
		}
	}
}
