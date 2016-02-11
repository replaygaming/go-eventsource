package main

import (
	"./cloud_metrics"
)

func main() {

	metrics, err := NewMetricAgent("eventsource")

	metrics.remote.ClientCount( 5 )
}
