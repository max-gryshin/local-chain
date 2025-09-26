package mapper

import (
	"fmt"

	"local-chain/internal/pkg/crypto"

	grpcPkg "local-chain/transport/gen/transport"

	"local-chain/internal/types"
)

type TransactionMapper struct{}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (tp *TransactionMapper) RpcToTransaction(req *grpcPkg.AddTransactionRequest) (*types.TransactionRequest, error) {
	receiver, err := crypto.PublicKeyFromBytes(req.Receiver)
	if err != nil {
		return nil, fmt.Errorf("public key is not ECDSA")
	}
	sender, err := crypto.PrivateKeyFromBytes(req.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return &types.TransactionRequest{
		Sender:   sender,
		Receiver: receiver,
		Amount:   types.Amount{Value: req.GetAmount().GetValue(), Unit: req.GetAmount().GetUnit()},
		Utxos:    rpcToUtxos(req.GetUtxos()),
	}, nil
}

func rpcToUtxos(rpcUtxos []*grpcPkg.Utxo) []*types.UTXO {
	utxos := make([]*types.UTXO, 0, len(rpcUtxos))
	for _, rpcUtxo := range rpcUtxos {
		utxo := &types.UTXO{
			TxHash: rpcUtxo.TxHash,
			Index:  rpcUtxo.Index,
		}
		utxos = append(utxos, utxo)
	}
	return utxos
}
