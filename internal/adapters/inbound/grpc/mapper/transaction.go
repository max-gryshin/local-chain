package mapper

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"local-chain/internal/types"
	grpcPkg "local-chain/transport/gen/transport"
)

type TransactionMapper struct{}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (tp *TransactionMapper) RpcToTransaction(req *grpcPkg.AddTransactionRequest) (*types.TransactionRequest, error) {
	amount := types.Amount{Value: req.GetAmount().GetValue(), Unit: req.GetAmount().GetUnit()}
	publicBlock, _ := pem.Decode([]byte(req.Receiver))
	if publicBlock == nil || publicBlock.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("invalid public key PEM")
	}
	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %v", err)
	}
	receiver, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not ECDSA")
	}

	privateBlock, _ := pem.Decode([]byte(req.Sender))
	if privateBlock == nil || privateBlock.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("invalid private key PEM")
	}
	sender, err := x509.ParseECPrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return &types.TransactionRequest{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
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
