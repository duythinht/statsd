package statsd

import (
	"testing"
	"time"
)

func not(t *testing.T, notExpected interface{}, actual interface{}, message string) {
	if notExpected == actual {
		t.Error(message, "Expected: not", notExpected, "Actually:", actual)
	}
}

func shouldBe(t *testing.T, expected interface{}, actual interface{}, message string) {
	if expected != actual {
		t.Error(message, "Expected: ", expected, "Actually:", actual)
	}
}

func TestParseImplements(t *testing.T) {
	line := []byte("test.metric:1|c")
	metric, err := ParseLine(line)

	if err != nil {
		t.Error("Error should not raised")
	}

	if metric.Key != "test.metric" {
		t.Error("Return wrong metric key")
	}

	if metric.Value != "1" {
		t.Error("Wrong metric value")
	}

	if metric.Type != "counter" {
		t.Error("Wrong metric type, expected", "c", "got", metric.Type)
	}

	if metric.SampleRate != 1.0 {
		t.Error("Wrong Sampling value, expected", 1.0, "got", metric.SampleRate)
	}
}

func TestParseLineShouldBeCorrect(t *testing.T) {
	lines := []string{
		"test.right:1|c",
		"test.right:2|ms",
		"test.right:3|g",
		"test.right:3|s",
		"test.right:1|c|@0.1",
		"test.right:1|ms|@0.9",
	}

	for _, line := range lines {
		_, err := ParseLine([]byte(line))
		shouldBe(t, nil, err, "Should be not error for"+line)
	}

}

func TestParseWrongFormat(t *testing.T) {
	lines := []string{
		"test.wrong:123:123",
		"wrong.parseSampleRate:1|g|@c",
		"wrong.sample.type:123|c|#abc",
	}

	for _, line := range lines {
		_, err := ParseLine([]byte(line))
		not(t, nil, err, "Should not nil with line"+line)
	}
}

func TestPerformance(t *testing.T) {
	line := []byte("test.right:1|c|@0.1")
	b := time.Now().UnixNano()

	for i := 0; i < 1000000; i++ {
		StatsdLine(line)
	}
	e := time.Now().UnixNano()
	if (e-b)/1000000 > 1000 {
		t.Error("Performance is bad")
	}
}
