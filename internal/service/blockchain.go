package service

import (
	"time"

	"local-chain/internal"

	"local-chain/internal/types"
)

// Blockchain represents a private blockchain.
type Blockchain struct {
	Blocks []*types.Block
}

// NewBlockchain creates a new blockchain with a genesis block.
func NewBlockchain() *Blockchain {
	genesisBlock := &types.Block{
		Timestamp: time.Now().UnixNano(),
		PrevHash:  []byte{},
		Hash:      []byte{},
	}
	return &Blockchain{
		Blocks: []*types.Block{genesisBlock},
	}
}

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(pool *types.Pool) error {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	txs := make([]*types.Transaction, 0, len(pool.Transactions))
	for _, txPool := range pool.Transactions {
		txs = append(txs, txPool.Tx)
	}
	merkleTree, err := internal.NewMerkleTree(txs)
	if err != nil {
		return err
	}

	newBlock := &types.Block{
		Timestamp:  time.Now().UnixNano(),
		PrevHash:   prevBlock.ComputeHash(),
		MerkleRoot: merkleTree.Root.Hash,
	}
	bc.Blocks = append(bc.Blocks, newBlock)

	// todo: assign block hash to transactions
	return nil
}
