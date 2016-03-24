package main

import (
	"log"
	"time"

	"github.com/luizbranco/eventsource"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/cloudmonitoring/v2beta2"
)

type GoogleCloudMonitoring struct {
	base_url string
	remote   cloudmonitoring.TimeseriesService
}

type Timeseries struct {
	base_url    string
	metric_name string
	start       time.Time
	end         time.Time
	value       float64
}

func createTimeseries(args Timeseries) *cloudmonitoring.TimeseriesPoint {

	var end_string string
	var start_string string

	start_string = args.start.Format(time.RFC3339)

	if &args.end != nil {
		end_string = args.end.Format(time.RFC3339)
	} else {
		end_string = start_string
	}

	description := cloudmonitoring.TimeseriesDescriptor{
		Labels: map[string]string{
			args.base_url + "implementation": "golang",
		},
		Metric:  args.base_url + args.metric_name,
		Project: "replay-gaming",
	}

	point := cloudmonitoring.Point{
		Start:       start_string,
		End:         end_string,
		DoubleValue: &args.value,
	}

	timeseries := cloudmonitoring.TimeseriesPoint{
		Point:          &point,
		TimeseriesDesc: &description,
	}

	return &timeseries
}

func pushMetrics(points []*cloudmonitoring.TimeseriesPoint, remote cloudmonitoring.TimeseriesService) {
	request := cloudmonitoring.WriteTimeseriesRequest{
		CommonLabels: map[string]string{
			"container.googleapis.com/container_name": "eventsource",
		},
		Timeseries: points,
	}

	response, err := remote.Write("replay-gaming", &request).Do()
	if err != nil {
		log.Printf("pushMetrics - Unable to write timeseries: %v", err)
	} else {
		log.Printf("pushMetrics - Response: %s", response)
	}
}

func createMetricDescriptor(prefix string, name string, description string) *cloudmonitoring.MetricDescriptor {
	metric_type := cloudmonitoring.MetricDescriptorTypeDescriptor{
		MetricType: "gauge",
		ValueType:  "double",
	}

	label := cloudmonitoring.MetricDescriptorLabelDescriptor{
		Description: "Application",
		Key:         "eventsource",
	}

	return &cloudmonitoring.MetricDescriptor{
		Description: description,
		Labels: []*cloudmonitoring.MetricDescriptorLabelDescriptor{
			&label,
		},
		Name:           "custom.cloudmonitoring.googleapis.com/" + prefix + "/" + name,
		Project:        "replay-gaming",
		TypeDescriptor: &metric_type,
	}
}

func createMetric(prefix string, metricDescriptorsService *cloudmonitoring.MetricDescriptorsService, name string, description string) {
	metricDescriptor := createMetricDescriptor(prefix, name, description)
	metricDescriptor, err := metricDescriptorsService.Create("replay-gaming", metricDescriptor).Do()
	if err != nil {
		log.Printf("Unable to create '%s' metric: %v", name, err)
	}
}

func createMetrics(prefix string, cloudmonitoringService *cloudmonitoring.Service) {
	metricDescriptorsService := cloudmonitoring.NewMetricDescriptorsService(cloudmonitoringService)

	createMetric(prefix, metricDescriptorsService, "clients", "Number of clients that the event was distributed to")
	createMetric(prefix, metricDescriptorsService, "avg_time", "Average time to send an event to all connected clients")
	createMetric(prefix, metricDescriptorsService, "connections", "Number of connected SSE clients (browser sessions)")
}

func NewMetrics(prefix string) (GoogleCloudMonitoring, error) {
	monitor := GoogleCloudMonitoring{base_url: "custom.cloudmonitoring.googleapis.com/" + prefix + "/"}

	client, err := google.DefaultClient(
		context.Background(),
		cloudmonitoring.MonitoringScope,
		cloudmonitoring.CloudPlatformScope,
	)
	if err != nil {
		log.Printf("Unable to get default client: %v", err)
	}

	cloudmonitoringService, err := cloudmonitoring.New(client)
	if err != nil {
		log.Printf("Unable to create monitoring service: %v", err)
	}

	createMetrics(prefix, cloudmonitoringService)

	timeseriesService := cloudmonitoring.NewTimeseriesService(cloudmonitoringService)

	monitor.remote = *timeseriesService
	return monitor, nil
}

func (monitor GoogleCloudMonitoring) ClientCount(count int) {
	log.Printf("[METRIC] %sconnections: %d\n", monitor.base_url, count)

	timeseries := createTimeseries(Timeseries{
		base_url:    monitor.base_url,
		metric_name: "connections",
		start:       time.Now().UTC(),
		value:       float64(count),
	})

	points := []*cloudmonitoring.TimeseriesPoint{
		timeseries,
	}

	pushMetrics(points, monitor.remote)
}

func (monitor GoogleCloudMonitoring) EventDone(event eventsource.Event, duration time.Duration, eventdurations []time.Duration) {
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

	clients_timeseries := createTimeseries(Timeseries{
		base_url:    monitor.base_url,
		metric_name: "clients",
		start:       time.Now().UTC(),
		value:       float64(count),
	})

	avg_time_timeseries := createTimeseries(Timeseries{
		base_url:    monitor.base_url,
		metric_name: "avg_time",
		start:       time.Now().UTC(),
		value:       avg,
	})

	points := []*cloudmonitoring.TimeseriesPoint{
		clients_timeseries,
		avg_time_timeseries,
	}

	pushMetrics(points, monitor.remote)
}
