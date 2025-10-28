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
	}, nil
}

func (tp *TransactionMapper) RpcToBalanceRequest(req *grpcPkg.GetBalanceRequest) (*types.BalanceRequest, error) {
	sender, err := crypto.PrivateKeyFromBytes(req.Sender)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return &types.BalanceRequest{
		Sender: sender,
	}, nil
}
