package txPool

import "local-chain/internal/types"

type Store interface {
	Put(key, value []byte)
	Get(key []byte) interface{}
	Delete(key []byte)
}
type TxPool struct {
	store Store
}

func NewTxPool(store Store) *TxPool {
	return &TxPool{
		store: store,
	}
}

// нужно подумать как хранить блокчейн
// 1. посмотреть как это делают в биткоине и эфириуме
// самый главный вопрос это формат
// как будто бы база должна состоять отдельно из цепочки блоков и транзакций (транзакции надо добавить id блока many-to-one)
//
// todo: заюзать key value db
// todo: implement serialization/deserialization - избыточно делать умное займет много времени (отложить на post MVP) покамест encode/decode в json
func (tp TxPool) Get(txHash []byte) *types.Transaction {
	res := tp.store.Get(txHash)

	// todo: implement serialization/deserialization
	return res.(*types.Transaction)
}
