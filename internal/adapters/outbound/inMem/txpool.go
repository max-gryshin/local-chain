package inMem

import "local-chain/internal/types"

type TxPoolMap map[string]*types.Transaction

func (pool TxPoolMap) InSlice() []*types.Transaction {
	txs := make([]*types.Transaction, 0, len(pool))
	for _, tx := range pool {
		txs = append(txs, tx)
	}
	return txs
}

type TxPool struct {
	pool TxPoolMap
}

func NewTxPool() *TxPool {
	return &TxPool{
		pool: make(TxPoolMap),
	}
}

func (tp TxPool) AddTx(tx *types.Transaction) {
	tp.pool[string(tx.GetHash())] = tx
}

func (tp TxPool) GetPool() TxPoolMap {
	return tp.pool
}
