package leveldb

import (
	"fmt"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
)

type TransactionStore struct {
	db Database
}

func NewTransactionStore(conn Database) *TransactionStore {
	return &TransactionStore{
		db: conn,
	}
}

func (s *TransactionStore) Get(txHash []byte) (*types.Transaction, error) {
	value, err := s.db.Get(txHash, nil)
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{}
	if err = rlp.DecodeBytes(value, tx); err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}

	return tx, nil
}

func (s *TransactionStore) Put(tx *types.Transaction) error {
	encoded, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}
	if err = s.db.Put(tx.GetHash(), encoded, nil); err != nil {
		return fmt.Errorf("failed to put transaction: %w", err)
	}

	return nil
}
