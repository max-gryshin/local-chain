package main

import (
	"context"
	"fmt"
	grpc2 "local-chain/internal/adapters/inbound/grpc"
	"local-chain/internal/adapters/inbound/grpc/mapper"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg"
	"local-chain/internal/runners"
	"local-chain/internal/service"
	transport2 "local-chain/transport/gen/transport"
	"log"
	"log/slog"
	"os"

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
	ctx := pkg.ContextWithServerID(context.Background(), raft.ServerID(raftAddr))
	cfg, err := NewConfig(serverID)
	if err != nil {
		log.Printf("error prepare configs: %v", err)
		return
	}

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
	tr, err := raft.NewTCPTransport(
		raftAddr,
		cfg.TCPTransport.Address,
		cfg.TCPTransport.MaxPool,
		cfg.TCPTransport.Timeout,
		cfg.TCPTransport.LogOutput,
	)
	if err != nil {
		log.Fatal(err)
	}
	r, err := raft.NewRaft(
		cfg.Raft,
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
		configureBootstrap(r, store)
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
