package service

import (
	"fmt"
	"local-chain/internal/adapters/outbound/inMem"
	"time"

	"local-chain/internal"

	"local-chain/internal/types"

	"github.com/hashicorp/raft"
)

const (
	applyTimeout = 1 * time.Minute
)

type BlockchainStore interface {
	Get() ([]*types.Block, error)
	Put(*types.Block) error
}

type txPool interface {
	GetPool() inMem.TxPoolMap
	Purge()
}

type RaftAPI interface {
	Apply(cmd []byte, timeout time.Duration) raft.ApplyFuture
}

// Blockchain represents a private blockchain.
type Blockchain struct {
	raftApi          RaftAPI
	blockchainStore  BlockchainStore
	transactionStore TransactionStore
	Blocks           []*types.Block
	txPool           txPool
}

// NewBlockchain creates a new blockchain with a genesis block.
func NewBlockchain(
	raftApi RaftAPI,
	blockchainStore BlockchainStore,
	txStore TransactionStore,
	txPool txPool,
) *Blockchain {
	b := &Blockchain{
		raftApi:          raftApi,
		blockchainStore:  blockchainStore,
		transactionStore: txStore,
		txPool:           txPool,
	}
	blocks, err := b.blockchainStore.Get()
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

// CreateBlock adds a new block to the blockchain.
func (bc *Blockchain) CreateBlock() error {
	//todo: get last block from store
	prevBlock := bc.Blocks[len(bc.Blocks)-1]

	txs := bc.txPool.GetPool().AsSlice()
	if len(txs) == 0 {
		return nil
	}

	merkleTree, err := internal.NewMerkleTree(txs)
	if err != nil {
		return fmt.Errorf("failed to create merkle tree: %w", err)
	}

	newBlock := &types.Block{
		Timestamp:  uint64(time.Now().UnixNano()),
		PrevHash:   prevBlock.ComputeHash(),
		MerkleRoot: merkleTree.Root.Hash,
	}
	//todo: get blocks from store
	bc.Blocks = append(bc.Blocks, newBlock)

	err = bc.blockchainStore.Put(newBlock)
	if err != nil {
		return fmt.Errorf("failed to put new block: %w", err)
	}
	blockHash := newBlock.ComputeHash()
	for _, tx := range txs {
		tx.BlockHash = blockHash
		err = bc.transactionStore.Put(tx)
		if err != nil {
			return fmt.Errorf("failed to put transaction: %w", err)
		}
	}

	blockBytes, err := newBlock.ToBytes()
	if err != nil {
		return fmt.Errorf("error while encoding block: %w", err)
	}
	envelopeBytes, err := types.NewEnvelope(types.EnvelopeTypeBlock, blockBytes).ToBytes()
	if err != nil {
		return fmt.Errorf("error while encoding envelope: %w", err)
	}
	if err := bc.raftApi.Apply(envelopeBytes, applyTimeout).Error(); err != nil {
		return fmt.Errorf("error while applying block to raft: %w", err)
	}
	bc.txPool.Purge()

	return nil
}
