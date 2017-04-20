package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	cvnStatsBlocks = 10
)

var (
	lastUpdate = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "last_update_unixtime",
		Subsystem: "faircoin2",
		Help:      "Last time faircoind was polled",
	})

	currentHeight = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "cvn_current_height",
		Subsystem: "faircoin2",
		Help:      "Height of the block chain",
	})

	lastBlocksSigned = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name:      "cvn_last_blocks_signed",
			Subsystem: "faircoin2",
			Help:      fmt.Sprintf("Signed blocks (of the last %d)", cvnStatsBlocks),
		},
		[]string{"node_id"},
	)

	cvnStatsBlocksMetric = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "cvn_stats_blocks",
		Subsystem: "faircoin2",
		Help:      "How many blocks are accounted to get CVN stats",
	})
)
