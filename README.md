# speedtest-prometheus-exporter

Do you want to monitor your home bandwidth and expose these metrics to Prometheus ? This projet makes it for you. It
starts a Go web server, exposes a `/metrics` endpoint where you can see some `speedtest_*` metrics !

## Why ?

I found some posts about bandwidth monitoring with [InfluxDB](https://www.influxdata.com)
and [Grafana](https://grafana.com).

- (The original post I found, in French) https://blog.eleven-labs.com/fr/monitorer-son-debit-internet/
- (A similar one in
  english) https://gonzalo123.com/2018/11/26/monitoring-the-bandwidth-with-grafana-influxdb-and-docker/

Just for fun (and to try Grafana), I set up the stack and started to have some beautiful dashboards.

But I already have a Prometheus on my own, and I didn't want to multiply time-series databases. So I built this,
a `speedtest-prometheus-exporter`

## How does it work ?

This project is designed has two parts:

- A server will receive the speedtest results at `/results` and expose them on `/metrics`
- A Kubernetes cronjob which will run periodically run speed tests and report results to `http://server:port/results`

By default, the job will run each 10 minutes (`*/10 * * * *`).

## Getting started

1. You can build your container or just mine `skynewz/speedtest-prometheus-exporter`.
2. Anyway, run it with `docker run -it --rm -e PORT=8080 -p 8080:8080 skynewz/speedtest-prometheus-exporter`
3. Set your Prometheus instance to add the brand-new target (see [`prometheus.yml`](prometheus-sample.yaml))
4. You can use the [dashboard example](Speedtest-1637677700797.json) to view results.

## Usage

```
Usage:
  speedtest-prometheus-exporter [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  server      Receive speedtest results and expose them as Prometheus metrics
  speedtest   Run speedtest and report to given server

Flags:
  -h, --help      help for speedtest-prometheus-exporter
      --verbose   Verbose output

Use "speedtest-prometheus-exporter [command] --help" for more information about a command.

```