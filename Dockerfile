FROM --platform=$TARGETPLATFORM debian:bullseye as speedtest

RUN apt-get update && \
  apt-get install -y curl && \
  curl -s https://install.speedtest.net/app/cli/install.deb.sh | bash && \
  apt-get install -y speedtest

FROM golang:1.17.3-alpine3.14 as exporter
ARG TARGETARCH
ARG TARGETVARIANT

WORKDIR /app

ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GOARCH=$TARGETARCH

COPY . .
RUN go mod vendor
RUN if [ $TARGETARCH == "arm" ]; then export GOARM=$(echo $TARGETVARIANT | tr -d "v"); fi && \
    echo "GOARCH: $GOARCH, GOOS: $GOOS, GOARM: $GOARM" && \
    go build -mod vendor -ldflags "-s -w" -a -installsuffix cgo -o /speedtest-prometheus-exporter .

FROM --platform=$TARGETPLATFORM gcr.io/distroless/static:nonroot

COPY --from=speedtest /usr/bin/speedtest /usr/bin/speedtest
COPY --from=exporter /speedtest-prometheus-exporter /speedtest-prometheus-exporter

ENTRYPOINT [ "/speedtest-prometheus-exporter" ]
CMD [ "--help" ]
