package main

import (
	"context"
	"flag"
	"fmt"
	grpc2 "local-chain/internal/adapters/inbound/grpc"
	"local-chain/internal/adapters/inbound/grpc/mapper"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg"
	"local-chain/internal/runners"
	"local-chain/internal/service"
	"local-chain/internal/types"
	transport2 "local-chain/transport/gen/transport"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/gotidy/ptr"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"

	fsm "local-chain/internal/adapters/inbound/raft"
	leveldbpkg "local-chain/internal/adapters/outbound/leveldb"
)

const (
	dbPath = "./db"

	fsmDbPath  = dbPath + "/fsm"
	logDb      = dbPath + "/log.dat"
	stableDb   = dbPath + "/stable.dat"
	snapshotDb = dbPath
)

var (
	myAddr        = flag.String("address", "127.0.0.1:8001", "TCP host+port for this node")
	raftId        = flag.String("raft_id", "10252f31-151b-457d-b8de-e4a6f1552b62", "Node id used by Raft")
	serverID      = raft.ServerID(ptr.ToString(raftId))
	raftBootstrap = flag.Bool("raft_bootstrap", true, "Whether to bootstrap the Raft cluster")
)

func main() {
	flag.Parse()

	fmt.Println("raftBootstrap:", *raftBootstrap)
	logger := slog.Default()
	ctx := pkg.ContextWithServerID(context.Background(), serverID)

	//ex, err := os.Executable()
	//if err != nil {
	//	panic(err)
	//}
	//exPath := filepath.Dir(ex)
	//fmt.Println(exPath)

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		log.Printf("error open db file: %v", err)
		return
	}
	defer db.Close() // nolint:errcheck

	store := leveldbpkg.New(db)
	fsmStore := fsm.New(store)

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
		ptr.ToString(myAddr),
		&net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 8001},
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
	if *raftBootstrap {
		configFuture := r.BootstrapCluster(raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      serverID,
					Address: raft.ServerAddress(ptr.ToString(myAddr)),
				},
			},
		})
		// genesis block
		if err = store.Blockchain().Put(&types.Block{
			Timestamp: 0,
			Hash:      []byte("genesis"),
		}); err != nil {
			log.Fatal(err)
		}
		if err = configFuture.Error(); err != nil {
			log.Fatal(err)
		}
	}
	transactor := service.NewTransactor(store.Transaction())
	tm := mapper.NewTransactionMapper()
	txPool := inMem.NewTxPool()
	localChainManager := grpc2.NewLocalChainManager(r, txPool, tm, transactor)

	grpcRunner := runners.New(9001, func(s *grpc.Server) {
		transport2.RegisterLocalChainManagerServer(s, localChainManager)
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
