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
	remote   cloudmonitoring.TimeseriesService
}

func NewMetrics(prefix string) (GoogleCloudMonitoring, error) {
	monitor := GoogleCloudMonitoring{base_url: "custom.cloudmonitoring.googleapis.com/" + prefix + "/"}

	client, err := google.DefaultClient(
		context.Background(),
		cloudmonitoring.MonitoringScope,
	)
	if err != nil {
		log.Fatal("Unable to get default client: %v", err)
	}

	cloudmonitoringService, err := cloudmonitoring.New(client)
	if err != nil {
		log.Fatal("Unable to create monitoring service: %v", err)
	}

	timeseriesService := cloudmonitoring.NewTimeseriesService(cloudmonitoringService)

	monitor.remote = *timeseriesService
	return monitor, nil
}

func (monitor GoogleCloudMonitoring) ClientCount(count int) {
	// n := strconv.Itoa(count)
	// s.statsd.Gauge(1, s.prefix+"connections", n)
	log.Printf("[METRIC] %sconnections: %d\n", monitor.base_url, count)

	count64 := int64( count )
	now, err := time.Now().UTC().MarshalText()
	if err != nil {
		log.Fatal("ClientCount - Unable to get current time: %v", err)
	}

	description := cloudmonitoring.TimeseriesDescriptor{
		Labels: map[string]string{
			monitor.base_url + "implementation":"golang",
		},
    Metric: monitor.base_url + "connections",
    Project: "replay-gaming",
	}

	point := cloudmonitoring.Point{
		Start: string(now),
		End: string(now),
		Int64Value: &count64,
	}

	timeseries := cloudmonitoring.TimeseriesPoint{
    Point: &point,
  	TimeseriesDesc: &description,
	}

	points := []*cloudmonitoring.TimeseriesPoint{
		&timeseries,
	}

	request := cloudmonitoring.WriteTimeseriesRequest{
		CommonLabels: map[string]string{
			"container.googleapis.com/container_name": "eventsource",
		},
		Timeseries: points,
	}

	response, err := monitor.remote.Write("replay-poker", &request).Do()
	if err != nil {
		log.Fatal("ClientCount - Unable to write timeseries: %v", err)
	}
	log.Printf("ClientCount - Response: %s", response)
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
