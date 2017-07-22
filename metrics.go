package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	cvnStatsBlocks = 10
)

var (
	lastUpdate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "last_update_unixtime",
		Subsystem: "faircoin",
		Help:      "Last time faircoind was polled",
	})

	currentHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "current_height",
		Subsystem: "faircoin",
		Help:      "Height of the block chain",
	})

	lastBlockSigned = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "cvn_last_block_signed",
			Subsystem: "faircoin",
			Help:      "Is the last block signed by this CVN",
		},
		[]string{"node_id"})
)

func registerMetrics() {
	prometheus.MustRegister(lastUpdate)
	prometheus.MustRegister(currentHeight)
	prometheus.MustRegister(lastBlockSigned)
}
