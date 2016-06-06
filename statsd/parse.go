package statsd

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
)

const VAL_SEPARATOR = ":"
const TYPE_SEPARATOR = "|"

type Metric struct {
	Key        string
	Value      string
	Type       string
	SampleRate float64
}

var re = regexp.MustCompile("\\d+")

func ParseLine(line []byte) (*Metric, error) {
	bits := bytes.Split(line, []byte(VAL_SEPARATOR))

	if len(bits) != 2 {
		return nil, errors.New("Failed to parse line:" + string(line))
	}

	key, rest := bits[0], bits[1]
	key = re.ReplaceAll(key, []byte("{number}"))

	segments := bytes.Split(rest, []byte(TYPE_SEPARATOR))

	//"123|c|@0.1" -> 3
	//"123|c" -> 2

	segLength := len(segments)

	if segLength > 3 || segLength < 2 {
		return nil, errors.New("Wrong metric format:" + string(line))
	}

	value := string(segments[0])

	metricType := getMetricType(string(segments[1]))

	// Get sampling rate when segments length == 3
	if segLength == 3 {
		sampleRate, err := parseSampleRate(segments[2])
		if err != nil {
			return nil, err
		}
		return &Metric{string(key), value, metricType, sampleRate}, nil
	}

	return &Metric{string(key), value, metricType, 1.0}, nil
}

func getMetricType(sign string) string {
	switch sign {
	case "ms":
		return "timing"
	case "c":
		return "counter"
	case "g":
		return "gauge"
	case "s":
		return "set"
	default:
		return "unknown"
	}
}

func parseSampleRate(segment []byte) (float64, error) {
	at, sample := segment[:1], segment[1:]
	if string(at) != "@" {
		return 0, errors.New("Wrong format for sampling: " + string(at) + string(sample))
	}
	rate, err := strconv.ParseFloat(string(sample), 64)

	if err != nil {
		return 0, errors.New("Wrong type of sampling rate: " + string(at) + string(sample))
	}

	if rate > 1.0 {
		rate = 1.0
	}

	return rate, nil
}
