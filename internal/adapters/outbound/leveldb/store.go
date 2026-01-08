package leveldb

import (
	"fmt"

	"local-chain/internal/service"

	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Database interface {
	Get(key []byte, ro *opt.ReadOptions) (value []byte, err error)
	Put(key, value []byte, wo *opt.WriteOptions) error
	Delete(key []byte, wo *opt.WriteOptions) error
	NewIterator(slice *util.Range, ro *opt.ReadOptions) iterator.Iterator
	Close() error
}

type Store struct {
	transaction       *transactionS
	blockchain        *blockchainS
	utxo              *utxoS
	user              *userS
	blockTransactions *blockTransactionsS
}

type dbF func(subPath string) Database

func New(newDB dbF) *Store {
	return &Store{
		transaction:       newTransactionStore(newDB("transaction")),
		blockchain:        newBlockchainStore(newDB("blockchain")),
		utxo:              newUtxoStore(newDB("utxo")),
		user:              newUserStore(newDB("user")),
		blockTransactions: newBlockTransactionsStore(newDB("block_transactions")),
	}
}

func (s *Store) Transaction() service.TransactionStore {
	return s.transaction
}

func (s *Store) Blockchain() service.BStore {
	return s.blockchain
}

func (s *Store) Utxo() service.UTXOStore {
	return s.utxo
}

func (s *Store) User() service.UserStore {
	return s.user
}

func (s *Store) BlockTransactions() service.BlockTxStore {
	return s.blockTransactions
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
