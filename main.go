package main

import (
	"flag"
	"net/http"
	"time"

	"log"

	client "github.com/mmoya/faircoin_rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress = flag.String("listen-address", ":9132", "The address to listen on for HTTP requests")
	rpcURL        = flag.String("rpc.url", "http://127.0.0.1:40405", "URL of the RPC server")
)

func init() {
	prometheus.MustRegister(lastUpdate)
	prometheus.MustRegister(currentHeight)
	prometheus.MustRegister(lastBlocksSigned)
	prometheus.MustRegister(cvnStatsBlocksMetric)

	cvnStatsBlocksMetric.Set(cvnStatsBlocks)
}

func updateState(c *client.FC2Client) {
	for {
		updateActiveCVNs(c)
		lastUpdate.Set(float64(time.Now().Unix()))

		time.Sleep(5 * time.Second)
	}
}

func main() {
	flag.Parse()

	c := client.New(*rpcURL, client.CookieCredential())

	log.Printf("connecting to faircoind rpc in %s\n", *rpcURL)

	go updateState(c)

	log.Printf("listening in %s\n", *listenAddress)

	http.Handle("/metrics", prometheus.Handler())
	http.ListenAndServe(*listenAddress, nil)
}
