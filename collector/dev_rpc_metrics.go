package collector

import (
	"fmt"

	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type DevNodeRpcMetrics struct {
	client           *nearapi.Client
	versionBuildDesc *prometheus.Desc
}

func NewDevNodeRpcMetrics(client *nearapi.Client) *DevNodeRpcMetrics {
	return &DevNodeRpcMetrics{
		client: client,
		versionBuildDesc: prometheus.NewDesc(
			"near_dev_version_build",
			"The Dev Near node version build",
			[]string{"version", "build"},
			nil,
		),
	}
}

func (collector *DevNodeRpcMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.versionBuildDesc
}

func (collector *DevNodeRpcMetrics) Collect(ch chan<- prometheus.Metric) {
	sr, err := collector.client.Get("status", nil)
	if err != nil {
		fmt.Println("failed dev node get")
		ch <- prometheus.NewInvalidMetric(collector.versionBuildDesc, err)
		return
	}

	versionBuildInt := HashString(sr.Status.Version.Build)
	ch <- prometheus.MustNewConstMetric(collector.versionBuildDesc, prometheus.GaugeValue, float64(versionBuildInt), sr.Status.Version.Version, sr.Status.Version.Build)
}
