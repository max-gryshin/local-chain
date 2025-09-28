package types

import "github.com/ethereum/go-ethereum/rlp"

type EnvelopeType string

const (
	EnvelopeTypeBlock       EnvelopeType = "block_type"
	EnvelopeTypeTransaction EnvelopeType = "transaction_type"
)

type Envelope struct {
	Type EnvelopeType
	Data []byte
}

func NewEnvelope(t EnvelopeType, data []byte) *Envelope {
	return &Envelope{
		Type: t,
		Data: data,
	}
}

func (e *Envelope) ToBytes() ([]byte, error) {
	return rlp.EncodeToBytes(e)
}

func EnvelopeFromBytes(data []byte) (*Envelope, error) {
	envelope := &Envelope{}
	err := rlp.DecodeBytes(data, envelope)

	return envelope, err
}
