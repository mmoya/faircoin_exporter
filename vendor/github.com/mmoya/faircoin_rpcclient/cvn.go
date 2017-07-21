package client

import "github.com/pkg/errors"

// Cvn is the state of a CVN as returned by getactivecvn rpc call
type Cvn struct {
	NodeID             string
	PubKey             string
	HeightAdded        int
	PredictedNextBlock int
	LastBlocksSigned   int
}

// ActiveCvns is the response of getactivecvn rpc call
type ActiveCvns struct {
	Count         int
	CurrentHeight int
	Cvns          []Cvn
}

// GetActiveCVNs is the wrapper around getactivecvn rpc call
func (c *Client) GetActiveCVNs() (*ActiveCvns, error) {
	response, err := c.c.Call("getactivecvns")
	if err != nil {
		return nil, errors.Wrap(err, "error calling getactivecvns")
	}

	activeCvns := ActiveCvns{}
	err = response.GetObject(&activeCvns)
	if err != nil {
		return nil, errors.Wrap(err, "error unmarshaling")
	}

	return &activeCvns, nil
}
