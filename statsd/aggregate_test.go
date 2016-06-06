package statsd

import (
	"testing"
)

func TestCounterProcess(t *testing.T) {
	points := []Point{}
	metrics := []string{
		"test.aggregate.counter:1|c",
		"test.aggregate.counter:2|c",
		"test.aggregate.counter:1|c",
		"test.aggregate.gauge:1|g",
	}

	for _, m := range metrics {
		StatsdLine([]byte(m))
	}

	ProcessCounter(&points)
	if len(points) != 1 {
		t.Error("len point is", len(points))
	}
	if points[0].Fields["value"] != 4 {
		t.Error("Wrong value")
	}
}

func TestTimingProcess(t *testing.T) {
	points := []Point{}
	metrics := []string{
		"test.aggregate.timing:1|ms",
		"test.aggregate.timing:2|ms",
		"test.aggregate.timing:3|ms",
		"test.aggregate.timing:4|ms",
		"test.aggregate.timing:5|ms",
		"test.aggregate.timing:6|ms",
		"test.aggregate.timing:7|ms",
		"test.aggregate.timing:8|ms",
		"test.aggregate.timing:9|ms",
		"test.aggregate.timing:10|ms",
		"test.aggregate.timing:11|ms",
		"test.aggregate.timing:12|ms",
	}

	for _, m := range metrics {
		StatsdLine([]byte(m))
	}

	ProcessTiming(&points)
}
