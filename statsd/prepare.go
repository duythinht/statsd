package statsd

import (
	"errors"
	"github.com/duythinht/statsd/config"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func PrepareTiming(metric Metric) error {

	value, err := strconv.Atoi(metric.Value)

	if err != nil {
		return err
	}

	if value < 0 {
		return errors.New("Value should be not negative")
	}

	if timing, ok := timings[metric.Key]; ok {
		return appendTiming(timing, value)
	}

	statsd, err := PrepareStatd(metric)

	if err != nil {
		return err
	}

	timings[metric.Key] = &timing{
		*statsd,
		value,
		1,
		0,
		0,
		value,
		value,
		[]int{value},
	}
	return nil
}

func appendTiming(t *timing, value int) error {

	t.Count += 1
	// https://en.wikipedia.org/wiki/Algorithms_for_calculating_variance
	x := value - t.FirstVal
	t.Sum += x
	t.SumSq += x * x

	if t.Upper < value {
		t.Upper = value
	}

	if t.Lower > value {
		t.Lower = value
	}

	t.data = append(t.data, value)

	return nil
}

func PrepareCounter(metric Metric) error {
	v, err := strconv.Atoi(metric.Value)
	if err != nil {
		return err
	}
	value := v * int(1.0/metric.SampleRate)

	if counter, ok := counters[metric.Key]; ok {
		counter.Value += value
		return nil
	}

	statsd, err := PrepareStatd(metric)

	if err != nil {
		return err
	}

	counters[metric.Key] = &counter{
		*statsd,
		value,
	}
	return nil
}

func PrepareGauge(metric Metric) error {

	value, err := strconv.Atoi(metric.Value)

	if err != nil {
		return err
	}

	if gauge, ok := gauges[metric.Key]; ok {
		if _, err := strconv.Atoi(metric.Value[:1]); err == nil {
			gauge.Value = value
		} else {
			gauge.Value += value
			//gauge.Value += value
		}
		return nil
	}

	statsd, err := PrepareStatd(metric)

	if err != nil {
		return err
	}

	gauges[metric.Key] = &gauge{
		*statsd,
		value,
	}
	return nil
}

func PrepareSet(metric Metric) error {

	if set, ok := sets[metric.Key]; ok {
		if _, exist := set.data[metric.Value]; exist {
			return nil
		}
		set.Value += 1
		set.data[metric.Value] = true
		return nil
	}

	statsd, err := PrepareStatd(metric)

	if err != nil {
		return err
	}

	sets[metric.Key] = &set{
		*statsd,
		1,
		map[string]bool{metric.Value: true},
	}

	return nil
}

func PrepareStatd(metric Metric) (*statsd, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	tags := map[string]string{
		"host":        hostname,
		"metric_type": metric.Type,
	}

	db := config.Root.DefaultDB

	bucket := strings.Split(metric.Key, ".")[0]

	for _, telemetry := range config.Telemetries {
		for _, key := range telemetry.Keys {
			re := regexp.MustCompile(key)
			db = telemetry.DB
			match := re.FindStringSubmatch(metric.Key)
			if len(match) > 0 {
				i := 1
				listTags := re.SubexpNames()
				for i <= re.NumSubexp() {
					if listTags[i] == "bucket" {
						bucket = match[i]
					} else {
						tags[listTags[i]] = match[i]
					}
					i += 1
				}
				return &statsd{db, bucket, tags}, nil
			}
		}
	}

	return nil, errors.New(metric.Key + " is not match any pattern")
}
