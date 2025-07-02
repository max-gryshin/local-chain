package service

import (
	"fmt"
	"time"

	"local-chain/internal"

	"local-chain/internal/types"
)

type BlockchainStore interface {
	Get() ([]*types.Block, error)
	Put(*types.Block) error
}

// Blockchain represents a private blockchain.
type Blockchain struct {
	BlockchainStore  BlockchainStore
	TransactionStore TransactionStore
	Blocks           []*types.Block
}

// NewBlockchain creates a new blockchain with a genesis block.
func NewBlockchain(blockchainStore BlockchainStore, txStore TransactionStore) *Blockchain {
	b := &Blockchain{
		BlockchainStore:  blockchainStore,
		TransactionStore: txStore,
	}
	blocks, err := b.BlockchainStore.Get()
	if err != nil {
		panic(err)
	}
	if len(blocks) == 0 {
		genesisBlock := &types.Block{
			Timestamp: uint64(time.Now().UnixNano()),
			PrevHash:  []byte{},
			Hash:      []byte{},
		}
		b.Blocks = append(b.Blocks, genesisBlock)
	} else {
		b.Blocks = blocks
	}

	return b
}

// AddBlock adds a new block to the blockchain.
func (bc *Blockchain) AddBlock(txs []*types.Transaction) error {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	merkleTree, err := internal.NewMerkleTree(txs)
	if err != nil {
		return err
	}

	newBlock := &types.Block{
		Timestamp:  uint64(time.Now().UnixNano()),
		PrevHash:   prevBlock.ComputeHash(),
		MerkleRoot: merkleTree.Root.Hash,
	}
	//todo: get blocks from store
	bc.Blocks = append(bc.Blocks, newBlock)

	err = bc.BlockchainStore.Put(newBlock)
	if err != nil {
		return fmt.Errorf("failed to put new block: %w", err)
	}
	blockHash := newBlock.ComputeHash()
	for _, tx := range txs {
		tx.BlockHash = blockHash
		err = bc.TransactionStore.Put(tx)
		if err != nil {
			return fmt.Errorf("failed to put transaction: %w", err)
		}
	}

	return nil
}
