package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type DevVersionBuild struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewDevVersionBuild(client *nearapi.Client) *DevVersionBuild {
	return &DevVersionBuild{
		client: client,
		desc: prometheus.NewDesc(
			"near_dev_version_build",
			"dev near node version build",
			nil,
			nil,
		),
	}
}

func (collector *DevVersionBuild) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *DevVersionBuild) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.Get("status", nil)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	v, err := getFloat64FromString(r.Status.Version.Build)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, v)
}
