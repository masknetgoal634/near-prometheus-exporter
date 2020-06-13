package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type VersionNumber struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewVersionNumber(client *nearapi.Client) *VersionNumber {
	return &VersionNumber{
		client: client,
		desc: prometheus.NewDesc(
			"near_version_number",
			"near node version number",
			nil,
			nil,
		),
	}
}

func (collector *VersionNumber) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *VersionNumber) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.VersionNumber()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	v, err := getFloatVersionFromString(r)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, v)
}
