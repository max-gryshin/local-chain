package main

import (
	"log"

	"local-chain/internal/pkg/crypto"

	"local-chain/internal/types"

	leveldbpkg "local-chain/internal/adapters/outbound/leveldb"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
)

func configureBootstrap(r *raft.Raft, store *leveldbpkg.Store) {
	configFuture := r.BootstrapCluster(raft.Configuration{
		Servers: []raft.Server{
			{
				ID:      serverID,
				Address: raft.ServerAddress(raftAddr),
			},
		},
	})
	genesisBlock := types.NewBlock(nil, []byte("genesis"))
	if err := store.Blockchain().Put(genesisBlock); err != nil {
		log.Fatal(err)
	}
	outputs := genesisOutputs()
	tx := genesisTx(genesisBlock, outputs)
	tx.ComputeHash()
	if err := store.Transaction().Put(tx); err != nil {
		log.Fatal(err)
	}
	for _, output := range outputs {
		pubKey, _ := crypto.PublicKeyFromBytes(output.PubKey)
		if err := store.Utxo().Put(crypto.PublicKeyToBytes(pubKey), types.NewUTXO(tx.GetHash(), 0)); err != nil {
			log.Fatal(err)
		}
	}
	if err := configFuture.Error(); err != nil {
		log.Fatal(err)
	}
}

func genesisTx(genesisBlock *types.Block, outputs []*types.TxOut) *types.Transaction {
	return &types.Transaction{
		BlockHash: genesisBlock.ComputeHash(),
		Outputs:   outputs,
		Hash:      []byte("genesis"),
	}
}

func genesisOutputs() []*types.TxOut {
	return []*types.TxOut{
		types.NewTxOut(
			uuid.MustParse("10252f31-151b-457d-b8de-e4a6f1552b62"),
			types.Amount{
				Value: 1_000_000_000_000,
				Unit:  100,
			},
			[]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEa/KaLpP9gikVe2ZXkp74RE+QmdDd
hJxRIN+5upGQgZyYFOqC7uwgXk0PS7GUNTl1aECoAKa2WEIWKL2PmTNZvg==
-----END PUBLIC KEY-----`)),
	}
}
