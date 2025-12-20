package runners

import (
	"context"
	"errors"
	"fmt"
	"time"

	"local-chain/internal/pkg"
)

const (
	blockInterval = 10 * time.Second
)

type blockchain interface {
	CreateBlock(ctx context.Context) error
}

type BlockchainRunner struct {
	errDone    chan struct{}
	blockchain blockchain
}

func NewBlockchainScheduler(blockchain blockchain) *BlockchainRunner {
	return &BlockchainRunner{
		blockchain: blockchain,
		errDone:    make(chan struct{}, 1),
	}
}

func (bs *BlockchainRunner) Run(ctx context.Context) error {
	tAddBlock := time.NewTicker(blockInterval)
	defer tAddBlock.Stop()

	sem := make(chan struct{}, 1)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-bs.errDone:
			return errors.New("blockchain scheduler reached the error threshold")
		case <-tAddBlock.C:
			pkg.GoWithRecoverAndSemaphore(func() {
				if err := bs.blockchain.CreateBlock(ctx); err != nil {
					bs.errDone <- struct{}{}
					fmt.Println("failed to add block:", err)
				}
			}, func(err error) {
				fmt.Println("failed to add block:", err)
			}, sem)
		}
	}
}
