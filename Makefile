build:
	CGO_ENABLED=0 GOOS=linux go build -o dist/statsd -a -tags netgo -ldflags '-w' .
