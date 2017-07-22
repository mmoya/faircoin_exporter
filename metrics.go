package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cvnStatsBlocks = 10
)

var (
	mLastUpdate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "last_update_unixtime",
		Subsystem: "faircoin",
		Help:      "Last time faircoind was polled",
	})

	mCurrentHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "current_height",
		Subsystem: "faircoin",
		Help:      "Height of the block chain",
	})

	mLastBlockSigned = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "cvn_last_block_signed",
			Subsystem: "faircoin",
			Help:      "Was the last block signed by a CVN",
		},
		[]string{"node_id"})

	mCvnCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "cvn_count",
		Subsystem: "faircoin",
		Help:      "How many CVNs are registered",
	})
)

func registerMetrics() {
	prometheus.MustRegister(mLastUpdate)
	prometheus.MustRegister(mCurrentHeight)
	prometheus.MustRegister(mLastBlockSigned)
	prometheus.MustRegister(mCvnCount)
}
