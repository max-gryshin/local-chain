package service_test

import (
	"local-chain/internal/service"

	"github.com/golang/mock/gomock"
)

type MockCustomStore struct {
	TransactionStore *MockTransactionStore
	BStore           *MockBStore
	UTXOStore        *MockUTXOStore
	UserStore        *MockUserStore
	BlockTxStore     *MockBlockTxStore
}

func (m MockCustomStore) Transaction() service.TransactionStore {
	return m.TransactionStore
}

func (m MockCustomStore) Blockchain() service.BStore {
	return m.BStore
}

func (m MockCustomStore) Utxo() service.UTXOStore {
	return m.UTXOStore
}

func (m MockCustomStore) User() service.UserStore {
	return m.UserStore
}

func (m MockCustomStore) BlockTransactions() service.BlockTxStore {
	return m.BlockTxStore
}

func NewMockCustomStore(ctrl *gomock.Controller) *MockCustomStore {
	return &MockCustomStore{
		TransactionStore: NewMockTransactionStore(ctrl),
		BStore:           NewMockBStore(ctrl),
		UTXOStore:        NewMockUTXOStore(ctrl),
		UserStore:        NewMockUserStore(ctrl),
		BlockTxStore:     NewMockBlockTxStore(ctrl),
	}
}
