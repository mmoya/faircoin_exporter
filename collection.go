package main

import (
	"log"

	client "gitlab.com/mmoya/faircoin2_rpcclient"
)

func updateActiveCVNs(c *client.FC2Client) {
	activeCvns, err := c.GetActiveCVNs()
	if err != nil {
		log.Fatal("GetActiveCvns: ", err)
	}

	for _, cvn := range activeCvns.Cvns {
		lastBlocksSigned.WithLabelValues(cvn.NodeID).Set(float64(cvn.LastBlocksSigned))
	}

	currentHeight.Set(float64(activeCvns.CurrentHeight))
}
