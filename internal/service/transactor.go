package service

import (
	"errors"
	"fmt"

	"local-chain/internal/types"
)

//go:generate mockgen -source transactor.go -destination transactor_mock_test.go -package service_test -mock_names TransactionStore=MockTransactionStore

type TransactionStore interface {
	Get(txHash []byte) (*types.Transaction, error)
	Put(*types.Transaction) error
}

type Transactor struct {
	txStore TransactionStore
}

func NewTransactor(txStore TransactionStore) *Transactor {
	return &Transactor{
		txStore: txStore,
	}
}

// CreateTx Creates new transaction
// The fact of checking the ownership of the output of the used transaction is the ability to sign the inputs
// of the new transaction, since it is impossible to sign the inputs of the new transaction without having the
// private key corresponding to the public key.
func (t *Transactor) CreateTx(txReq *types.TransactionRequest) (*types.Transaction, error) {
	newTx := types.NewTransaction()
	inputs := make([]*types.TxIn, 0, len(txReq.Utxos))
	balance := types.Amount{}
	for id, utxo := range txReq.Utxos {
		tx, err := t.txStore.Get(utxo.TxHash)
		if err != nil {
			return nil, fmt.Errorf("get utxo tx hash err: %v", err)
		}
		if int(utxo.Index) >= len(tx.Outputs) {
			return nil, fmt.Errorf("UTXO index %d is out of bounds for transaction %s", utxo.Index, string(utxo.TxHash))
		}
		output := tx.Outputs[utxo.Index]
		// Check if the output belongs to the sender
		if output.PubKey != types.ToHashString(&txReq.Sender.PublicKey) {
			return nil, fmt.Errorf("sender do not own transaction's output: tx:%s", string(utxo.TxHash))
		}

		r, s, err := utxo.Sign(txReq.Sender)
		if err != nil {
			return nil, fmt.Errorf("sign UTXO:%s err:%s", string(utxo.TxHash), err.Error())
		}

		if !utxo.Verify(txReq.Sender.PublicKey, r, s) {
			return nil, fmt.Errorf("can not verify UTXO:%s err: not valid private key", string(utxo.TxHash))
		}

		input := &types.TxIn{
			Prev:       utxo,
			PubKey:     &txReq.Sender.PublicKey,
			SignatureR: r,
			SignatureS: s,
			NSequence:  uint32(id),
		}
		inputs = append(inputs, input)
		// calculate balance
		balance.Value += output.Amount.Value
	}

	if balance.Value < txReq.Amount.Value {
		return nil, errors.New("insufficient balance")
	}
	newTx = newTx.WithInputs(inputs...)
	var outputs []*types.TxOut
	outputs = append(outputs, &types.TxOut{TxID: newTx.ID, Amount: txReq.Amount, PubKey: types.ToHashString(txReq.Receiver)})
	if balance.Value > txReq.Amount.Value {
		balance.Value -= txReq.Amount.Value
		outputs = append(outputs, &types.TxOut{TxID: newTx.ID, Amount: balance, PubKey: types.ToHashString(&txReq.Sender.PublicKey)})
	}
	newTx = newTx.WithOutputs(outputs...)
	newTx.ComputeHash()

	return newTx, nil
}
