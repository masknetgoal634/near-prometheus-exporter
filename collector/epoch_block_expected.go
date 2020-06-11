package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type EpochBlockExpected struct {
	accountId string
	client    *nearapi.Client
	desc      *prometheus.Desc
}

func NewEpochBlockExpected(client *nearapi.Client, accountId string) *EpochBlockExpected {
	return &EpochBlockExpected{
		accountId: accountId,
		client:    client,
		desc: prometheus.NewDesc(
			"near_epoch_block_expected_number",
			"the number of block expected in epoch",
			nil,
			nil,
		),
	}
}

func (collector *EpochBlockExpected) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EpochBlockExpected) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.ExpectedBlocks(collector.accountId)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	value := float64(r)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
