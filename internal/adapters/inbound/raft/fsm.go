package raft

import (
	"encoding/binary"
	"fmt"
	"io"
	"time"

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
			return fmt.Errorf("failed to decode envelope: %v", err)
		}
		switch envelope.Type {
		case types.EnvelopeTypeBlock:
			if err = f.addBlock(envelope.Data); err != nil {
				return fmt.Errorf("add block error: %v", err)
			}
		case types.EnvelopeTypeTransaction:

		}
		return nil
	default:
		return nil
	}
}

func (f *Fsm) addBlock(blockBytes []byte) error {
	blockTxsEnvelope := types.NewBlockTxsEnvelope(nil, nil)
	if err := blockTxsEnvelope.FromBytes(blockBytes); err != nil {
		return fmt.Errorf("failed to decode block: %w", err)
	}
	// should we check if the block already exists?
	if err := f.store.Blockchain().Put(blockTxsEnvelope.Block); err != nil {
		return fmt.Errorf("failed to save block: %w", err)
	}
	blockHash := blockTxsEnvelope.Block.ComputeHash()
	for _, tx := range blockTxsEnvelope.Txs {
		tx.BlockHash = blockHash
		if err := f.store.Transaction().Put(tx); err != nil {
			return fmt.Errorf("failed to put transaction: %w", err)
		}
		if len(tx.Outputs) != 2 {
			return fmt.Errorf("invalid number of outputs in transaction: %d", len(tx.Outputs))
		}
		receiver := tx.Outputs[0].PubKey
		receiverUtxos, err := f.store.Utxo().Get(receiver)
		if err != nil {
			return fmt.Errorf("failed to get utxos for receiver: %w", err)
		}
		receiverUtxos = append(receiverUtxos, &types.UTXO{TxHash: tx.GetHash(), Index: 0})
		if err = f.store.Utxo().Put(receiver, receiverUtxos); err != nil {
			return fmt.Errorf("failed to put receiver's utxo: %w", err)
		}
		if err = f.store.Utxo().Put(tx.Outputs[1].PubKey, []*types.UTXO{{TxHash: tx.GetHash(), Index: 1}}); err != nil {
			return fmt.Errorf("failed to put sender's utxo: %w", err)
		}
	}
	fmt.Printf("\nblock added - timestamp: %s\n", time.Unix(0, int64(blockTxsEnvelope.Block.Timestamp)))
	f.txPool.Purge()
	return nil
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
