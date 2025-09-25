package types

import (
	"crypto/sha512"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/rlp"
)

type Block struct {
	Timestamp  uint64
	PrevHash   []byte
	Hash       []byte
	MerkleRoot []byte
}

func NewBlock(prevHash []byte, merkleRoot []byte) *Block {
	return &Block{
		Timestamp: uint64(time.Now().UnixNano()),
		PrevHash:  prevHash,
		Hash:      merkleRoot,
	}
}

// ComputeHash computes the hash of a block.
func (b *Block) ComputeHash() []byte {
	hash := sha512.New()
	hash.Write([]byte(strconv.FormatUint(b.Timestamp, 10)))
	hash.Write(b.PrevHash)
	hash.Write(b.MerkleRoot)
	return hash.Sum(nil)
}

func (b *Block) ToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *Block) FromBytes(data []byte) error {
	return rlp.DecodeBytes(data, b)
}

type Blocks []*Block

func (b *Blocks) ToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *Blocks) FromBytes(data []byte) error {
	return rlp.DecodeBytes(data, b)
}

type BlockTxsEnvelope struct {
	Block *Block
	Txs   Transactions
}

func NewBlockTxsEnvelope(block *Block, txs Transactions) *BlockTxsEnvelope {
	return &BlockTxsEnvelope{
		Block: block,
		Txs:   txs,
	}
}

func (b *BlockTxsEnvelope) ToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(b)
}

func (b *BlockTxsEnvelope) FromBytes(data []byte) error {
	return rlp.DecodeBytes(data, b)
}
