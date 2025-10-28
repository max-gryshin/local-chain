package inMem

import (
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
	tp.pool[string(tx.GetHash())] = tx
	for index, output := range tx.Outputs {
		if index == 0 {
			utxos := tp.utxosPool[string(output.PubKey)]
			tp.utxosPool[string(output.PubKey)] = append(utxos, types.NewUTXO(tx.GetHash(), uint32(index)))
			continue
		}
		tp.utxosPool[string(output.PubKey)] = types.UTXOs{types.NewUTXO(tx.GetHash(), uint32(index))}
	}
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
