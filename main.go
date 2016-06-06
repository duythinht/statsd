package main

import (
	"bytes"
	"flag"
	"github.com/duythinht/statsd/config"
	"github.com/duythinht/statsd/influx"
	"github.com/duythinht/statsd/statsd"
	"log"
	"net"
	"time"
)

const BUFFER_SIZE = 65536

func init() {
	var configFile string
	flag.StringVar(&configFile, "config", "/etc/statsd.hcl", "Need a config")
	flag.Parse()
	config.Load(configFile)
}

func main() {
	log.Println("Start UDP")
	buffer := make([]byte, BUFFER_SIZE)

	address, _ := net.ResolveUDPAddr("udp", config.Root.Listen)

	log.Printf("listening on %s", address)
	listener, err := net.ListenUDP("udp", address)

	if err != nil {
		log.Println(err)
	}
	// Make aggregate and submit server
	go func() {
		client := influx.CreateClient()
		log.Println(client.Ping())
		client.Init()
		for {
			time.Sleep(10 * time.Second)
			points := []statsd.Point{}
			statsd.Aggregate(&points)
			//log.Println(points)
			num := client.Submit(points)
			log.Println("Wrote", num, "metrics")
		}
	}()

	for {
		num, _ := listener.Read(buffer)
		log.Println("Gather metrics key:", num, string(buffer[:num]))
		lines := bytes.Split(buffer[:num], []byte("\n"))

		for _, line := range lines {
			if len(line) > 0 {
				err := statsd.StatsdLine(line)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}
