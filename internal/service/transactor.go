package service

import (
	"errors"
	"fmt"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg/crypto"
	"local-chain/internal/types"
	"slices"
)

//go:generate mockgen -source transactor.go -destination transactor_mock_test.go -package service_test -mock_names TransactionStore=MockTransactionStore,TxPool=MockTxPool

type TransactionStore interface {
	Get(txHash []byte) (*types.Transaction, error)
	Put(*types.Transaction) error
}

type UTXOStore interface {
	Get(pubKey []byte) ([]*types.UTXO, error)
	Put(pubKey []byte, utxos ...*types.UTXO) error
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
		tx           *types.Transaction
		balance      = types.Amount{}
		senderPubKey = crypto.PublicKeyToBytes(&txReq.Sender.PublicKey)
	)
	utxos, err := t.getUTXOs(senderPubKey)
	if err != nil {
		return nil, fmt.Errorf("error getting utxos : %v", err)
	}
	for id, utxo := range utxos {
		tx, err = t.getTx(utxo.TxHash)
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

func (t *Transactor) GetBalance(pubKeyBytes []byte) (*types.Amount, error) {
	pubKeyEcdsa, err := crypto.PublicKeyFromBytes(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("public key is not ECDSA")
	}
	pubKey := crypto.PublicKeyToBytes(pubKeyEcdsa)
	utxos, err := t.getUTXOs(pubKey)
	if err != nil {
		return nil, fmt.Errorf("error getting utxos : %v", err)
	}
	balance := &types.Amount{}
	for _, utxo := range utxos {
		tx, err := t.getTx(utxo.TxHash)
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

// GetUTXOs gets unspent transaction outputs for public key
func (t *Transactor) getUTXOs(pubKeyBytes []byte) (types.UTXOs, error) {
	utxos, err := t.utxoStore.Get(pubKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error getting utxos : %v", err)
	}
	// todo: add timestamp for utxos in pool to guaranty we get the oldest one utxo with index > 0
	utxosPool := t.txPool.GetUTXOs(pubKeyBytes)
	var rest *types.UTXO
	for _, utxo := range utxosPool {
		if utxo.Index == 0 {
			utxos = append(utxos, utxo)
		}
		if utxo.Index > 0 {
			if rest == nil {
				utxos = slices.DeleteFunc(utxos, func(u *types.UTXO) bool {
					return u.Index > 0
				})
			}
			rest = utxo
		}
	}
	if rest != nil {
		utxos = append(utxos, rest)
	}
	return utxos, nil
}

func (t *Transactor) getTx(hash []byte) (*types.Transaction, error) {
	var err error
	tx, ok := t.txPool.GetPool()[string(hash)]
	if !ok {
		tx, err = t.txStore.Get(hash)
		if err != nil {
			return nil, fmt.Errorf("error getting transaction : %v", err)
		}
	}

	return tx, nil
}
