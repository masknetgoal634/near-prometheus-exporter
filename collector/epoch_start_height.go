package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type EpochStartHeight struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewEpochStartHeight(client *nearapi.Client) *EpochStartHeight {
	return &EpochStartHeight{
		client: client,
		desc: prometheus.NewDesc(
			"near_epoch_start_height",
			"near epoch start height",
			nil,
			nil,
		),
	}
}

func (collector *EpochStartHeight) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EpochStartHeight) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.EpochStartHeight()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	value := float64(r)

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
