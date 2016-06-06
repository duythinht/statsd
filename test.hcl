influx_url = "http://localhost:8086"

/*influx_url = "udp://192.168.99.100:8089"*/

influx_username = "admin"

influx_password = "password"

listen = ":8127"

default_db = "statsd"

telemetry {
  db   = "response"
  keys = ["^response.(?P<service>[^.]+).(?P<method>[A-Z]+).(?P<url>[^.]+).(?P<status>\\d+)"]
}
