package collector

import (
	"fmt"
	"strconv"

	nearapi "github.com/bisontrails/near-exporter/client"
	"github.com/prometheus/client_golang/prometheus"
)

type NodeRpcMetrics struct {
	accountId                 string
	internalClient            *nearapi.Client
	externalClient            *nearapi.Client
	epochBlockBroducedDesc    *prometheus.Desc
	epochBlockExpectedDesc    *prometheus.Desc
	seatPriceDesc             *prometheus.Desc
	currentStakeDesc          *prometheus.Desc
	epochStartHeightDesc      *prometheus.Desc
	blockHeightInternalDesc   *prometheus.Desc
	blockLagDesc              *prometheus.Desc
	blocksMissedDesc          *prometheus.Desc
	syncingDesc               *prometheus.Desc
	versionBuildDesc          *prometheus.Desc
	currentValidatorStakeDesc *prometheus.Desc
	nextValidatorStakeDesc    *prometheus.Desc
	prevEpochKickoutDesc      *prometheus.Desc
	currentProposalsDesc      *prometheus.Desc
}

func NewNodeRpcMetrics(
	internalClient *nearapi.Client,
	externalClient *nearapi.Client,
	accountId string) *NodeRpcMetrics {

	return &NodeRpcMetrics{
		accountId:      accountId,
		internalClient: internalClient,
		externalClient: externalClient,
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
		blockHeightInternalDesc: prometheus.NewDesc(
			"near_internal_block_height",
			"The head of the NEAR chain",
			nil,
			nil,
		),
		blockLagDesc: prometheus.NewDesc(
			"near_block_lag",
			"The number of blocks behind rpc endpoint block head.",
			nil,
			nil,
		),
		blocksMissedDesc: prometheus.NewDesc(
			"near_blocks_missed",
			"The number of blocks missed while validating in the active set.",
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
			[]string{"account_id", "reason", "produced", "expected", "stake_u128", "threshold_u128"},
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
	ch <- collector.blockHeightInternalDesc
	ch <- collector.blockLagDesc
	ch <- collector.blocksMissedDesc
	ch <- collector.syncingDesc
	ch <- collector.versionBuildDesc
	ch <- collector.currentValidatorStakeDesc
	ch <- collector.nextValidatorStakeDesc
	ch <- collector.currentProposalsDesc
	ch <- collector.prevEpochKickoutDesc
}

func (collector *NodeRpcMetrics) Collect(ch chan<- prometheus.Metric) {
	sr, err := collector.internalClient.Get("status", nil)
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.versionBuildDesc, err)
		return
	}

	srExt, err := collector.externalClient.Get("status", nil)
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
	ch <- prometheus.MustNewConstMetric(collector.blockHeightInternalDesc, prometheus.GaugeValue, float64(blockHeight))

	fmt.Printf("External BlockHeght: %d", srExt.Status.SyncInfo.LatestBlockHeight)
	fmt.Printf("Internal BlockHeght: %d", sr.Status.SyncInfo.LatestBlockHeight)
	blockLag := srExt.Status.SyncInfo.LatestBlockHeight - sr.Status.SyncInfo.LatestBlockHeight
	ch <- prometheus.MustNewConstMetric(collector.blockLagDesc, prometheus.GaugeValue, float64(blockLag))

	versionBuildInt := HashString(sr.Status.Version.Build)
	ch <- prometheus.MustNewConstMetric(collector.versionBuildDesc, prometheus.GaugeValue, float64(versionBuildInt), sr.Status.Version.Version, sr.Status.Version.Build)

	r, err := collector.internalClient.Get("validators", []uint64{blockHeight})
	if err != nil {
		ch <- prometheus.NewInvalidMetric(collector.epochBlockBroducedDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.epochBlockExpectedDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.seatPriceDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.currentStakeDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.epochStartHeightDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.blockHeightInternalDesc, err)
		ch <- prometheus.NewInvalidMetric(collector.blocksMissedDesc, err)
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
	ch <- prometheus.MustNewConstMetric(collector.blocksMissedDesc, prometheus.GaugeValue, eb - pb)
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
		if reason, ok := v.Reason["NotEnoughStake"]; ok {
			if threshold, ok2 := reason["threshold_u128"]; ok2 {
				// set seat price if we have "threshold_u128"
				seatPrice = GetStakeFromString(threshold.(string))
			}
			if stake, ok2 := reason["stake_u128"]; ok2 {
				ch <- prometheus.MustNewConstMetric(collector.prevEpochKickoutDesc, prometheus.GaugeValue,
					GetStakeFromString(stake.(string)), v.AccountId, "NotEnoughStake", "", "", stake.(string), reason["threshold_u128"].(string))
			}

		} else if val, ok := v.Reason["NotEnoughBlocks"]; ok {
			if produced, ok2 := val["produced"]; ok2 {
				ch <- prometheus.MustNewConstMetric(collector.prevEpochKickoutDesc, prometheus.GaugeValue,
					float64(produced.(float64)), v.AccountId, "NotEnoughBlocks", fmt.Sprintf("%v", produced.(float64)), fmt.Sprintf("%v", val["expected"].(float64)), "", "")
			}
		}
	}
}
