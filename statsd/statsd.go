package statsd

import (
	"errors"
	"sync"
)

var mutex = sync.Mutex{}

type statsd struct {
	DB     string
	Bucket string
	Tags   map[string]string
}

type timing struct {
	statsd
	FirstVal int
	Count    int
	Sum      int
	SumSq    int
	Upper    int
	Lower    int
	data     []int
}

type gauge struct {
	statsd
	Value int
}

type counter struct {
	statsd
	Value int
}

type set struct {
	statsd
	Value int
	data  map[string]bool
}

var (
	counters = map[string]*counter{}
	gauges   = map[string]*gauge{}
	timings  = map[string]*timing{}
	sets     = map[string]*set{}
)

func StatsdLine(line []byte) error {
	mutex.Lock()
	defer mutex.Unlock()
	metric, err := ParseLine(line)
	if err != nil {
		return err
	}

	switch metric.Type {
	case "timing":
		for i := 0; i < int(1.0/metric.SampleRate); i++ {
			err = PrepareTiming(*metric)
			if err != nil {
				return err
			}
		}
	case "counter":
		return PrepareCounter(*metric)
	case "gauge":
		return PrepareGauge(*metric)
	case "set":
		return PrepareSet(*metric)
	default:
		return errors.New("Unknown metric type: " + metric.Type)
	}
	return nil
}
