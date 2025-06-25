package service

import (
	"fmt"
	"local-chain/internal/adapters/outbound/leveldb"
	"time"

	"local-chain/internal"

	"local-chain/internal/types"
)

// Blockchain represents a private blockchain.
type Blockchain struct {
	BlockchainStore  leveldb.BlockchainStore
	TransactionStore leveldb.TransactionStore
	Blocks           []*types.Block
}

// NewBlockchain creates a new blockchain with a genesis block.
func NewBlockchain(blockchainStore *leveldb.BlockchainStore, txStore *leveldb.TransactionStore) *Blockchain {
	b := &Blockchain{
		BlockchainStore:  *blockchainStore,
		TransactionStore: *txStore,
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
		Timestamp:  uint64(time.Now().UnixNano()),
		PrevHash:   prevBlock.ComputeHash(),
		MerkleRoot: merkleTree.Root.Hash,
	}
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
