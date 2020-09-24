package collector

import (
	"fmt"
	"strconv"

	nearapi "github.com/masknetgoal634/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type NodeRpcMetrics struct {
	accountId                 string
	client                    *nearapi.Client
	epochBlockBroducedDesc    *prometheus.Desc
	epochBlockExpectedDesc    *prometheus.Desc
	seatPriceDesc             *prometheus.Desc
	currentStakeDesc          *prometheus.Desc
	epochStartHeightDesc      *prometheus.Desc
	blockNumberDesc           *prometheus.Desc
	syncingDesc               *prometheus.Desc
	versionBuildDesc          *prometheus.Desc
	currentValidatorStakeDesc *prometheus.Desc
	nextValidatorStakeDesc    *prometheus.Desc
	prevEpochKickoutDesc      *prometheus.Desc
	currentProposalsDesc      *prometheus.Desc
}

func NewNodeRpcMetrics(client *nearapi.Client, accountId string) *NodeRpcMetrics {
	return &NodeRpcMetrics{
		accountId: accountId,
		client:    client,
		epochBlockBroducedDesc: prometheus.NewDesc(
			"near_epoch_block_produced_number",
			"The number of block produced in epoch",
			nil,
			nil,
		),
		epochBlockExpectedDesc: prometheus.NewDesc(
			"near_epoch_block_expected_number",
			"The number of block expected in epoch",
			nil,
			nil,
		),
		seatPriceDesc: prometheus.NewDesc(
			"near_seat_price",
			"Validator seat price",
			nil,
			nil,
		),
		currentStakeDesc: prometheus.NewDesc(
			"near_current_stake",
			"Current stake of a given account id",
			nil,
			nil,
		),
		epochStartHeightDesc: prometheus.NewDesc(
			"near_epoch_start_height",
			"Near epoch start height",
			nil,
			nil,
		),
		blockNumberDesc: prometheus.NewDesc(
			"near_block_number",
			"The number of most recent block",
			nil,
			nil,
		),
		syncingDesc: prometheus.NewDesc(
			"near_sync_state",
			"Sync state",
			nil,
			nil,
		),
		versionBuildDesc: prometheus.NewDesc(
			"near_version_build",
			"The Near node version build",
			[]string{"version", "build"},
			nil,
		),
		currentValidatorStakeDesc: prometheus.NewDesc(
			"near_current_validator_stake",
			"Current amount of validator stake",
			[]string{"account_id", "public_key", "slashed", "shards", "num_produced_blocks", "num_expected_blocks"},
			nil,
		),
		nextValidatorStakeDesc: prometheus.NewDesc(
			"near_next_validator_stake",
			"The next validators",
			[]string{"account_id", "public_key", "shards"},
			nil,
		),
		currentProposalsDesc: prometheus.NewDesc(
			"near_current_proposals_stake",
			"Current proposals",
			[]string{"account_id", "public_key"},
			nil,
		),
		prevEpochKickoutDesc: prometheus.NewDesc(
			"near_prev_epoch_kickout",
			"Near previous epoch kicked out validators",
			[]string{"account_id", "reason"},
			nil,
		),
	}
}

func (collector *NodeRpcMetrics) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.epochBlockBroducedDesc
	ch <- collector.epochBlockExpectedDesc
	ch <- collector.seatPriceDesc
	ch <- collector.currentStakeDesc
	ch <- collector.epochStartHeightDesc
	ch <- collector.blockNumberDesc
	ch <- collector.syncingDesc
	ch <- collector.versionBuildDesc
	ch <- collector.currentValidatorStakeDesc
	ch <- collector.nextValidatorStakeDesc
	ch <- collector.currentProposalsDesc
	ch <- collector.prevEpochKickoutDesc
}

func (collector *NodeRpcMetrics) Collect(ch chan<- prometheus.Metric) {
	sr, err := collector.client.Get("status", nil)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.versionBuildDesc, err)
		return
	}
	syn := sr.Status.SyncInfo.Syncing
	var isSyncing int
	if syn {
		isSyncing = 1
	} else {
		isSyncing = 0
	}
	ch <- prometheus.MustNewConstMetric(collector.syncingDesc, prometheus.GaugeValue, float64(isSyncing))

	blockHeight := sr.Status.SyncInfo.LatestBlockHeight
	ch <- prometheus.MustNewConstMetric(collector.blockNumberDesc, prometheus.GaugeValue, float64(blockHeight))

	versionBuildInt := HashString(sr.Status.Version.Build)
	ch <- prometheus.MustNewConstMetric(collector.versionBuildDesc, prometheus.GaugeValue, float64(versionBuildInt), sr.Status.Version.Version, sr.Status.Version.Build)

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
		ch <- prometheus.NewInvalidMetric(collector.currentValidatorStakeDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.nextValidatorStakeDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.currentProposalsDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.prevEpochKickoutDesc, err)
		return
	}

	ch <- prometheus.MustNewConstMetric(collector.epochStartHeightDesc, prometheus.GaugeValue, float64(r.Validators.EpochStartHeight))

	var pb, eb, seatPrice, currentStake float64
	for _, v := range r.Validators.CurrentValidators {

		ch <- prometheus.MustNewConstMetric(collector.currentValidatorStakeDesc, prometheus.GaugeValue,
			float64(GetStakeFromString(v.Stake)), v.AccountId, v.PublicKey, strconv.FormatBool(v.IsSlashed), strconv.Itoa(len(v.Shards)), strconv.Itoa(int(v.NumProducedBlocks)), strconv.Itoa(int(v.NumExpectedBlocks)))

		t := GetStakeFromString(v.Stake)
		if seatPrice == 0 {
			seatPrice = t
		}
		if seatPrice > t {
			seatPrice = t
		}
		if v.AccountId == collector.accountId {
			pb = float64(v.NumProducedBlocks)
			eb = float64(v.NumExpectedBlocks)
			currentStake = t
		}
	}
	ch <- prometheus.MustNewConstMetric(collector.epochBlockBroducedDesc, prometheus.GaugeValue, pb)
	ch <- prometheus.MustNewConstMetric(collector.epochBlockExpectedDesc, prometheus.GaugeValue, eb)
	ch <- prometheus.MustNewConstMetric(collector.seatPriceDesc, prometheus.GaugeValue, seatPrice)
	ch <- prometheus.MustNewConstMetric(collector.currentStakeDesc, prometheus.GaugeValue, currentStake)

	for _, v := range r.Validators.NextValidators {
		ch <- prometheus.MustNewConstMetric(collector.nextValidatorStakeDesc, prometheus.GaugeValue,
			float64(GetStakeFromString(v.Stake)), v.AccountId, v.PublicKey, strconv.Itoa(len(v.Shards)))
	}

	for _, v := range r.Validators.CurrentProposals {
		ch <- prometheus.MustNewConstMetric(collector.currentProposalsDesc, prometheus.GaugeValue,
			float64(GetStakeFromString(v.Stake)), v.AccountId, v.PublicKey)
	}

	for _, v := range r.Validators.PrevEpochKickOut {
		ch <- prometheus.MustNewConstMetric(collector.prevEpochKickoutDesc, prometheus.GaugeValue, 0, v.AccountId, fmt.Sprintf("%v", v.Reason))
	}
}
