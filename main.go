package main

import (
	"flag"
	"net/http"

	"log"

	faircoin "github.com/mmoya/faircoin_rpcclient"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress = flag.String("listen-address", ":9132", "The address to listen on for HTTP requests")

	rpcURL      = flag.String("rpc.url", "http://127.0.0.1:40405", "URL of the RPC server")
	rpcUser     = flag.String("rpc.user", "", "User to authenticate to RPC server")
	rpcPassword = flag.String("rpc.password", "", "Password to authenticate to RPC server")

	zmqURL = flag.String("zmq.url", "tcp://127.0.0.1:28332", "URI of ZMQ server")
)

func main() {
	flag.Parse()

	log.Printf("Connecting to rpc=%s zmq=%s", *rpcURL, *zmqURL)

	setFlags := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { setFlags[f.Name] = true })

	var cred *faircoin.Credential
	if setFlags["rpc.user"] || setFlags["rpc.password"] {
		cred = faircoin.NewCredential(*rpcUser, *rpcPassword)
	} else {
		cred = faircoin.CookieCredential()
	}

	collector := NewCollector(*rpcURL, *zmqURL, cred)
	collector.Start()

	registerMetrics()

	http.Handle("/metrics", prometheus.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
				<head><title>FairCoin Exporter</title></head>
				<body>
				<h1>FairCoin Exporter</h1>
				<p><a href="/metrics">Metrics</a></p>
				</body>
				</html`))
	})

	log.Println("Listening on", *listenAddress)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
}
