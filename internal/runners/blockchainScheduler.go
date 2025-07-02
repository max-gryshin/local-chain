package runners

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/raft"

	"local-chain/internal/types"
)

const (
	blockInterval = 10 * time.Second
)

type RaftAPI interface {
	Apply(cmd []byte, timeout time.Duration) raft.ApplyFuture
}

type blockchain interface {
	AddBlock(txs []*types.Transaction) error
}

type txPool interface {
	GetPool() map[string]*types.Transaction
}

type BlockchainScheduler struct {
	raftApi    RaftAPI
	errDone    chan struct{}
	blockchain blockchain
	txPool     txPool
}

func NewBlockchainScheduler(raftAPI RaftAPI, blockchain blockchain, txPool txPool) *BlockchainScheduler {
	return &BlockchainScheduler{
		raftApi:    raftAPI,
		blockchain: blockchain,
		txPool:     txPool,
		errDone:    make(chan struct{}, 1),
	}
}

func (bs *BlockchainScheduler) Run(ctx context.Context) error {
	tAddBlock := time.NewTicker(blockInterval)
	defer tAddBlock.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-bs.errDone:
			return errors.New("blockchain scheduler reached the error threshold")
		case <-tAddBlock.C:
			pool := bs.txPool.GetPool()
			txs := make([]*types.Transaction, 0, len(pool))
			for _, tx := range pool {
				txs = append(txs, tx)
			}
			if err := bs.blockchain.AddBlock(txs); err != nil {
				bs.errDone <- struct{}{}
				return fmt.Errorf("error while adding block: %w", err)
			}
			// todo: send the block to the raft
			//bs.raftApi.Apply()
		}
	}
}
