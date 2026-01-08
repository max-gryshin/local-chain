package leveldb

import (
	"errors"
	"fmt"
	"strconv"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
	leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"
)

type blockchainS struct {
	db Database
}

func newBlockchainStore(conn Database) *blockchainS {
	return &blockchainS{
		db: conn,
	}
}

func (s *blockchainS) GetAll() (types.Blocks, error) {
	keys, err := s.getKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}
	blocks := make(types.Blocks, 0, len(keys))
	for _, key := range keys {
		raw, err := s.db.Get(key, nil)
		if err != nil && !errors.As(err, &leveldberrors.ErrNotFound) { // nolint:govet
			return nil, fmt.Errorf("blockchainStore.GetAll get block error: %w", err)
		}
		if raw == nil {
			return nil, nil
		}

		var block *types.Block
		if err = rlp.DecodeBytes(raw, &block); err != nil {
			return nil, fmt.Errorf("failed to decode block: %w", err)
		}
		blocks = append(blocks, block)
	}

	return blocks, nil
}

func (s *blockchainS) GetByTimestamp(t uint64) (*types.Block, error) {
	raw, err := s.db.Get([]byte(strconv.Itoa(int(t))), nil)
	if err != nil && !errors.As(err, &leveldberrors.ErrNotFound) { // nolint:govet
		return nil, fmt.Errorf("blockchainStore.GetByTimestamp get block error: %w", err)
	}
	if raw == nil {
		return nil, nil
	}

	var block *types.Block
	if err = rlp.DecodeBytes(raw, &block); err != nil {
		return nil, fmt.Errorf("failed to decode block: %w", err)
	}

	return block, nil
}

func (s *blockchainS) Put(block *types.Block) error {
	existingBlock, err := s.GetByTimestamp(block.Timestamp)
	if err != nil {
		return fmt.Errorf("blockchainStore.Put get existing block error: %w", err)
	}
	if existingBlock != nil {
		return fmt.Errorf("blockchainStore.Put blockchain with timestamp %d already exists", block.Timestamp)
	}
	encoded, err := rlp.EncodeToBytes(block)
	if err != nil {
		return fmt.Errorf("failed to encode block: %w", err)
	}
	if err = s.db.Put([]byte(strconv.Itoa(int(block.Timestamp))), encoded, nil); err != nil {
		return fmt.Errorf("failed to put new block: %w", err)
	}
	return nil
}

func (s *blockchainS) Delete() error {
	keys, err := s.getKeys()
	if err != nil {
		return fmt.Errorf("failed to get keys for deletion: %w", err)
	}

	for _, key := range keys {
		if err := s.db.Delete(key, nil); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", string(key), err)
		}
	}

	return nil
}

func (s *blockchainS) GetKeys() ([]uint64, error) {
	keys, err := s.getKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	var timestamps []uint64
	for _, key := range keys {
		timestamp, err := strconv.ParseUint(string(key), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key %s to uint64: %w", string(key), err)
		}
		timestamps = append(timestamps, timestamp)
	}

	return timestamps, nil
}

func (s *blockchainS) getKeys() ([][]byte, error) {
	iterator := s.db.NewIterator(nil, nil)
	defer iterator.Release()

	var keys [][]byte
	for iterator.Next() {
		// Make a copy of the key since the iterator reuses the buffer
		key := make([]byte, len(iterator.Key()))
		copy(key, iterator.Key())
		keys = append(keys, key)
	}

	if err := iterator.Error(); err != nil {
		return nil, fmt.Errorf("failed to iterate over keys: %w", err)
	}
	return keys, nil
}
