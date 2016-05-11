package main

import (
	"log"
	"time"

	"github.com/replaygaming/eventsource"
)

type LogMonitoring struct {
	prefix string
}

func NewMetrics(prefix string) (LogMonitoring, error) {
	monitor := LogMonitoring{prefix: prefix + "/"}

	return monitor, nil
}

func (monitor LogMonitoring) ClientCount(count int) {
	log.Printf("[METRIC] %sconnections: %d\n", monitor.base_url, count)
}

func (monitor LogMonitoring) EventDone(event eventsource.Event, duration time.Duration, eventdurations []time.Duration) {
	var sum int64
	var count int64
	var avg float64

	for _, d := range eventdurations {
		if d > 0 {
			sum += d.Nanoseconds()
		}
	}

	count = int64(len(eventdurations))

	if count > 0 {
		avg = float64(sum) / float64(count)
	}

	log.Printf("[METRIC] %s.event_distributed.clients: %d\n", monitor.base_url, count)
	log.Printf("[METRIC] %s.event_distributed.avg_time: %dns\n", monitor.base_url, avg)
}
