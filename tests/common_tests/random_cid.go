package common_tests

import (
	"crypto/rand"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func RandomCid() string {
	data := make([]byte, 32)
	_, _ = rand.Read(data)

	hash, err := multihash.Sum(data, multihash.SHA2_256, -1)
	if err != nil {
		panic(err)
	}

	return cid.NewCidV1(cid.Raw, hash).String()
}
