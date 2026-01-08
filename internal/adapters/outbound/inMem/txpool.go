package inMem

import (
	"sync"

	"local-chain/internal/types"

	"github.com/google/uuid"
)

type (
	Pool      map[uuid.UUID]*types.Transaction
	utxosPool map[string]types.UTXOs
)

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

func (txp *TxPool) AddTx(tx *types.Transaction) error {
	txp.mtx.Lock()
	defer txp.mtx.Unlock()

	txp.pool[tx.ID] = tx
	return nil
}

func (txp *TxPool) AddUtxos(pubKey []byte, utxos ...*types.UTXO) {
	txp.mtx.Lock()
	defer txp.mtx.Unlock()
	txp.utxosPool[string(pubKey)] = utxos
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
