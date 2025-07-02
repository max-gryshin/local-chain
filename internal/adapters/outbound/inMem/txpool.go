package inMem

import "local-chain/internal/types"

type TxPool struct {
	pool map[string]*types.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		pool: make(map[string]*types.Transaction),
	}
}

func (tp TxPool) AddTx(tx *types.Transaction) error {
	tp.pool[string(tx.GetHash())] = tx
	return nil
}

func (tp TxPool) GetPool() map[string]*types.Transaction {
	return tp.pool
}
