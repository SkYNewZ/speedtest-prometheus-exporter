package main

import "github.com/prometheus/client_golang/prometheus"

type speedTestCollector struct {
	pingJitterMetric  *prometheus.Desc
	pingLatencyMetric *prometheus.Desc

	downloadBandwithMetric *prometheus.Desc
	downloadBytesMetric    *prometheus.Desc
	downloadElapsedMetric  *prometheus.Desc

	uploadBandwithMetric *prometheus.Desc
	uploadBytesMetric    *prometheus.Desc
	uploadElapsedMetric  *prometheus.Desc

	packetLossMetric *prometheus.Desc
}

func newspeedTestCollector() *speedTestCollector {
	return &speedTestCollector{
		// Ping
		pingJitterMetric:  prometheus.NewDesc("speedtest_ping_jitter", "Ping Jitter in milliseconds", nil, nil),
		pingLatencyMetric: prometheus.NewDesc("speedtest_ping_latency", "Ping latency in milliseconds", nil, nil),

		// Download
		downloadBandwithMetric: prometheus.NewDesc("speedtest_download_bandwith", "Download bandwith in megabits", nil, nil),
		downloadBytesMetric:    prometheus.NewDesc("speedtest_download_bytes", "Download bytes", nil, nil),
		downloadElapsedMetric:  prometheus.NewDesc("speedtest_download_elapsed", "Elapsed time in seconds to perform this test", nil, nil),

		// Upload
		uploadBandwithMetric: prometheus.NewDesc("speedtest_upload_bandwith", "Upload bandwith in megabits", nil, nil),
		uploadBytesMetric:    prometheus.NewDesc("speedtest_upload_bytes", "Upload bytes", nil, nil),
		uploadElapsedMetric:  prometheus.NewDesc("speedtest_upload_elapsed", "Elapsed time in seconds to perform this test", nil, nil),

		// PacketLoss
		packetLossMetric: prometheus.NewDesc("speedtest_packet_loss", "Packets loss during the test", nil, nil),
	}
}

func (collector *speedTestCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.pingJitterMetric
	ch <- collector.pingLatencyMetric

	ch <- collector.downloadBandwithMetric
	ch <- collector.downloadBytesMetric
	ch <- collector.downloadElapsedMetric

	ch <- collector.uploadBandwithMetric
	ch <- collector.uploadBytesMetric
	ch <- collector.uploadElapsedMetric

	ch <- collector.packetLossMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *speedTestCollector) Collect(ch chan<- prometheus.Metric) {

	// Get speedtest
	speedtest := execute()

	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.

	// Ping
	ch <- prometheus.MustNewConstMetric(collector.pingJitterMetric, prometheus.GaugeValue, speedtest.Ping.Jitter)
	ch <- prometheus.MustNewConstMetric(collector.pingLatencyMetric, prometheus.GaugeValue, speedtest.Ping.Latency)

	// Download
	ch <- prometheus.MustNewConstMetric(collector.downloadBandwithMetric, prometheus.GaugeValue, float64(speedtest.Download.Bandwidth))
	ch <- prometheus.MustNewConstMetric(collector.downloadBytesMetric, prometheus.GaugeValue, float64(speedtest.Download.Bytes))
	ch <- prometheus.MustNewConstMetric(collector.downloadElapsedMetric, prometheus.GaugeValue, float64(speedtest.Download.Elapsed))

	// Upload
	ch <- prometheus.MustNewConstMetric(collector.uploadBandwithMetric, prometheus.GaugeValue, float64(speedtest.Upload.Bandwidth))
	ch <- prometheus.MustNewConstMetric(collector.uploadBytesMetric, prometheus.GaugeValue, float64(speedtest.Upload.Bytes))
	ch <- prometheus.MustNewConstMetric(collector.uploadElapsedMetric, prometheus.GaugeValue, float64(speedtest.Upload.Elapsed))

	// Packet loss
	ch <- prometheus.MustNewConstMetric(collector.packetLossMetric, prometheus.GaugeValue, speedtest.PacketLoss)
}
