package client

// GetBlockCount is the wrapper around getblockcount rpc call
func (c *FC2Client) GetBlockCount() (int, error) {
	response, err := c.c.Call("getblockcount")
	if err != nil {
		return 0, err
	}

	var blockCount int
	err = response.GetObject(&blockCount)
	if err != nil {
		return 0, err
	}

	return blockCount, nil
}
