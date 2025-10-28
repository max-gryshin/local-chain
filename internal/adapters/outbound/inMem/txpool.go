package inMem

import (
	"slices"
	"sync"

	"local-chain/internal/types"
)

type Pool map[string]*types.Transaction
type utxosPool map[string]types.UTXOs

func (pool Pool) AsSlice() types.Transactions {
	txs := make(types.Transactions, 0, len(pool))
	for _, tx := range pool {
		txs = append(txs, tx)
	}
	return txs
}

type TxPool struct {
	// general pool with transactions by tx hash
	pool Pool
	// unspent transactions pool by owner
	utxosPool utxosPool
	mtx       sync.Mutex
}

func NewTxPool() *TxPool {
	return &TxPool{
		pool:      make(Pool),
		utxosPool: make(utxosPool),
	}
}

func (tp *TxPool) AddTx(tx *types.Transaction) {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	for index, output := range tx.Outputs {
		utxos := tp.utxosPool[string(output.PubKey)]
		if index > 0 {
			for _, txFromPool := range tp.pool {
				if len(txFromPool.Outputs) > 0 && string(txFromPool.Outputs[1].PubKey) == string(output.PubKey) {
					txFromPool.Outputs[1].Amount = *types.NewAmount(0)
				}
			}
			utxos = slices.DeleteFunc(utxos, func(utxo *types.UTXO) bool { return utxo.Index > 0 })
		}
		tp.utxosPool[string(output.PubKey)] = append(utxos, types.NewUTXO(tx.GetHash(), uint32(index)))
	}
	tp.pool[string(tx.GetHash())] = tx
}

func (tp *TxPool) GetPool() Pool {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	return tp.pool
}

func (tp *TxPool) GetUTXOs(pubKey []byte) types.UTXOs {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	utxos, ok := tp.utxosPool[string(pubKey)]
	if !ok {
		return nil
	}
	return utxos
}

func (tp *TxPool) GetUnspentTx() Pool {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	return tp.pool
}

func (tp *TxPool) Purge() {
	tp.mtx.Lock()
	defer tp.mtx.Unlock()
	tp.pool = make(Pool)
	tp.utxosPool = make(utxosPool)
}
