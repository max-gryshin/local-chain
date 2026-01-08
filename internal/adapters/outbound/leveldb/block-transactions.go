package leveldb

import (
	"errors"
	"fmt"
	"local-chain/internal/types"
	"strconv"

	"github.com/ethereum/go-ethereum/rlp"
	leveldberrors "github.com/syndtr/goleveldb/leveldb/errors"
)

type blockTransactionsS struct {
	db Database
}

func newBlockTransactionsStore(conn Database) *blockTransactionsS {
	return &blockTransactionsS{
		db: conn,
	}
}

func (s *blockTransactionsS) Put(envelope *types.BlockTxsEnvelope) error {
	txs, err := rlp.EncodeToBytes(envelope.Txs)
	if err != nil {
		return fmt.Errorf("failed to encode rlp: %w", err)
	}
	if err = s.db.Put([]byte(strconv.Itoa(int(envelope.Block.Timestamp))), txs, nil); err != nil {
		return fmt.Errorf("failed to put block transactions: %w", err)
	}
	return nil
}

func (s *blockTransactionsS) GetByBlockTimestamp(t uint64) (types.Transactions, error) {
	raw, err := s.db.Get([]byte(strconv.Itoa(int(t))), nil)
	if err != nil && !errors.Is(err, leveldberrors.ErrNotFound) {
		return nil, fmt.Errorf("blockTransactionsStore.GetByBlockTimestamp get block transactions error: %w", err)
	}
	if raw == nil {
		return nil, ErrNotFound
	}

	var txs types.Transactions
	if err = rlp.DecodeBytes(raw, &txs); err != nil {
		return nil, fmt.Errorf("failed to decode block transactions: %w", err)
	}

	return txs, nil
}
