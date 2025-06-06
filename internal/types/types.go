package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const CurrencyUnit = 100000000

type Transaction struct {
	ID        uuid.UUID
	Timestamp int64
	nLockTime uint32 // ??? зачем?

	Salt [16]byte
	Hash []byte

	Inputs  []*TxIn
	Outputs []*TxOut
}

func NewTransaction() *Transaction {
	tx := &Transaction{
		ID:        uuid.New(),
		Timestamp: time.Now().UnixNano(),
		Salt:      [16]byte(uuid.New()),
	}
	tx.Hash = tx.computeHash()
	return tx
}

func (tx *Transaction) WithInput(inputs ...*TxIn) {
	tx.Inputs = inputs
}

func (tx *Transaction) WithOutput(outputs ...*TxOut) {
	tx.Outputs = outputs
}

func (tx *Transaction) computeHash() []byte {
	data := tx.ID.String() + strconv.Itoa(int(tx.Timestamp)) + strconv.Itoa(int(tx.nLockTime)) + string(tx.Salt[:])

	hash := sha512.New()
	hash.Write([]byte(data))

	return hash.Sum(nil)
}

type TxIn struct {
	Prev       *UTXO
	PubKey     ecdsa.PublicKey
	SignatureR *big.Int
	SignatureS *big.Int
	NSequence  int32
}

type TxOut struct {
	TxID   uuid.UUID // ID транзакции, к которой относится этот выход
	Amount Amount
	PubKey ecdsa.PublicKey
}

type UTXO struct {
	TxHash []byte
	// Index of the output in the transaction
	Index int32
}

func (u *UTXO) Sign(key *ecdsa.PrivateKey) (*big.Int, *big.Int, error) {
	r, s, err := ecdsa.Sign(rand.Reader, key, u.TxHash)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("failed to sign UTXO %s: %s", u.TxHash, err.Error()))
	}

	return r, s, nil
}

func (u *UTXO) Verify(pubKey *ecdsa.PublicKey, r, s *big.Int) bool {
	return ecdsa.Verify(pubKey, u.TxHash, r, s)
}

type Amount struct {
	Value int64
	Unit  int32
}
