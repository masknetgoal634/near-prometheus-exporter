package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type NearSyncing struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewNearSyncing(client *nearapi.Client) *NearSyncing {
	return &NearSyncing{
		client: client,
		desc: prometheus.NewDesc(
			"near_sync_state",
			"sync state",
			nil,
			nil,
		),
	}
}

func (collector *NearSyncing) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *NearSyncing) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.Syncing()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	var result int
	if r {
		// syncing
		result = 1
	} else {
		// not syncing,  its good
		result = 0
	}
	value := float64(result)
	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
