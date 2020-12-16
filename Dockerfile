FROM golang:1.15-buster AS builder

WORKDIR /go/src/b2exporter
COPY b2exporter.go .
RUN go get -d -v ./...
#RUN go get -u ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -o b2exporter ./...

#

FROM alpine:latest
RUN apk --no-cache add ca-certificates

COPY --from=builder /go/src/b2exporter/b2exporter /

EXPOSE 8080
ENTRYPOINT ["/b2exporter"]
CMD [ "-period", "1h30m" ]
