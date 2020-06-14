package collector

import (
	"fmt"

	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type NodeRpcMetrics struct {
	accountId              string
	client                 *nearapi.Client
	epochBlockBroducedDesc *prometheus.Desc
	epochBlockExpectedDesc *prometheus.Desc
	seatPriceDesc          *prometheus.Desc
	currentStakeDesc       *prometheus.Desc
	versionNumberDesc      *prometheus.Desc
	epochStartHeightDesc   *prometheus.Desc
	blockNumberDesc        *prometheus.Desc
	syncingDesc            *prometheus.Desc
	versionBuildDesc       *prometheus.Desc
}

func NewNodeRpcMetrics(client *nearapi.Client, accountId string) *NodeRpcMetrics {
	return &NodeRpcMetrics{
		accountId: accountId,
		client:    client,
		epochBlockBroducedDesc: prometheus.NewDesc(
			"near_epoch_block_produced_number",
			"the number of block produced in epoch",
			nil,
			nil,
		),
		epochBlockExpectedDesc: prometheus.NewDesc(
			"near_epoch_block_expected_number",
			"the number of block expected in epoch",
			nil,
			nil,
		),
		seatPriceDesc: prometheus.NewDesc(
			"near_seat_price",
			"validator seat price",
			nil,
			nil,
		),
		currentStakeDesc: prometheus.NewDesc(
			"near_current_stake",
			"current stake of a given account id",
			nil,
			nil,
		),
		versionNumberDesc: prometheus.NewDesc(
			"near_version_number",
			"near node version number",
			nil,
			nil,
		),
		epochStartHeightDesc: prometheus.NewDesc(
			"near_epoch_start_height",
			"near epoch start height",
			nil,
			nil,
		),
		blockNumberDesc: prometheus.NewDesc(
			"near_block_number",
			"the number of most recent block",
			nil,
			nil,
		),
		syncingDesc: prometheus.NewDesc(
			"near_sync_state",
			"sync state",
			nil,
			nil,
		),
		versionBuildDesc: prometheus.NewDesc(
			"near_version_build",
			"near node version build",
			nil,
			nil,
		),
	}
}

func (collector *NodeRpcMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.epochBlockBroducedDesc
	ch <- collector.epochBlockExpectedDesc
	ch <- collector.seatPriceDesc
	ch <- collector.currentStakeDesc
	ch <- collector.versionNumberDesc
	ch <- collector.epochStartHeightDesc
	ch <- collector.blockNumberDesc
	ch <- collector.syncingDesc
	ch <- collector.versionBuildDesc
}

func (collector *NodeRpcMetrics) Collect(ch chan<- prometheus.Metric) {
	sr, err := collector.client.Get("status", nil)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.versionNumberDesc, err)
		return
	}
	syn := sr.Status.SyncInfo.Syncing
	var isSyncing int
	if syn {
		isSyncing = 1
	} else {
		isSyncing = 0
	}

	blockHeight := sr.Status.SyncInfo.LatestBlockHeight
	versionNumber := sr.Status.Version.Version
	vn, err := GetFloatVersionFromString(versionNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	versionBuildStr := sr.Status.Version.Build
	versionBuildInt := HashString(versionBuildStr)

	r, err := collector.client.Get("validators", []uint64{blockHeight})
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.epochBlockBroducedDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.epochBlockExpectedDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.seatPriceDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.currentStakeDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.epochStartHeightDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.blockNumberDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.syncingDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.versionBuildDesc, err)
		return
	}

	epochStartHeight := r.Validators.EpochStartHeight

	var pb, eb, seatPrice, currentStake int64
	for _, v := range r.Validators.CurrentValidators {
		t := StringToInt64(v.Stake)
		if seatPrice == 0 {
			seatPrice = t
		}
		if seatPrice > t {
			seatPrice = t
		}
		if v.AccountId == collector.accountId {
			pb = v.NumProducedBlocks
			eb = v.NumExpectedBlocks
			currentStake = t
		}
	}

	value := float64(pb)
	ch <- prometheus.MustNewConstMetric(collector.epochBlockBroducedDesc, prometheus.GaugeValue, value)

	value = float64(eb)
	ch <- prometheus.MustNewConstMetric(collector.epochBlockExpectedDesc, prometheus.GaugeValue, value)

	value = float64(seatPrice)
	ch <- prometheus.MustNewConstMetric(collector.seatPriceDesc, prometheus.GaugeValue, value)

	value = float64(currentStake)
	ch <- prometheus.MustNewConstMetric(collector.currentStakeDesc, prometheus.GaugeValue, value)

	value = float64(vn)
	ch <- prometheus.MustNewConstMetric(collector.versionNumberDesc, prometheus.GaugeValue, value)

	value = float64(epochStartHeight)
	ch <- prometheus.MustNewConstMetric(collector.epochStartHeightDesc, prometheus.GaugeValue, value)

	value = float64(blockHeight)
	ch <- prometheus.MustNewConstMetric(collector.blockNumberDesc, prometheus.GaugeValue, value)

	value = float64(isSyncing)
	ch <- prometheus.MustNewConstMetric(collector.syncingDesc, prometheus.GaugeValue, value)

	value = float64(versionBuildInt)
	ch <- prometheus.MustNewConstMetric(collector.versionBuildDesc, prometheus.GaugeValue, value)
}
