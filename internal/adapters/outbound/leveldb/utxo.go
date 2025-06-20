package leveldb

import (
	"fmt"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
)

type UtxoStore struct {
	db Database
}

func NewUtxoStore(conn Database) *UtxoStore {
	return &UtxoStore{
		db: conn,
	}
}

func (s *UtxoStore) Get(pubKeyHash []byte) ([]*types.UTXO, error) {
	value, err := s.db.Get(pubKeyHash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get utxos: %w", err)
	}
	var utxos []*types.UTXO
	if err = rlp.DecodeBytes(value, &utxos); err != nil {
		return nil, fmt.Errorf("failed to decode utxos: %w", err)
	}

	return utxos, nil
}

func (s *UtxoStore) Put(pubKeyHash []byte, utxos []*types.UTXO) error {
	encoded, err := rlp.EncodeToBytes(utxos)
	if err != nil {
		return fmt.Errorf("failed to encode utxos: %w", err)
	}
	if err = s.db.Put(pubKeyHash, encoded, nil); err != nil {
		return fmt.Errorf("failed to put utxos: %w", err)
	}

	return nil
}
