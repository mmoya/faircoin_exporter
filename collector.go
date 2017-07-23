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
		block, err := c.updateStateFromBlock(blockHash)
		if err != nil {
			return errors.Wrap(err, "updateStateFromBlock")
		}

		mLastBlockHeardTimestamp.Set(float64(block.Time))
		blockHash = block.NextBlockHash
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

	mCvnCount.Set(float64(len(activeCvns.Cvns)))
	c.state.lastBlockSigned = make(map[string]bool, len(activeCvns.Cvns))
	for _, cvn := range activeCvns.Cvns {
		c.state.lastBlockSigned[cvn.NodeID] = false
	}

	return activeCvns.CurrentHeight, nil
}

func (c *Collector) updateStateFromBlock(blockHash string) (*faircoin.Block, error) {
	block, err := c.fc.GetBlock(blockHash, true, 1)
	if err != nil {
		return nil, errors.Wrap(err, "GetBlock")
	}

	mLastBlockTimestamp.Set(float64(block.Time))

	if strings.Contains(block.Payload, "cvninfo") {
		log.Println("cnvinfo blog, resetting cvn state")

		_, err := c.updateCvnList()
		if err != nil {
			return nil, errors.Wrap(err, "updateCvnList from a cvninfo block")
		}
	}

	c.state.height = block.Height
	mCurrentHeight.Set(float64(block.Height))

	missingIds := make(map[string]bool)
	for _, id := range block.MissingCreatorIds {
		missingIds[id] = true
	}

	for id := range c.state.lastBlockSigned {
		_, missing := missingIds[id]

		c.state.lastBlockSigned[id] = !missing

		lastBlockSigned := 1
		if missing {
			lastBlockSigned = 0
		}

		mLastBlockSigned.WithLabelValues(id).Set(float64(lastBlockSigned))
	}

	return block, nil
}

func (c *Collector) startZmqListener() {
	c.sk, _ = zmq.NewSocket(zmq.SUB)
	c.sk.SetSubscribe("hashblock")
	c.sk.Connect(c.zmqURL)
	defer c.sk.Close()

	for {
		msg, err := c.sk.RecvMessageBytes(0)
		if err != nil {
			log.Println("zmqlistener", err)
		}

		if len(msg) < 2 {
			continue
		}

		topic := msg[0]
		body := msg[1]

		switch string(topic) {
		case "hashblock":
			mLastBlockHeardTimestamp.Set(float64(time.Now().Unix()))

			blockHash := string(hexlify(body))
			_, err := c.updateStateFromBlock(blockHash)
			if err != nil {
				log.Println("err during hashblock message", err)
			}
		}
	}
}
