package statsd

import (
	"math"
	"sort"
)

type Point struct {
	DB     string
	Bucket string
	Tags   map[string]string
	Fields map[string]interface{}
}

func Aggregate(points *[]Point) {
	mutex.Lock()
	defer mutex.Unlock()

	ProcessCounter(points)
	ProcessTiming(points)
	ProcessGauge(points)
	ProcessSet(points)
	//Then reset them
	counters = map[string]*counter{}
	gauges = map[string]*gauge{}
	timings = map[string]*timing{}
	sets = map[string]*set{}
}

func ProcessCounter(points *[]Point) {
	for _, c := range counters {
		*points = append(*points, Point{
			c.DB,
			c.Bucket,
			c.Tags,
			map[string]interface{}{
				"value": c.Value,
			},
		})
	}
}

func ProcessGauge(points *[]Point) {
	for _, c := range gauges {
		*points = append(*points, Point{
			c.DB,
			c.Bucket,
			c.Tags,
			map[string]interface{}{
				"value": c.Value,
			},
		})
	}
}

func ProcessSet(points *[]Point) {
	for _, c := range sets {
		*points = append(*points, Point{
			c.DB,
			c.Bucket,
			c.Tags,
			map[string]interface{}{
				"value": c.Value,
			},
		})
	}
}

func ProcessTiming(points *[]Point) {
	for _, c := range timings {
		mean := c.FirstVal + (c.Sum / c.Count)
		variance := (float64(c.SumSq) - float64(c.Sum*c.Sum)/float64(c.Count)) / float64(c.Count)
		stddev := math.Sqrt(variance)
		sort.Ints(c.data)
		percentile_90 := c.data[c.Count*9/10]
		*points = append(*points, Point{
			c.DB,
			c.Bucket,
			c.Tags,
			map[string]interface{}{
				"mean":          mean,
				"stddev":        stddev,
				"upper":         c.Upper,
				"lower":         c.Lower,
				"count":         c.Count,
				"percentile_90": percentile_90,
			},
		})
	}
}
