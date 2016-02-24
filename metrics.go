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
}

func (monitor GoogleCloudMonitoring) EventDone(event eventsource.Event, duration time.Duration, eventdurations []time.Duration) {
	// s.statsd.Counter(1, s.prefix+"publish.count", 1)
	// s.statsd.Timing(sampleRate, s.prefix+"publish.timing", d*roundUpMs)
	// for _, t := range c {
	// 	if t > 0 {
	// 		s.statsd.Timing(sampleRate, s.prefix+"publish.connection_write.timing", t*roundUpMs)
	// 	}
	// }

	var sum float64
	var count float64
	var avg float64
	for _, d := range eventdurations {
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

//// If you need a oauth2.TokenSource, use the DefaultTokenSource function:

//ts, err := google.DefaultTokenSource(ctx, scope1, scope2, ...)
//if err != nil {
//  // Handle error.
//}
//httpClient := oauth2.NewClient(ctx, ts)
