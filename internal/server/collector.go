package server

import (
	"github.com/SkYNewZ/speedtest-prometheus-exporter/internal/speedtest"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = (*collector)(nil)

type collector struct {
	data                    *speedtest.Result // carry our latest data
	pingJitterMetric        *prometheus.Desc
	pingLatencyMetric       *prometheus.Desc
	downloadBandwidthMetric *prometheus.Desc
	downloadBytesMetric     *prometheus.Desc
	downloadElapsedMetric   *prometheus.Desc
	uploadBandwidthMetric   *prometheus.Desc
	uploadBytesMetric       *prometheus.Desc
	uploadElapsedMetric     *prometheus.Desc
	packetLossMetric        *prometheus.Desc
}

func newCollector() *collector {
	return &collector{
		data: new(speedtest.Result),

		// Ping
		pingJitterMetric:  prometheus.NewDesc("speedtest_ping_jitter", "Ping Jitter in milliseconds", nil, nil),
		pingLatencyMetric: prometheus.NewDesc("speedtest_ping_latency", "Ping latency in milliseconds", nil, nil),

		// Download
		downloadBandwidthMetric: prometheus.NewDesc("speedtest_download_bandwidth", "Download bandwidth in megabits", nil, nil),
		downloadBytesMetric:     prometheus.NewDesc("speedtest_download_bytes", "Download bytes", nil, nil),
		downloadElapsedMetric:   prometheus.NewDesc("speedtest_download_elapsed", "Elapsed time in seconds to perform this test", nil, nil),

		// Upload
		uploadBandwidthMetric: prometheus.NewDesc("speedtest_upload_bandwidth", "Upload bandwidth in megabits", nil, nil),
		uploadBytesMetric:     prometheus.NewDesc("speedtest_upload_bytes", "Upload bytes", nil, nil),
		uploadElapsedMetric:   prometheus.NewDesc("speedtest_upload_elapsed", "Elapsed time in seconds to perform this test", nil, nil),

		// PacketLoss
		packetLossMetric: prometheus.NewDesc("speedtest_packet_loss", "Packets loss during the test", nil, nil),
	}
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.pingJitterMetric
	ch <- c.pingLatencyMetric
	ch <- c.downloadBandwidthMetric
	ch <- c.downloadBytesMetric
	ch <- c.downloadElapsedMetric
	ch <- c.uploadBandwidthMetric
	ch <- c.uploadBytesMetric
	ch <- c.uploadElapsedMetric
	ch <- c.packetLossMetric
}

//Collect implements required collect function for all promehteus collectors
func (c *collector) Collect(ch chan<- prometheus.Metric) {
	// Ping
	ch <- prometheus.MustNewConstMetric(c.pingJitterMetric, prometheus.GaugeValue, c.data.Ping.Jitter)
	ch <- prometheus.MustNewConstMetric(c.pingLatencyMetric, prometheus.GaugeValue, c.data.Ping.Latency)

	// Download
	ch <- prometheus.MustNewConstMetric(c.downloadBandwidthMetric, prometheus.GaugeValue, float64(c.data.Download.Bandwidth))
	ch <- prometheus.MustNewConstMetric(c.downloadBytesMetric, prometheus.GaugeValue, float64(c.data.Download.Bytes))
	ch <- prometheus.MustNewConstMetric(c.downloadElapsedMetric, prometheus.GaugeValue, float64(c.data.Download.Elapsed))

	// Upload
	ch <- prometheus.MustNewConstMetric(c.uploadBandwidthMetric, prometheus.GaugeValue, float64(c.data.Upload.Bandwidth))
	ch <- prometheus.MustNewConstMetric(c.uploadBytesMetric, prometheus.GaugeValue, float64(c.data.Upload.Bytes))
	ch <- prometheus.MustNewConstMetric(c.uploadElapsedMetric, prometheus.GaugeValue, float64(c.data.Upload.Elapsed))

	// Packet loss
	ch <- prometheus.MustNewConstMetric(c.packetLossMetric, prometheus.GaugeValue, c.data.PacketLoss)
}
