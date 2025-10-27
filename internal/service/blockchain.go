package service

import (
	"context"
	"fmt"
	"time"

	"local-chain/internal/pkg"

	"local-chain/internal"

	"local-chain/internal/types"

	"github.com/hashicorp/raft"
)

const (
	applyTimeout = 10 * time.Second
)

type BlockchainStore interface {
	Get() (types.Blocks, error)
	Put(*types.Block) error
}

type RaftAPI interface {
	Apply(cmd []byte, timeout time.Duration) raft.ApplyFuture
	LeaderWithID() (raft.ServerAddress, raft.ServerID)
}

// Blockchain represents a private blockchain.
type Blockchain struct {
	raftApi          RaftAPI
	blockchainStore  BlockchainStore
	transactionStore TransactionStore
	prevBlock        *types.Block
	txPool           TxPool
}

// NewBlockchain creates a new blockchain with a genesis block.
func NewBlockchain(
	raftApi RaftAPI,
	blockchainStore BlockchainStore,
	txStore TransactionStore,
	txPool TxPool,
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
		if err = b.blockchainStore.Put(genesisBlock); err != nil {
			panic(err)
		}
	}
	for _, block := range blocks {
		if b.prevBlock == nil {
			b.prevBlock = block
		}
		if b.prevBlock.Timestamp < block.Timestamp {
			b.prevBlock = block
		}
	}

	return b
}

// CreateBlock adds a new block to the blockchain.
func (bc *Blockchain) CreateBlock(ctx context.Context) error {
	if _, leaderID := bc.raftApi.LeaderWithID(); leaderID != pkg.ServerIDFromContext(ctx) {
		return nil
	}
	txs := bc.txPool.GetPool().AsSlice()
	if len(txs) == 0 {
		return nil
	}

	merkleTree, err := internal.NewMerkleTree(txs...)
	if err != nil {
		return fmt.Errorf("failed to create merkle tree: %w", err)
	}

	block := types.NewBlock(bc.prevBlock.ComputeHash(), merkleTree.Root.Hash)
	blockTxsEnvelope := types.NewBlockTxsEnvelope(block, txs)
	bytes, err := blockTxsEnvelope.ToBytes()
	if err != nil {
		return fmt.Errorf("error while encoding block with txs: %w", err)
	}
	envelopeBytes, err := types.NewEnvelope(types.EnvelopeTypeBlock, bytes).ToBytes()
	if err != nil {
		return fmt.Errorf("error while encoding envelope: %w", err)
	}
	future := bc.raftApi.Apply(envelopeBytes, applyTimeout)
	if err = future.Error(); err != nil {
		return fmt.Errorf("error while applying block to raft: %w", err)
	}
	if response := future.Response(); response != nil {
		if err, ok := response.(error); ok {
			return fmt.Errorf("FSM failed to apply block: %w", err)
		}
	}

	bc.prevBlock = block

	return nil
}
