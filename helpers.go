package main

import (
	"encoding/hex"
	"log"
)

func hexToString(in []byte) string {
	dst := make([]byte, hex.DecodedLen(len(in)))
	n, err := hex.Decode(dst, in)
	if err != nil {
		log.Fatal(err)
	}

	return string(dst[:n])
}
