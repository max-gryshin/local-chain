package types

import (
	"crypto/sha512"
	"strconv"
)

type Block struct {
	Timestamp  int64
	PrevHash   []byte
	Hash       []byte
	MerkleRoot []byte
}

// ComputeHash computes the hash of a block.
func (b *Block) ComputeHash() []byte {
	hash := sha512.New()
	hash.Write([]byte(strconv.FormatInt(b.Timestamp, 10)))
	hash.Write(b.PrevHash)
	hash.Write(b.MerkleRoot)
	return hash.Sum(nil)
}
