FROM golang:latest
LABEL maintainer="jon@yakshaver.dev"

ARG UID=1000
ARG GID=1000
RUN groupadd -g ${GID} -r app && useradd -u ${UID} -r -g app app

WORKDIR /go/src/b2exporter
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

USER app
ENTRYPOINT ["b2exporter"]
