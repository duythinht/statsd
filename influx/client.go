package influx

import (
	"github.com/duythinht/statsd/config"
	"github.com/duythinht/statsd/statsd"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"net/url"
	"time"
)

type Client struct {
	client  client.Client
	batches map[string]client.BatchPoints
}

func CreateClient() Client {
	uri, err := url.Parse(config.Root.InfluxURL)
	if err != nil {
		log.Fatalln(err)
	}

	switch uri.Scheme {
	case "udp":
		clnt, err := client.NewUDPClient(client.UDPConfig{
			Addr: uri.Host,
		})
		if err != nil {
			log.Fatalln(err)
		}

		return Client{clnt, map[string]client.BatchPoints{}}

	case "http":
		clnt, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     config.Root.InfluxURL,
			Username: config.Root.InfluxUsername,
			Password: config.Root.InfluxPassword,
		})

		if err != nil {
			log.Fatalln(err)
		}

		return Client{clnt, map[string]client.BatchPoints{}}

	default:
		log.Fatalln("Wrong protocol of influxdb")
	}
	return Client{}
}

func (c *Client) Init() {
	err := c.CreateDB(config.Root.DefaultDB)

	if err != nil {
		log.Fatalln(err)
	}

	for _, telemetry := range config.Telemetries {
		c.CreateDB(telemetry.DB)
	}
}

func (c *Client) CreateDB(db string) error {
	log.Println("CreateDB", db)
	q := client.Query{
		Command: "CREATE DATABASE IF NOT EXISTS " + db,
	}
	_, err := c.client.Query(q)
	return err
}

func (c *Client) Submit(points []statsd.Point) int {

	batches := map[string]client.BatchPoints{}

	for _, point := range points {
		pt, _ := client.NewPoint(point.Bucket, point.Tags, point.Fields, time.Now())

		var batch client.BatchPoints

		if bpoint, ok := batches[point.DB]; ok {
			batch = bpoint
		} else {

			newbatch, err := client.NewBatchPoints(client.BatchPointsConfig{
				Database:  point.DB,
				Precision: "s",
			})

			if err != nil {
				log.Fatal(err)
			}
			batch = newbatch
			batches[point.DB] = batch
		}
		batch.AddPoint(pt)
	}

	for _, batch := range batches {
		e := c.client.Write(batch)
		if e != nil {
			log.Println("Error when write to db:", e)
		}
	}
	return len(points)
}

func (c *Client) Ping() (time.Duration, string, error) {
	return c.client.Ping(5 * time.Second)
}
