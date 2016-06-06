package statsd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatsd(t *testing.T) {
	metrics := []string{
		"test.sample.counter:1|c",
		"test.sample.timing:2|ms",
		"test.sample.set:1|s",
		"test.sample.gauge:+1|g",
	}

	for _, m := range metrics {
		err := StatsdLine([]byte(m))
		assert.Nil(t, err, m+" has been failed to parse")
	}
}

func TestTiming(t *testing.T) {
	metrics := []string{
		"test.value.timing:1|ms",
		"test.value.timing:2|ms",
		"test.value.timing:3|ms",
	}

	for _, m := range metrics {
		err := StatsdLine([]byte(m))
		assert.Nil(t, err, m+" has been failed to parse")
	}

	assert.Equal(t, 1, timings["test.value.timing"].FirstVal)
	assert.Equal(t, 3, timings["test.value.timing"].Count)
	assert.Equal(t, 3, timings["test.value.timing"].Sum)
	assert.Equal(t, 5, timings["test.value.timing"].SumSq)
	assert.Equal(t, 3, timings["test.value.timing"].Upper)
	assert.Equal(t, 1, timings["test.value.timing"].Lower)
}

func TestCounter(t *testing.T) {
	metrics := []string{
		"test.value.counter:1|c",
		"test.value.counter:1|c|@0.1",
	}

	for _, m := range metrics {
		err := StatsdLine([]byte(m))
		assert.Nil(t, err, m+" has been failed to parse")
	}

	assert.Equal(t, 11, counters["test.value.counter"].Value)
}

func TestSet(t *testing.T) {
	metrics := []string{
		"test.value.set:1|s",
		"test.value.set:1|s",
		"test.value.set:2|s",
	}

	for _, m := range metrics {
		err := StatsdLine([]byte(m))
		assert.Nil(t, err, m+" has been failed to parse")
	}

	assert.Equal(t, 2, sets["test.value.set"].Value)

}

func TestGause(t *testing.T) {
	err := StatsdLine([]byte("test:+10|g"))
	assert.Nil(t, err, "Error should be nil")

	gs := []string{
		"test.g.value:1|g",
		"test.g.value:+2|g",
		"test.g.value:+3|g",
		"test.g.value:+4|g",
		"test.g.value:-5|g",
	}
	for _, m := range gs {
		StatsdLine([]byte(m))
	}
	assert.Equal(t, 5, gauges["test.g.value"].Value)
}
