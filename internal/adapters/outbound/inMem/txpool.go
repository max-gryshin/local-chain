package inMem

import (
	"local-chain/internal/types"
	"sync"
)

type TxPoolMap map[string]*types.Transaction

func (pool TxPoolMap) AsSlice() []*types.Transaction {
	txs := make([]*types.Transaction, 0, len(pool))
	for _, tx := range pool {
		txs = append(txs, tx)
	}
	return txs
}

type TxPool struct {
	pool TxPoolMap
	mtx  sync.Mutex
}

func NewTxPool() *TxPool {
	return &TxPool{
		pool: make(TxPoolMap),
	}
}

func (tp *TxPool) AddTx(tx *types.Transaction) {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	tp.pool[string(tx.GetHash())] = tx
}

func (tp *TxPool) GetPool() TxPoolMap {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	return tp.pool
}

func (tp *TxPool) Purge() {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	tp.pool = make(TxPoolMap)
}
