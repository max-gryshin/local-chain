package types

import "github.com/google/uuid"

type Pool struct {
	Transactions []*TransactionPool
}

type TransactionPool struct {
	NodeID uuid.UUID
	Tx     *Transaction
}

func NewPool() *Pool {
	return &Pool{
		Transactions: []*TransactionPool{},
	}
}

func (p *Pool) Add(nodeID uuid.UUID, tx *Transaction) {
	p.Transactions = append(p.Transactions, &TransactionPool{NodeID: nodeID, Tx: tx})
}
