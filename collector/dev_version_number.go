package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type DevVersionNumber struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewDevVersionNumber(client *nearapi.Client) *DevVersionNumber {
	return &DevVersionNumber{
		client: client,
		desc: prometheus.NewDesc(
			"near_dev_version_number",
			"dev near node version number",
			nil,
			nil,
		),
	}
}

func (collector *DevVersionNumber) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *DevVersionNumber) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.Get("status", nil)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	v, err := getFloatVersionFromString(r.Status.Version.Version)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, v)
}
