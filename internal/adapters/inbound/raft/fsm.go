package raft

import (
	"encoding/binary"
	"fmt"
	"io"
	"local-chain/internal/adapters/outbound/inMem"

	"local-chain/internal/types"

	"github.com/hashicorp/raft"

	"local-chain/internal/adapters/outbound/leveldb"
)

type txPool interface {
	GetPool() inMem.TxPoolMap
	Purge()
}

type Fsm struct {
	store  *leveldb.Store
	txPool txPool
}

func New(store *leveldb.Store, txPool txPool) *Fsm {
	return &Fsm{
		store:  store,
		txPool: txPool,
	}
}

func (f *Fsm) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		envelope, err := types.EnvelopeFromBytes(log.Data)
		if err != nil {
			return err
		}
		switch envelope.Type {
		case types.EnvelopeTypeBlock:
			blockTxsEnvelope := types.NewBlockTxsEnvelope(nil, nil)
			if err = blockTxsEnvelope.FromBytes(envelope.Data); err != nil {
				return fmt.Errorf("failed to decode block: %w", err)
			}
			// should we check if the block already exists?
			if err = f.store.Blockchain().Put(blockTxsEnvelope.Block); err != nil {
				return fmt.Errorf("failed to save block: %w", err)
			}
			blockHash := blockTxsEnvelope.Block.ComputeHash()
			for _, tx := range blockTxsEnvelope.Txs {
				tx.BlockHash = blockHash
				err = f.store.Transaction().Put(tx)
				if err != nil {
					return fmt.Errorf("failed to put transaction: %w", err)
				}
			}
			f.txPool.Purge()
		case types.EnvelopeTypeTransaction:

		}
		return nil
	default:
		return nil
	}
}

func (f *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	blocks, err := f.store.Blockchain().Get()
	if err != nil {
		return nil, err
	}
	return &FsmSnapshot{blocks: blocks}, nil
}

func (f *Fsm) Restore(snapshot io.ReadCloser) error {
	// Clear the current blockchain blocks to avoid conflicts
	if err := f.store.Blockchain().Delete(); err != nil {
		return fmt.Errorf("failed to clear blockchain blocks: %w", err)
	}
	for {
		var length uint32
		if err := binary.Read(snapshot, binary.BigEndian, &length); err != nil {
			if err == io.EOF {
				break // End of snapshot
			}
			return fmt.Errorf("failed to read block length: %w", err)
		}
		blockBytes := make([]byte, length)
		if _, err := io.ReadFull(snapshot, blockBytes); err != nil {
			return fmt.Errorf("failed to read block data: %w", err)
		}
		block := &types.Block{}
		if err := block.FromBytes(blockBytes); err != nil {
			return fmt.Errorf("failed to deserialize block: %w", err)
		}
		if err := f.store.Blockchain().Put(block); err != nil {
			return fmt.Errorf("failed to store block %d: %w", block.Timestamp, err)
		}
	}

	if err := snapshot.Close(); err != nil {
		return fmt.Errorf("failed to close snapshot: %w", err)
	}

	return nil
}

type FsmSnapshot struct {
	blocks types.Blocks
}

func (s *FsmSnapshot) Persist(sink raft.SnapshotSink) error {
	for _, block := range s.blocks {
		blockBytes, err := block.ToBytes()
		if err != nil {
			return fmt.Errorf("failed to serialize block %d: %w", block.Timestamp, err)
		}
		// Write block length (to allow deserialization)
		length := uint32(len(blockBytes))
		if err := binary.Write(sink, binary.BigEndian, length); err != nil {
			return fmt.Errorf("failed to write block length: %w", err)
		}
		// Write block data
		if _, err := sink.Write(blockBytes); err != nil {
			return fmt.Errorf("failed to write block data: %w", err)
		}
	}
	return sink.Cancel() // nolint:errcheck
}

func (s *FsmSnapshot) Release() {
	// release resources if needed
}
