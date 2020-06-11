package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type EpochBlockProduced struct {
	accountId string
	client    *nearapi.Client
	desc      *prometheus.Desc
}

func NewEpochBlockProduced(client *nearapi.Client, accountId string) *EpochBlockProduced {
	return &EpochBlockProduced{
		accountId: accountId,
		client:    client,
		desc: prometheus.NewDesc(
			"near_epoch_block_produced_number",
			"the number of block produced in epoch",
			nil,
			nil,
		),
	}
}

func (collector *EpochBlockProduced) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EpochBlockProduced) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.ProducedBlocks(collector.accountId)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	value := float64(r)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
