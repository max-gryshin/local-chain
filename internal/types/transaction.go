package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

const CurrencyUnit = 100000000

type Transaction struct {
	ID        uuid.UUID
	Timestamp uint64
	nLockTime uint32
	BlockHash []byte

	Salt [16]byte
	Hash []byte

	Inputs  []*TxIn
	Outputs []*TxOut
}

func NewTransaction() *Transaction {
	tx := &Transaction{
		ID:        uuid.New(),
		Timestamp: uint64(time.Now().UnixNano()),
		Salt:      [16]byte(uuid.New()),
	}
	return tx
}

func (tx *Transaction) WithInputs(inputs ...*TxIn) *Transaction {
	tx.Inputs = append(tx.Inputs, inputs...)
	return tx
}

func (tx *Transaction) WithInput(input *TxIn) *Transaction {
	tx.Inputs = append(tx.Inputs, input)

	return tx
}

func (tx *Transaction) WithOutput(amount *Amount, key *ecdsa.PublicKey) *Transaction {
	tx.Outputs = append(tx.Outputs, &TxOut{
		TxID:   tx.ID,
		Amount: *amount,
		PubKey: ToHashString(key),
	})
	return tx
}

func (tx *Transaction) WithOutputs(outputs ...*TxOut) *Transaction {
	tx.Outputs = append(tx.Outputs, outputs...)

	return tx
}

func (tx *Transaction) ComputeHash() {
	data := make([]byte, 0, 256)
	data = append(data, tx.ID[:]...)
	timestamp := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(tx.Timestamp))
	data = append(data, timestamp...)
	nLockTime := make([]byte, 8)
	binary.LittleEndian.PutUint64(timestamp, uint64(tx.nLockTime))
	data = append(data, nLockTime...)
	for _, out := range tx.Outputs {
		data = append(data, out.TxID[:]...)
		data = append(data, []byte(out.PubKey)...)
		data = append(data, out.Amount.ToBytes()...)
	}
	hash := sha512.New()
	hash.Write(data)
	tx.Hash = hash.Sum(nil)
}

func (tx *Transaction) GetHash() []byte {
	if tx.Hash == nil {
		tx.ComputeHash()
	}
	return tx.Hash
}

type TxIn struct {
	Prev       *UTXO
	PubKey     *ecdsa.PublicKey
	SignatureR *big.Int
	SignatureS *big.Int
	NSequence  uint32
}

type TxOut struct {
	TxID   uuid.UUID
	Amount Amount
	PubKey string
}

func ToHashString(pubKey *ecdsa.PublicKey) string {
	if pubKey == nil {
		return ""
	}
	pubKeyBytes := append(pubKey.X.Bytes(), pubKey.Y.Bytes()...)
	hash := sha512.Sum512(pubKeyBytes)

	return string(hash[:])
}

type UTXO struct {
	TxHash []byte
	// Index of the output in the transaction
	Index uint32
}

func (u *UTXO) Sign(key *ecdsa.PrivateKey) (*big.Int, *big.Int, error) {
	r, s, err := ecdsa.Sign(rand.Reader, key, u.TxHash)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to sign UTXO %s: %s", u.TxHash, err.Error())
	}

	return r, s, nil
}

func (u *UTXO) Verify(pubKey ecdsa.PublicKey, r, s *big.Int) bool {
	return ecdsa.Verify(&pubKey, u.TxHash, r, s)
}

type Amount struct {
	Value uint64
	Unit  uint32
}

func (a *Amount) ToBytes() []byte {
	amount := make([]byte, 8)
	binary.LittleEndian.PutUint64(amount, uint64(a.Value))
	unit := make([]byte, 8)
	binary.LittleEndian.PutUint32(unit, uint32(a.Unit))

	return append(amount, unit...)
}

func NewAmount(value uint64) *Amount {
	return &Amount{
		Value: value,
		Unit:  CurrencyUnit,
	}
}

type TransactionRequest struct {
	Sender   *ecdsa.PrivateKey
	Receiver *ecdsa.PublicKey
	Amount   Amount
	Utxos    []*UTXO
}
