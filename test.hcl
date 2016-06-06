influx_url = "http://localhost:8086"

/*influx_url = "udp://192.168.99.100:8089"*/

influx_username = "admin"

influx_password = "password"

listen = ":8127"

default_db = "statsd"

telemetry {
  db   = "nomad"
  keys = ["^test.(?P<bucket>hello).(?P<name>[^ ]*)$", "def"]
}

telemetry {
  db   = "consul"
  keys = ["123", "456"]
}
