package leveldb

import (
	"errors"
	"fmt"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
	leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"
)

const BlockchainKey = "blockchain"

type BlockchainStore struct {
	db Database
}

func NewBlockchainStore(conn Database) *BlockchainStore {
	return &BlockchainStore{
		db: conn,
	}
}

// Get todo: store as map, key is block hash
func (s *BlockchainStore) Get() (types.Blocks, error) {
	raw, err := s.db.Get([]byte(BlockchainKey), nil)
	if err != nil {
		return nil, fmt.Errorf("blockchainStore.Get get blockchain error: %w", err)
	}

	var blocks types.Blocks
	if err = rlp.DecodeBytes(raw, &blocks); err != nil {
		return nil, fmt.Errorf("failed to decode blockchain: %w", err)
	}

	return blocks, nil
}

func (s *BlockchainStore) Put(block *types.Block) error {
	blockchain, err := s.Get()
	if err != nil && !errors.Is(err, leveldberrors.ErrNotFound) {
		return fmt.Errorf("blockchainStore.Put get blockchain error: %w", err)
	}

	blockchain = append(blockchain, block)
	encoded, err := rlp.EncodeToBytes(blockchain)
	if err != nil {
		return fmt.Errorf("failed to encode blockchain: %w", err)
	}
	if err = s.db.Put([]byte(BlockchainKey), encoded, nil); err != nil {
		return fmt.Errorf("failed to put new block: %w", err)
	}
	return nil
}

func (s *BlockchainStore) Delete() error {
	if err := s.db.Delete([]byte(BlockchainKey), nil); err != nil {
		return fmt.Errorf("failed to clear blockchain: %w", err)
	}
	return nil
}
