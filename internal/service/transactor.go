package service

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg/crypto"
	"local-chain/internal/types"
	"math/big"
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
	newTx := types.NewTransaction()
	balance, err := t.getBalance(
		txReq.Sender,
		func(utxo *types.UTXO, senderPubKey []byte, r, s *big.Int, id uint32) {
			newTx.AddInput(types.NewTxIn(utxo, senderPubKey, r, s, id))
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error getting balance : %v", err)
	}

	if balance.Value < txReq.Amount.Value {
		return nil, errors.New("insufficient balance")
	}

	newTx.AddOutput(types.NewTxOut(newTx.ID, txReq.Amount, crypto.PublicKeyToBytes(txReq.Receiver)))
	if balance.Value > txReq.Amount.Value {
		balance.Value -= txReq.Amount.Value
		// this output contains actual sender balance
		newTx.AddOutput(types.NewTxOut(newTx.ID, *balance, crypto.PublicKeyToBytes(&txReq.Sender.PublicKey)))
	}
	newTx.ComputeHash()
	t.txPool.AddTx(newTx)

	return newTx, nil
}

func (t *Transactor) GetBalance(req *types.BalanceRequest) (*types.Amount, error) {
	balance, err := t.getBalance(req.Sender, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting balance : %v", err)
	}
	return balance, nil
}

func (t *Transactor) getBalance(key *ecdsa.PrivateKey, f func(utxo *types.UTXO, senderPubKey []byte, r, s *big.Int, id uint32)) (*types.Amount, error) {
	pubKey := crypto.PublicKeyToBytes(&key.PublicKey)
	utxos, err := t.getUTXOs(pubKey)
	if err != nil {
		return nil, fmt.Errorf("error getting utxos : %v", err)
	}
	balance := types.NewAmount(0)
	for id, utxo := range utxos {
		tx, err := t.getTx(utxo.TxHash)
		if err != nil {
			return nil, fmt.Errorf("get utxo tx hash err: %v", err)
		}
		if int(utxo.Index) >= len(tx.Outputs) {
			return nil, fmt.Errorf("UTXO index %d is out of bounds for transaction %s", utxo.Index, string(utxo.TxHash))
		}
		output := tx.Outputs[utxo.Index]
		// Check if the output belongs to the key owner
		outputPubKey, err := crypto.PublicKeyFromBytes(output.PubKey)
		if err != nil {
			return nil, fmt.Errorf("get output public key err: %v", err)
		}
		if !outputPubKey.Equal(&key.PublicKey) {
			return nil, fmt.Errorf("sender do not own transaction's output: tx: %s", string(utxo.TxHash))
		}

		r, s, err := utxo.Sign(key)
		if err != nil {
			return nil, fmt.Errorf("sign UTXO:%s err: %s", string(utxo.TxHash), err.Error())
		}

		if !utxo.Verify(key.PublicKey, r, s) {
			return nil, fmt.Errorf("can not verify UTXO:%s err: not valid private key", string(utxo.TxHash))
		}
		if f != nil {
			f(utxo, pubKey, r, s, uint32(id))
		}
		balance.Value += output.Amount.Value
		// assume all outputs have the same unit
		balance.Unit = output.Amount.Unit
	}
	return balance, nil
}

// GetUTXOs gets unspent transaction outputs for public key
func (t *Transactor) getUTXOs(pubKey []byte) (types.UTXOs, error) {
	utxos, err := t.utxoStore.Get(pubKey)
	if err != nil {
		return nil, fmt.Errorf("error getting utxos : %v", err)
	}
	// we need to get utxos from the pool as well to avoid double spending
	// also we need to get the utxo with index > 0 only once (the rest are change utxos)
	utxosPool := t.txPool.GetUTXOs(pubKey)
	var rest *types.UTXO
	for _, utxo := range utxosPool {
		if utxo.Index == 0 {
			utxos = append(utxos, utxo)
			continue
		}
		if utxo.Index > 0 {
			rest = utxo
		}
	}
	if rest != nil {
		return types.UTXOs{rest}, nil
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
