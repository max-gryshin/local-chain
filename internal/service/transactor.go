package service

import (
	"errors"
	"fmt"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg/crypto"

	"local-chain/internal/types"
)

//go:generate mockgen -source transactor.go -destination transactor_mock_test.go -package service_test -mock_names TransactionStore=MockTransactionStore,TxPool=MockTxPool

type TransactionStore interface {
	Get(txHash []byte) (*types.Transaction, error)
	Put(*types.Transaction) error
}

type UTXOStore interface {
	Get(pubKey []byte) ([]*types.UTXO, error)
	Put(pubKey []byte, utxos []*types.UTXO) error
}

type TxPool interface {
	GetPool() inMem.Pool
	Purge()
	AddTx(tx *types.Transaction)
	GetUTXOs(pubKey []byte) types.UTXOs
}

type Transactor struct {
	txStore   TransactionStore
	utxoStore UTXOStore
	txPool    TxPool
}

func NewTransactor(txStore TransactionStore, UTXOStore UTXOStore, txPool TxPool) *Transactor {
	return &Transactor{
		txStore:   txStore,
		utxoStore: UTXOStore,
		txPool:    txPool,
	}
}

// CreateTx Creates new transaction
// The fact of checking the ownership of the output of the used transaction is the ability to sign the inputs
// of the new transaction, since it is impossible to sign the inputs of the new transaction without having the
// private key corresponding to the public key.
func (t *Transactor) CreateTx(txReq *types.TransactionRequest) (*types.Transaction, error) {
	var (
		err          error
		newTx        = types.NewTransaction()
		balance      = types.Amount{}
		senderPubKey = crypto.PublicKeyToBytes(&txReq.Sender.PublicKey)
		// check is utxos exists in TxPool for sender - to prevent double spending
		utxos = t.txPool.GetUTXOs(senderPubKey)
	)
	if utxos == nil {
		utxos, err = t.utxoStore.Get(senderPubKey)
		if err != nil {
			return nil, fmt.Errorf("error getting sender's utxos : %v", err)
		}
	}
	for id, utxo := range utxos {
		tx, err := t.txStore.Get(utxo.TxHash)
		if err != nil {
			return nil, fmt.Errorf("get utxo tx hash err: %v", err)
		}
		if int(utxo.Index) >= len(tx.Outputs) {
			return nil, fmt.Errorf("UTXO index %d is out of bounds for transaction %s", utxo.Index, string(utxo.TxHash))
		}
		output := tx.Outputs[utxo.Index]
		// Check if the output belongs to the sender
		outputPubKey, err := crypto.PublicKeyFromBytes(output.PubKey)
		if err != nil {
			return nil, fmt.Errorf("get output public key err: %v", err)
		}
		if !outputPubKey.Equal(&txReq.Sender.PublicKey) {
			return nil, fmt.Errorf("sender do not own transaction's output: tx: %s", string(utxo.TxHash))
		}

		r, s, err := utxo.Sign(txReq.Sender)
		if err != nil {
			return nil, fmt.Errorf("sign UTXO:%s err: %s", string(utxo.TxHash), err.Error())
		}

		if !utxo.Verify(txReq.Sender.PublicKey, r, s) {
			return nil, fmt.Errorf("can not verify UTXO:%s err: not valid private key", string(utxo.TxHash))
		}
		newTx.AddInput(types.NewTxIn(utxo, senderPubKey, r, s, uint32(id)))
		balance.Value += output.Amount.Value
	}

	if balance.Value < txReq.Amount.Value {
		return nil, errors.New("insufficient balance")
	}

	newTx.AddOutput(types.NewTxOut(newTx.ID, txReq.Amount, crypto.PublicKeyToBytes(txReq.Receiver)))
	if balance.Value > txReq.Amount.Value {
		balance.Value -= txReq.Amount.Value
		// this output contains actual balance of the sender
		newTx.AddOutput(types.NewTxOut(newTx.ID, balance, crypto.PublicKeyToBytes(&txReq.Sender.PublicKey)))
	}
	newTx.ComputeHash()
	t.txPool.AddTx(newTx)

	return newTx, nil
}

func (t *Transactor) GetBalance(pubKey []byte) (*types.Amount, error) {
	pubKeyEcdsa, err := crypto.PublicKeyFromBytes(pubKey)
	if err != nil {
		return nil, fmt.Errorf("public key is not ECDSA")
	}
	utxos, err := t.utxoStore.Get(crypto.PublicKeyToBytes(pubKeyEcdsa))
	if err != nil {
		return nil, fmt.Errorf("error getting utxos : %v", err)
	}
	balance := &types.Amount{}
	for _, utxo := range utxos {
		tx, err := t.txStore.Get(utxo.TxHash)
		if err != nil {
			return nil, fmt.Errorf("get utxo tx hash err: %v", err)
		}
		if int(utxo.Index) >= len(tx.Outputs) {
			return nil, fmt.Errorf("UTXO index %d is out of bounds for transaction %s", utxo.Index, string(utxo.TxHash))
		}
		output := tx.Outputs[utxo.Index]
		balance.Value += output.Amount.Value
		// assume all outputs have the same unit
		balance.Unit = output.Amount.Unit
	}
	return balance, nil
}
