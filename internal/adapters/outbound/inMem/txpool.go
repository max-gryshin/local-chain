package inMem

import (
	"local-chain/internal/types"
	"sync"
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

func (txp *TxPool) AddTx(tx *types.Transaction) {
	txp.mtx.Lock()
	defer txp.mtx.Unlock()

	for index, output := range tx.Outputs {
		utxos := txp.utxosPool[string(output.PubKey)]
		if index > 0 {
			for _, txFromPool := range txp.pool {
				if len(txFromPool.Outputs) > 0 &&
					string(txFromPool.Outputs[1].PubKey) == string(output.PubKey) {
					txFromPool.Outputs[1].Amount = *types.NewAmount(0)
				}
			}
		}
		txp.utxosPool[string(output.PubKey)] = append(utxos, types.NewUTXO(tx.GetHash(), uint32(index)))
	}
	txp.pool[string(tx.GetHash())] = tx
}

func (txp *TxPool) GetPool() Pool {
	txp.mtx.Lock()
	defer txp.mtx.Unlock()
	return txp.pool
}

func (txp *TxPool) GetUTXOs(pubKey []byte) types.UTXOs {
	txp.mtx.Lock()
	defer txp.mtx.Unlock()
	utxos, ok := txp.utxosPool[string(pubKey)]
	if !ok {
		return nil
	}
	return utxos
}

func (txp *TxPool) GetUnspentTx() Pool {
	txp.mtx.Lock()
	defer txp.mtx.Unlock()
	return txp.pool
}

func (txp *TxPool) Purge() {
	txp.mtx.Lock()
	defer txp.mtx.Unlock()
	txp.pool = make(Pool)
	txp.utxosPool = make(utxosPool)
}
