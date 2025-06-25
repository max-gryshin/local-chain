package txPool

import "local-chain/internal/types"

type TxPool struct {
	pool map[string]*types.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{}
}

func (tp TxPool) AddTx(tx *types.Transaction) error {
	tp.pool[string(tx.GetHash())] = tx
	return nil
}
