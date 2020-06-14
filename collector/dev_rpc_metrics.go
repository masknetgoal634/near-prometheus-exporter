package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type DevNodeRpcMetrics struct {
	client            *nearapi.Client
	versionNumberDesc *prometheus.Desc
	versionBuildDesc  *prometheus.Desc
}

func NewDevNodeRpcMetrics(client *nearapi.Client) *DevNodeRpcMetrics {
	return &DevNodeRpcMetrics{
		client: client,
		versionNumberDesc: prometheus.NewDesc(
			"near_dev_version_number",
			"dev near node version number",
			nil,
			nil,
		),
		versionBuildDesc: prometheus.NewDesc(
			"near_dev_version_build",
			"dev near node version build",
			nil,
			nil,
		),
	}
}

func (collector *DevNodeRpcMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.versionBuildDesc
	ch <- collector.versionNumberDesc
}

func (collector *DevNodeRpcMetrics) Collect(ch chan<- prometheus.Metric) {
	sr, err := collector.client.Get("status", nil)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.versionNumberDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.versionBuildDesc, err)
		return
	}

	versionNumber := sr.Status.Version.Version
	vn, err := GetFloatVersionFromString(versionNumber)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.versionNumberDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.versionBuildDesc, err)
		return
	}

	versionBuildStr := sr.Status.Version.Build
	versionBuildInt := HashString(versionBuildStr)

	value := float64(vn)
	ch <- prometheus.MustNewConstMetric(collector.versionNumberDesc, prometheus.GaugeValue, value)

	value = float64(versionBuildInt)
	ch <- prometheus.MustNewConstMetric(collector.versionBuildDesc, prometheus.GaugeValue, value)
}
