package runners

import (
	"context"
	"errors"
	"fmt"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg"
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
	CreateBlock(txs []*types.Transaction) (*types.Block, error)
}

type txPool interface {
	GetPool() inMem.TxPoolMap
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
			pkg.GoWithRecover(func() {
				if err := bs.addBlock(); err != nil {
					fmt.Println("failed to add block:", err)
				}
			}, func(err error) {
				fmt.Println("failed to add block:", err)
			})
		}
	}
}

func (bs *BlockchainScheduler) addBlock() error {
	txs := bs.txPool.GetPool().InSlice()
	if len(txs) == 0 {
		return nil
	}
	block, err := bs.blockchain.CreateBlock(txs)
	if err != nil {
		bs.errDone <- struct{}{}
		return fmt.Errorf("error while adding block: %w", err)
	}
	blockBytes, err := block.ToBytes()
	if err != nil {
		bs.errDone <- struct{}{}
		return fmt.Errorf("error while encoding block: %w", err)
	}
	envelope := types.NewEnvelope(types.EnvelopeTypeBlock, blockBytes)
	envelopeBytes, err := envelope.ToBytes()
	if err != nil {
		bs.errDone <- struct{}{}
		return fmt.Errorf("error while encoding envelope: %w", err)
	}
	bs.raftApi.Apply(envelopeBytes, blockInterval)
	return nil
}
