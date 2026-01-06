package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"runtime/debug"

	"local-chain/internal/pkg/grpc/interceptors"

	grpc2 "local-chain/internal/adapters/inbound/grpc"
	"local-chain/internal/adapters/inbound/grpc/mapper"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg"
	"local-chain/internal/runners"
	"local-chain/internal/service"
	transport2 "local-chain/transport/gen/transport"

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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic recovered in main: %v\nstack: %s", r, string(debug.Stack()))
		}
	}()
	fmt.Println("raftBootstrap:", bootstrap)
	logger := slog.Default()
	ctx := pkg.ContextWithServerID(context.Background(), raft.ServerID(nodeID))
	cfg, err := NewConfig(serverID)
	if err != nil {
		log.Printf("error prepare configs: %v", err)
		return
	}

	dbFunc := func(subPath string) leveldbpkg.Database {
		path := fmt.Sprintf("%s/%s", dbDir, subPath)
		db, err := leveldb.OpenFile(path, nil)
		if err != nil {
			log.Printf("error open db file: %v", err)
			panic(err)
		}
		return db
	}
	store := leveldbpkg.New(dbFunc)
	defer func() {
		if err := store.Close(); err != nil {
			log.Printf("error closing store: %v", err)
		}
	}()
	txPool := inMem.NewTxPool()
	fsmStore := fsm.New(store, txPool)

	logStore, err := raftboltdb.NewBoltStore(logDb)
	if err != nil {
		log.Printf("error create logStore: %v", err)
		return
	}
	defer func() {
		if err := logStore.Close(); err != nil {
			log.Printf("error closing logStore: %v", err)
		}
	}()
	stableStore, err := raftboltdb.NewBoltStore(stableDb)
	if err != nil {
		log.Printf("error create stableStore: %v", err)
		return
	}
	defer func() {
		if err := stableStore.Close(); err != nil {
			log.Printf("error closing stableStore: %v", err)
		}
	}()

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

	user := service.NewUserService(store.User())
	um := mapper.NewUserMapper()
	superUser := initSuperUser(store.User())

	if bootstrap {
		configureBootstrap(r, store, superUser)
	}
	transactor := service.NewTransactor(store.Transaction(), store.Utxo(), txPool)
	tm := mapper.NewTransactionMapper()
	bm := mapper.NewBlockMapper()

	localChainManager := grpc2.NewLocalChain(serverID, r, tm, transactor, user, um, store.Blockchain(), store.Transaction(), bm)

	leaderRedirectInterceptor := interceptors.NewLeaderRedirectInterceptor(serverID, r)
	grpcRunner := runners.New(
		grpcAddr, func(s *grpc.Server) {
			transport2.RegisterLocalChainServer(s, localChainManager)
		},
		*logger,
		leaderRedirectInterceptor.UnaryInterceptor(),
	)

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
