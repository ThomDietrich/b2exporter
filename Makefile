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
	docker build -t shavenyak/b2exporter:dev --build-arg APP_VERSION=$(APP_VERSION) .
