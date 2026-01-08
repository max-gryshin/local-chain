package leveldb

import (
	"fmt"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
	"github.com/google/uuid"
)

type transactionS struct {
	db Database
}

func newTransactionStore(conn Database) *transactionS {
	return &transactionS{
		db: conn,
	}
}

func (s *transactionS) Get(id uuid.UUID) (*types.Transaction, error) {
	value, err := s.db.Get([]byte(id.String()), nil)
	if err != nil {
		return nil, err
	}
	tx := &types.Transaction{}
	if err = rlp.DecodeBytes(value, tx); err != nil {
		return nil, fmt.Errorf("failed to decode transaction: %w", err)
	}

	return tx, nil
}

func (s *transactionS) Put(tx *types.Transaction) error {
	encoded, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}
	if err = s.db.Put([]byte(tx.ID.String()), encoded, nil); err != nil {
		return fmt.Errorf("failed to put transaction: %w", err)
	}

	return nil
}
