package main

import (
	"time"

	"github.com/luizbranco/eventsource"
	"github.com/quipo/statsd"
)

type Stats struct {
	graphite *statsd.StatsdBuffer
}

func NewStats(addr string, prefix string) *Stats {
	stats := &Stats{}
	interval := time.Second * 2 // aggregate stats and flush every 2 seconds
	statsdclient := statsd.NewStatsdClient(addr, prefix)
	stats.graphite = statsd.NewStatsdBuffer(interval, statsdclient)
	stats.graphite.CreateSocket()
	return stats
}

func (s *Stats) ClientsCount(count int) {
	s.graphite.Gauge("connections", int64(count))
}

func (s *Stats) EventSent(stats eventsource.EventStats) {
	duration := stats.End.Sub(stats.Start)
	s.graphite.PrecisionTiming("publish.timing", duration)
}

func (s *Stats) EventEnd(stats eventsource.EventStats) {
	duration := stats.End.Sub(stats.Start)
	s.graphite.PrecisionTiming("publish.connection_write.timing", duration)
}
