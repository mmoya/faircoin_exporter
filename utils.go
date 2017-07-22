package main

import "encoding/hex"

func hexlify(in []byte) []byte {
	out := make([]byte, hex.EncodedLen(len(in)))
	hex.Encode(out, in)
	return out
}
