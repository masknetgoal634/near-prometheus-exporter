package collector

import (
	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type SeatPrice struct {
	client *nearapi.Client
	desc   *prometheus.Desc
}

func NewSeatPrice(client *nearapi.Client) *SeatPrice {
	return &SeatPrice{
		client: client,
		desc: prometheus.NewDesc(
			"near_seat_price",
			"validator seat price",
			nil,
			nil,
		),
	}
}

func (collector *SeatPrice) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *SeatPrice) Collect(ch chan<- prometheus.Metric) {
	r, err := collector.client.SeatPrice()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}
	value := float64(r)

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, value)
}
