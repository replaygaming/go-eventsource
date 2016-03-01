package main

import (
	"log"
	"time"

	"github.com/luizbranco/eventsource"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudmonitoring/v2beta2"
)

//const sampleRate = 0.25
//const roundUpMs = 1e3

type GoogleCloudMonitoring struct {
	base_url string
	remote   cloudmonitoring.Service
}

func NewMetrics(prefix string) (GoogleCloudMonitoring, error) {
	monitor := GoogleCloudMonitoring{base_url: "custom.cloudmonitoring.googleapis.com/" + prefix + "/"}

	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.Background()

	client, err := google.DefaultClient(
		ctx,
		cloudmonitoring.MonitoringScope,
	)
	if err != nil {
		log.Fatal(err)
	}

	cloudmonitoringService, err := cloudmonitoring.New(client)
	if err != nil {
		log.Fatal(err)
	}

	monitor.remote = *cloudmonitoringService
	return monitor, nil
}

func (monitor GoogleCloudMonitoring) ClientCount(count int) {
	// n := strconv.Itoa(count)
	// s.statsd.Gauge(1, s.prefix+"connections", n)
	log.Printf("[METRIC] %s.connections: %d\n", monitor.base_url, count)
}

func (monitor GoogleCloudMonitoring) EventDone(event eventsource.Event, duration time.Duration, eventdurations []time.Duration) {
	// s.statsd.Counter(1, s.prefix+"publish.count", 1)
	// s.statsd.Timing(sampleRate, s.prefix+"publish.timing", d*roundUpMs)
	// for _, t := range c {
	// 	if t > 0 {
	// 		s.statsd.Timing(sampleRate, s.prefix+"publish.connection_write.timing", t*roundUpMs)
	// 	}
	// }

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

//// If you need a oauth2.TokenSource, use the DefaultTokenSource function:

//ts, err := google.DefaultTokenSource(ctx, scope1, scope2, ...)
//if err != nil {
//  // Handle error.
//}
//httpClient := oauth2.NewClient(ctx, ts)
