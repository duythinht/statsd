package config

type BasicConfig struct {
	InfluxURL      string
	DefaultDB      string
	InfluxUsername string
	InfluxPassword string
	DefaultBucket  string
	Listen         string
}

type Telemetry struct {
	DB   string
	Keys []string
}
