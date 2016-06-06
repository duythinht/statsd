build:
	CGO_ENABLED=0 GOOS=linux go build -o dist/statsd -a -tags netgo -ldflags '-w' .
test:
	go build -o dist/statsd-osx
	./dist/statsd-osx -version
	./dist/statsd-osx -config test.hcl
