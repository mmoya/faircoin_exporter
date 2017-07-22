package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/pkg/errors"

	faircoin "github.com/mmoya/faircoin_rpcclient"
	zmq "github.com/pebbe/zmq4"
)

type cvnState struct {
	synchronized    bool
	height          int
	lastBlockSigned map[string]bool
}

// Collector represents a FairCoin metrics collector
type Collector struct {
	zmqURL string

	fc *faircoin.Client
	sk *zmq.Socket

	state cvnState
}

// NewCollector returns a new Collector instance
func NewCollector(rpcURL, zmqURL string, cred *faircoin.Credential) *Collector {
	fc := faircoin.NewClient(rpcURL, cred)
	collector := Collector{
		zmqURL: zmqURL,
		fc:     fc,
	}

	return &collector
}

// Start collection of metrics
func (c *Collector) Start() {
	c.syncState()
	go c.startZmqListener()
}

func (c *Collector) syncState() {
	err := c.syncStateOnce()

	for err != nil {
		fmt.Println("syncStateOnce", err)
		time.Sleep(3 * time.Second)

		err = c.syncStateOnce()
	}
}

func (c *Collector) syncStateOnce() error {
	currentHeight, err := c.updateCvnList()
	if err != nil {
		return errors.Wrap(err, "updateCvnList")
	}

	blockHash, err := c.fc.GetBlockHash(currentHeight)
	if err != nil {
		return errors.Wrap(err, "GetBlockHash")
	}

	for len(blockHash) > 0 {
		blockHash, err = c.updateStateFromBlock(blockHash)
		if err != nil {
			return errors.Wrap(err, "updateStateFromBlock")
		}
	}

	c.state.synchronized = true
	return nil
}

func (c *Collector) updateCvnList() (int, error) {
	c.state.synchronized = false

	activeCvns, err := c.fc.GetActiveCVNs()
	if err != nil {
		return 0, err
	}

	currentHeight := activeCvns.CurrentHeight
	c.state.lastBlockSigned = make(map[string]bool)
	for _, cvn := range activeCvns.Cvns {
		c.state.lastBlockSigned[cvn.NodeID] = false
	}

	return currentHeight, nil
}

func (c *Collector) updateStateFromBlock(blockHash string) (string, error) {
	block, err := c.fc.GetBlock(blockHash, true, 1)
	if err != nil {
		return "", errors.Wrap(err, "GetBlock")
	}

	if strings.Contains(block.Payload, "cvninfo") {
		_, err := c.updateCvnList()
		if err != nil {
			return "", errors.Wrap(err, "updateCvnList from a cvninfo block")
		}
	}

	c.state.height = block.Height
	currentHeight.Set(float64(block.Height))

	missingIds := make(map[string]bool)
	for _, id := range block.MissingCreatorIds {
		missingIds[id] = true
	}

	for id := range c.state.lastBlockSigned {
		_, missing := missingIds[id]

		c.state.lastBlockSigned[id] = !missing

		lbs := 1
		if missing {
			lbs = 0
		}

		lastBlockSigned.WithLabelValues(id).Set(float64(lbs))
	}

	return block.NextBlockHash, nil
}

func (c *Collector) startZmqListener() {
	c.sk, _ = zmq.NewSocket(zmq.SUB)
	c.sk.SetSubscribe("hashblock")
	c.sk.Connect(c.zmqURL)
	defer c.sk.Close()

	for {
		msg, err := c.sk.RecvMessageBytes(0)
		if err != nil {
			log.Println("inside zmqlistener", err)
		}

		if len(msg) < 2 {
			continue
		}

		topic := msg[0]
		body := msg[1]

		switch string(topic) {
		case "hashblock":
			blockHash := string(hexlify(body))
			_, err := c.updateStateFromBlock(blockHash)
			if err != nil {
				log.Println("err during hashblock message", err)
			}
		}
	}
}
