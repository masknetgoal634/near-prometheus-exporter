package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type VersionBuild struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewVersionBuild(client *nearapi.Client) *VersionBuild {
	return &VersionBuild{
		client: client,
		desc: prometheus.NewDesc(
			"near_version_build",
			"near node version build",
			nil,
			nil,
		),
	}
}

func (collector *VersionBuild) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *VersionBuild) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.VersionBuild()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	v, err := getFloat64FromString(r)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, v)
}
