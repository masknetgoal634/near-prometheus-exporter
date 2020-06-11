package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type CurrentStake struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewCurrentStake(client *nearapi.Client) *CurrentStake {
	return &CurrentStake{
		client: client,
		desc: prometheus.NewDesc(
			"near_current_stake",
			"current stake of a given account id",
			nil,
			nil,
		),
	}
}

func (collector *CurrentStake) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *CurrentStake) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.CurrentStake()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	value := float64(r)

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
