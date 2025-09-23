package leveldb

import (
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Database interface {
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	Put(key, value []byte, wo *opt.WriteOptions) error
	Delete(key []byte, wo *opt.WriteOptions) error
}

type Store struct {
	transaction *TransactionStore
	blockchain  *BlockchainStore
	utxo        *UtxoStore
}

func New(db Database) *Store {
	return &Store{
		transaction: NewTransactionStore(db),
		blockchain:  NewBlockchainStore(db),
		utxo:        NewUtxoStore(db),
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
