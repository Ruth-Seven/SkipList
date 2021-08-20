package skiplist

import (
	"crypto/rand"
	"log"
	"math/big"
)

const (
	FormatLen  = 5
	FormatBits = 64
)

func randUint64(limit uint64) uint64 {
	max := &big.Int{}
	max.SetUint64(limit)
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		log.Printf("randUint64: fail to generate rand number. err: %v", err)
	}
	return num.Uint64()
}
