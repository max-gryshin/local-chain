package raft

import (
	"encoding/json"
	"io"

	"local-chain/internal/types"

	"github.com/hashicorp/raft"

	"local-chain/internal/adapters/outbound/leveldb"
)

type Fsm struct {
	store *leveldb.Store
}

func New(store *leveldb.Store) *Fsm {
	return &Fsm{
		store: store,
	}
}

func (f *Fsm) Apply(log *raft.Log) interface{} {
	switch log.Type {
	case raft.LogCommand:
		// todo:
		// switch on types of messages (block, transaction, utxo, etc.)
		// deserialize log into msg
		// put msg into appropriate store
		return nil
	default:
		return nil
	}
}

func (f *Fsm) Snapshot() (raft.FSMSnapshot, error) {
	// todo: what is the snapshot?
	blocks, err := f.store.Blockchain().Get()
	if err != nil {
		return nil, err
	}
	return &FsmSnapshot{state: blocks}, nil
}

func (f *Fsm) Restore(rc io.ReadCloser) error {
	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &f.store)
}

type FsmSnapshot struct {
	state []*types.Block
}

func (s *FsmSnapshot) Persist(sink raft.SnapshotSink) error {
	// save the snapshot
	data, err := json.Marshal(s.state)
	if err != nil {
		return err
	}
	_, err = sink.Write(data)
	if err != nil {
		sink.Cancel() // nolint:errcheck
		return err
	}
	return sink.Close()
}

func (s *FsmSnapshot) Release() {
	// release resources if needed
}
