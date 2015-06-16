package main

import (
	"github.com/luizbranco/eventsource"
	"github.com/peterbourgon/g2s"
)

const sampleRate = 0.25

type Stats struct {
	prefix string
	statsd g2s.Statter
}

func NewStats(addr string, prefix string) (*Stats, error) {
	stats := &Stats{prefix: prefix + "."}
	statsd, err := g2s.Dial("udp", addr)
	if err != nil {
		return stats, err
	}
	stats.statsd = statsd
	return stats, nil
}

func (s *Stats) ClientsCount(count int) {
	s.statsd.Gauge(1, s.prefix+"connections", string(count))
}

func (s *Stats) EventSent(stats eventsource.EventStats) {
	duration := stats.End.Sub(stats.Start)
	s.statsd.Timing(sampleRate, s.prefix+"publish.timing", duration)
}

func (s *Stats) EventEnd(stats eventsource.EventStats) {
	duration := stats.End.Sub(stats.Start)
	s.statsd.Timing(sampleRate, s.prefix+"publish.connection_write.timing", duration)
}
