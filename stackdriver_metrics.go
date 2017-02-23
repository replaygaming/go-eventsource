package main

import (
	"log"
	"math/rand"
	"time"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	googlepb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/replaygaming/eventsource"
	"golang.org/x/net/context"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

const projectID = "replay-gaming"
const bufferSize = 10

type StackdriverMetrics struct {
	prefix           string
	subscriptionName string
	client           monitoring.MetricClient
	rand             *rand.Rand
}

type simpleDataPoint interface {
	toTypedValue() *monitoringpb.TypedValue
	getTime() int64
}

type int64DataPoint struct {
	value int64
	time  int64
}

type float64DataPoint struct {
	value float64
	time  int64
}

func (dataPoint int64DataPoint) getTime() int64 {
	return dataPoint.time
}

func (dataPoint int64DataPoint) toTypedValue() *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_Int64Value{
			Int64Value: dataPoint.value,
		},
	}
}

func (dataPoint float64DataPoint) getTime() int64 {
	return dataPoint.time
}

func (dataPoint float64DataPoint) toTypedValue() *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_DoubleValue{
			DoubleValue: dataPoint.value,
		},
	}
}

func NewStackdriverMetrics(prefix string, subscriptionName string) (StackdriverMetrics, error) {
	monitor := StackdriverMetrics{
		prefix:           prefix,
		subscriptionName: subscriptionName,
		rand:             rand.New(rand.NewSource(time.Now().Unix())),
	}

	// Use oauth2.NoContext if there isn't a good context to pass in.
	ctx := context.Background()

	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	monitor.client = *client

	return monitor, nil
}

func (monitor StackdriverMetrics) ClientCount(count int) {
	monitor.sendDataPoint("client_count", &int64DataPoint{
		value: int64(count),
		time:  time.Now().Unix(),
	}, 1)

	log.Printf("[METRIC] Current connections: %d\n", count)
}

func (monitor StackdriverMetrics) EventDone(event eventsource.Event, duration time.Duration, eventdurations []time.Duration) {
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

	log.Printf("[METRIC] Event completed - clients %.f, avg time %.2f ns\n", count, avg)

	monitor.sendDataPoint("client_message_average_time", float64DataPoint{
		value: avg,
		time:  time.Now().Unix(),
	}, 0.3)
}

func (monitor StackdriverMetrics) sendDataPoint(metric string, dataPoint simpleDataPoint, sampleRate float32) {
	// optionally limit amount of data we send
	if monitor.rand.Float32() > sampleRate {
		return
	}

	log.Printf("[METRIC] Sending time series for %s", metric)

	err := monitor.client.CreateTimeSeries(context.Background(), &monitoringpb.CreateTimeSeriesRequest{
		Name: monitoring.MetricProjectPath(projectID),
		TimeSeries: []*monitoringpb.TimeSeries{
			{
				Metric: &metricpb.Metric{
					Type: "custom.googleapis.com/eventsource/" + monitor.prefix + "/" + metric,
					Labels: map[string]string{
						"subscription_name": monitor.subscriptionName,
					},
				},
				Resource: &monitoredrespb.MonitoredResource{
					Type: "global",
					Labels: map[string]string{
						"project_id": projectID,
					},
				},
				Points: []*monitoringpb.Point{
					{
						Interval: &monitoringpb.TimeInterval{
							EndTime: &googlepb.Timestamp{
								Seconds: dataPoint.getTime(),
							},
						},
						Value: dataPoint.toTypedValue(),
					},
				},
			},
		},
	})

	if err != nil {
		log.Printf("[METRIC] [ERROR] Error sending time series: %v", err)
	}
}
