package main

import (
	"context"
	"fmt"
	"local-chain/internal/pkg/crypto"
	"log"
	"log/slog"
	"net"
	"net/netip"
	"os"
	"time"

	grpc2 "local-chain/internal/adapters/inbound/grpc"
	"local-chain/internal/adapters/inbound/grpc/mapper"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg"
	"local-chain/internal/runners"
	"local-chain/internal/service"
	transport2 "local-chain/transport/gen/transport"

	"local-chain/internal/types"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"

	fsm "local-chain/internal/adapters/inbound/raft"
	leveldbpkg "local-chain/internal/adapters/outbound/leveldb"
)

var (
	nodeID   = os.Getenv("NODE_ID")
	raftAddr = os.Getenv("RAFT_ADDR")
	grpcAddr = os.Getenv("GRPC_ADDR")
	dbDir    = os.Getenv("DATA_DIR")

	logDb      = dbDir + "/log.dat"
	stableDb   = dbDir + "/stable.dat"
	snapshotDb = dbDir

	bootstrap = os.Getenv("BOOTSTRAP") == "true"
	serverID  = raft.ServerID(nodeID)
)

func main() {
	fmt.Println("raftBootstrap:", bootstrap)
	logger := slog.Default()
	ctx := pkg.ContextWithServerID(context.Background(), serverID)

	//ex, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//exPath := filepath.Dir(ex)
	//fmt.Println(exPath)

	db, err := leveldb.OpenFile(dbDir, nil)
	if err != nil {
		log.Printf("error open db file: %v", err)
		return
	}
	defer db.Close() // nolint:errcheck

	store := leveldbpkg.New(db)
	txPool := inMem.NewTxPool()
	fsmStore := fsm.New(store, txPool)

	logStore, err := raftboltdb.NewBoltStore(logDb)
	if err != nil {
		log.Printf("error create logStore: %v", err)
		return
	}
	defer logStore.Close() // nolint:errcheck
	stableStore, err := raftboltdb.NewBoltStore(stableDb)
	if err != nil {
		log.Printf("error create stableStore: %v", err)
		return
	}
	defer stableStore.Close() // nolint:errcheck

	snapshotStore, err := raft.NewFileSnapshotStore(snapshotDb, 3, os.Stderr)
	if err != nil {
		log.Printf("error create snapshotStore: %v", err)
		return
	}

	raftAddrPort, err := netip.ParseAddrPort(raftAddr)
	if err != nil {
		log.Printf("error parse raft addr: %v", err)
		return
	}

	tr, err := raft.NewTCPTransport(
		raftAddr,
		net.TCPAddrFromAddrPort(raftAddrPort),
		3,
		10*time.Second,
		os.Stderr,
	)
	if err != nil {
		log.Fatal(err)
	}
	raftConfig := &raft.Config{
		ProtocolVersion:    raft.ProtocolVersionMax,
		HeartbeatTimeout:   1000 * time.Millisecond,
		ElectionTimeout:    1000 * time.Millisecond,
		CommitTimeout:      50 * time.Millisecond,
		MaxAppendEntries:   64,
		ShutdownOnRemove:   true,
		TrailingLogs:       10240,
		SnapshotInterval:   120 * time.Second,
		SnapshotThreshold:  8192,
		LeaderLeaseTimeout: 500 * time.Millisecond,
		LogLevel:           "DEBUG",
		LocalID:            serverID,
	}
	r, err := raft.NewRaft(
		raftConfig,
		fsmStore,
		logStore,
		stableStore,
		snapshotStore,
		tr,
	)
	if err != nil {
		log.Fatal(err)
	}
	if bootstrap {
		configFuture := r.BootstrapCluster(raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      serverID,
					Address: raft.ServerAddress(raftAddr),
				},
			},
		})
		genesisBlock := types.NewBlock(nil, []byte("genesis"))
		if err = store.Blockchain().Put(genesisBlock); err != nil {
			log.Fatal(err)
		}
		outputs := genesisOutputs()
		tx := genesisTx(genesisBlock, outputs)
		tx.ComputeHash()
		if err = store.Transaction().Put(tx); err != nil {
			log.Fatal(err)
		}
		for _, output := range outputs {
			utxos := make([]*types.UTXO, 0)
			pubKey, err := crypto.PublicKeyFromBytes(output.PubKey)
			if err != nil {
				log.Fatal(err)
			}
			utxos = append(utxos, &types.UTXO{TxHash: tx.GetHash(), Index: 0})
			if err = store.Utxo().Put(crypto.PublicKeyToBytes(pubKey), utxos); err != nil {
				log.Fatal(err)
			}
		}
		if err = configFuture.Error(); err != nil {
			log.Fatal(err)
		}
	}
	transactor := service.NewTransactor(store.Transaction(), store.Utxo())
	tm := mapper.NewTransactionMapper()
	localChainManager := grpc2.NewLocalChain(serverID, r, txPool, tm, transactor)

	grpcRunner := runners.New(grpcAddr, func(s *grpc.Server) {
		transport2.RegisterLocalChainServer(s, localChainManager)
	}, *logger)

	blockchain := service.NewBlockchain(r, store.Blockchain(), store.Transaction(), txPool)
	blockchainScheduler := runners.NewBlockchainScheduler(blockchain)

	runnable := []pkg.Runner{
		grpcRunner,
		blockchainScheduler,
	}

	firstError := pkg.Run(ctx, logger, runnable...)

	if firstError != nil {
		logger.Error("runner finished with an error", slog.Any("error", firstError))
	} else {
		logger.Info("runner finished successfully")
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
				Value: 100000000,
				Unit:  100,
			},
			[]byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEa/KaLpP9gikVe2ZXkp74RE+QmdDd
hJxRIN+5upGQgZyYFOqC7uwgXk0PS7GUNTl1aECoAKa2WEIWKL2PmTNZvg==
-----END PUBLIC KEY-----`)),
	}
}
