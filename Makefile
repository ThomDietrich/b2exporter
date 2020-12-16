.PHONY: build
build: b2exporter

.PHONY: clean
clean:
	go clean
	rm -f b2exporter

b2exporter: b2exporter.go
	go build -o b2exporter b2exporter.go

.PHONY: test
test:
	go test

.PHONY: run
run:
	go run b2exporter.go

image:
	docker build -t b2exporter:local .
