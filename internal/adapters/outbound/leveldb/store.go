package leveldb

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Database interface {
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	Put(key, value []byte, wo *opt.WriteOptions) error
	Delete(key []byte, wo *opt.WriteOptions) error
	Close() error
}

type Store struct {
	transaction *TransactionStore
	blockchain  *BlockchainStore
	utxo        *UtxoStore
	user        *UserStore
}

type dbF func(subPath string) Database

func New(newDB dbF) *Store {
	return &Store{
		transaction: NewTransactionStore(newDB("transaction")),
		blockchain:  NewBlockchainStore(newDB("blockchain")),
		utxo:        NewUtxoStore(newDB("utxo")),
		user:        NewUserStore(newDB("user")),
	}
}

func (s *Store) Transaction() *TransactionStore {
	return s.transaction
}

func (s *Store) Blockchain() *BlockchainStore {
	return s.blockchain
}

func (s *Store) Utxo() *UtxoStore {
	return s.utxo
}

func (s *Store) User() *UserStore {
	return s.user
}

func (s *Store) Close() error {
	if err := s.blockchain.db.Close(); err != nil {
		return fmt.Errorf("error closing blockchain store: %w", err)
	}

	if err := s.transaction.db.Close(); err != nil {
		return fmt.Errorf("error closing transaction store: %w", err)
	}

	if err := s.utxo.db.Close(); err != nil {
		return fmt.Errorf("error closing utxo store: %w", err)
	}

	if err := s.user.db.Close(); err != nil {
		return fmt.Errorf("error closing user store: %w", err)
	}

	return nil
}
