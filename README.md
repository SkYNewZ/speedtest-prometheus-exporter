# speedtest-prometheus-exporter

Do you want to monitor your home bandwith and expose these metrics to Prometheus ? This projet makes it for you.
It starts a Go web server, exposes a `/metrics` endpoint where you can see some `speedtest_*` metrics !

## Why ?

I found some posts about bandwith monitoring with [InfluxDB](https://www.influxdata.com) and [Grafana](https://grafana.com).

- (The original post I found, in french) https://blog.eleven-labs.com/fr/monitorer-son-debit-internet/
- (A similar one in english) https://gonzalo123.com/2018/11/26/monitoring-the-bandwidth-with-grafana-influxdb-and-docker/

Just for fun (and to try Grafana), I set up the stack and started to have some beautiful dashboards.

But I already have a Prometheus on my own and I didn't want to multiply time-series databases. So I built this, a `speedtest-prometheus-exporter`

## How does it work ?

Because a speedtest can be very long, when you start the container (or the Go project - ensure you have the [`speedtest` cli](https://www.speedtest.net/apps/cli)), it immediately start a first speed test and store it on local memory.
Every minute, the process will restart.

- Execute a new speed test if last result is _older than 1 minute_
- Otherwise, read the previous result from cache

In the meantime, the web server is always running and Prometheus can continue to scrape your speed test metrics.

## Getting started

1. You can build your container or just mine `skynewz/speedtest-prometheus-exporter`.
2. Anyway, run it with `docker run -it --rm -p 8080:8080 skynewz/speedtest-prometheus-exporter`
3. Set your Prometheus instance to add the brand new target (see [`prometheus.yml`](/prometheus.yml#32))

Be care of the `scrape_interval` and `scrape_timeout`. Because a speed test can be long, Prometheus timeout is important.
Indeed, if your `scrape_timeout` is too low, and the cache is empty (or expired), the request can take more than 25 secondes.
For `scrape_interval`, set what you want, but this value must be greater or equal than `scrape_timeout`.

## Usage

```
Usage of speedtest-prometheus-exporter:
  -speedtest-path string
    	the speedtest cli path (default "/usr/bin/speedtest")
```
