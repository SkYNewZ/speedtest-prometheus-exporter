FROM debian:buster-slim as speedtest-builder

RUN apt-get update && apt-get install -y curl jq gnupg1 apt-transport-https dirmngr && \
    apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 379CE192D401AB61 && \
    echo "deb https://ookla.bintray.com/debian buster main" | tee /etc/apt/sources.list.d/speedtest.list && \
    apt-get update && apt-get install -y speedtest


FROM golang:1.15.2-alpine3.12 as go-builder
WORKDIR /go/src/github.com/SkYNewZ/speedtest-prometheus-exporter

COPY go.* ./
RUN go mod download

COPY . .
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    go build -a -installsuffix cgo -o /speedtest-prometheus-exporter


FROM scratch

# Get speedtest-cli
COPY --from=speedtest-builder /usr/bin/speedtest /usr/bin/speedtest
COPY --from=speedtest-builder /etc/ssl/certs /etc/ssl/certs

# Get my app
COPY --from=go-builder /speedtest-prometheus-exporter /speedtest-prometheus-exporter

ENTRYPOINT [ "/speedtest-prometheus-exporter" ]
CMD [ "-speedtest-path", "/usr/bin/speedtest" ]