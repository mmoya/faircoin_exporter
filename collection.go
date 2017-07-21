package main

import (
	"log"

	faircoin "github.com/mmoya/faircoin_rpcclient"
)

func updateActiveCVNs(c *faircoin.Client) {
	activeCvns, err := c.GetActiveCVNs()
	if err != nil {
		log.Fatal("GetActiveCvns: ", err)
	}

	for _, cvn := range activeCvns.Cvns {
		lastBlocksSigned.WithLabelValues(cvn.NodeID).Set(float64(cvn.LastBlocksSigned))
	}

	currentHeight.Set(float64(activeCvns.CurrentHeight))
}
