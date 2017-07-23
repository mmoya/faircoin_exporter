package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cvnStatsBlocks = 10
)

var (
	// blockchain related
	mCurrentHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "current_height",
		Subsystem: "faircoin",
		Help:      "Height of the block chain",
	})

	mLastBlockTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "last_block_timestamp",
		Subsystem: "faircoin",
		Help:      "Timestamp of the last block",
	})

	mLastBlockHeardTimestamp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "last_block_heard_timestamp",
		Subsystem: "faircoin",
		Help:      "Timestamp the last block was heard",
	})

	// cvn related
	mCvnCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "cvn_count",
		Subsystem: "faircoin",
		Help:      "How many CVNs are registered",
	})

	mLastBlockSigned = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "cvn_last_block_signed",
			Subsystem: "faircoin",
			Help:      "Was the last block signed by a CVN",
		},
		[]string{"node_id"})
)

func registerMetrics() {
	prometheus.MustRegister(mCurrentHeight)
	prometheus.MustRegister(mLastBlockTimestamp)
	prometheus.MustRegister(mLastBlockHeardTimestamp)

	prometheus.MustRegister(mCvnCount)
	prometheus.MustRegister(mLastBlockSigned)
}
