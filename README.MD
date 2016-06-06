## Go statsd with pattern
Config: Using HCL config syntax, with key using regular expression

    influx_url = "http://localhost:8086"

    /*influx_url = "udp://192.168.99.100:8089"*/

    influx_username = "admin"

    influx_password = "password"

    listen = ":8127"

    default_db = "statsd"

    telemetry {
      db   = "test"
      keys = ["^test.(?P<bucket>hello).(?P<name>[^ ]*)$", "def"]
    }

    telemetry {
      db   = "metrics"
      keys = ["the.metric", "lmao.hello.abc"]
    }

RUN:

    $ statsd -config [path_to_config_file]
