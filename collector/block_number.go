package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type BlockNumber struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewBlockNumber(client *nearapi.Client) *BlockNumber {
	return &BlockNumber{
		client: client,
		desc: prometheus.NewDesc(
			"near_block_number",
			"the number of most recent block",
			nil,
			nil,
		),
	}
}

func (collector *BlockNumber) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *BlockNumber) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.LatestBlockHeight()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	value := float64(r)

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
