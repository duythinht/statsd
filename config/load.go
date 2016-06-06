package config

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/hcl"
	"io/ioutil"
	"log"
)

var Root = BasicConfig{}
var Telemetries = []Telemetry{}

func Load(configFile string) {
	fdata, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Fatalln(err)
	}

	tree := map[string]interface{}{}

	hcl.Decode(&tree, string(fdata))

	if dfdb, ok := tree["default_db"]; ok {
		Root.DefaultDB = dfdb.(string)
	} else {
		log.Fatalln("Must have config default_db")
	}

	if influx_url, ok := tree["influx_url"]; ok {
		Root.InfluxURL = influx_url.(string)
	} else {
		log.Fatalln("Must has config influx_url")
	}

	if listen, ok := tree["listen"]; ok {
		Root.Listen = listen.(string)
	} else {
		log.Fatalln("Must has config listen")
	}

	if influx_username, ok := tree["influx_username"]; ok {
		Root.InfluxUsername = influx_username.(string)
	} else {
		Root.InfluxUsername = "admin"
	}
	if influx_pass, ok := tree["influx_password"]; ok {
		Root.InfluxPassword = influx_pass.(string)
	} else {
		Root.InfluxPassword = "admin"
	}

	if telemetries, ok := tree["telemetry"]; ok {
		for _, telemetry := range telemetries.([]map[string]interface{}) {

			// Process raw keys
			rks := telemetry["keys"].([]interface{})
			if len(rks) < 0 {
				log.Fatal("telemetry for", telemetry["db"], "must be specs")
			}
			keys := []string{}

			for _, rk := range rks {
				keys = append(keys, rk.(string))
			}

			Telemetries = append(Telemetries, Telemetry{
				telemetry["db"].(string),
				keys,
			})
		}
	}

	if err != nil {
		log.Fatalln(err)
	}

	log.Println(spew.Sdump(Root))
	log.Println(spew.Sdump(Telemetries))

}
