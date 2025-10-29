package leveldb

import (
	"errors"
	"fmt"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
)

type UtxoStore struct {
	db Database
}

func NewUtxoStore(conn Database) *UtxoStore {
	return &UtxoStore{
		db: conn,
	}
}

func (s *UtxoStore) Get(pubKey []byte) ([]*types.UTXO, error) {
	utxos := make([]*types.UTXO, 0)
	value, err := s.db.Get(pubKey, nil)
	if err != nil {
		if errors.As(err, &leveldbErrors.ErrNotFound) { // nolint:govet
			return nil, nil
		}
		return utxos, fmt.Errorf("failed to get utxos: %w", err)
	}
	if err = rlp.DecodeBytes(value, &utxos); err != nil {
		return nil, fmt.Errorf("failed to decode utxos: %w", err)
	}

	return utxos, nil
}

func (s *UtxoStore) Put(pubKey []byte, utxos ...*types.UTXO) error {
	encoded, err := rlp.EncodeToBytes(utxos)
	if err != nil {
		return fmt.Errorf("failed to encode utxos: %w", err)
	}
	if err = s.db.Put(pubKey, encoded, nil); err != nil {
		return fmt.Errorf("failed to put utxos: %w", err)
	}

	return nil
}
