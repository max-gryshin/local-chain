package main

import (
	"context"
	"flag"
	grpc2 "local-chain/internal/adapters/inbound/grpc"
	"local-chain/internal/adapters/inbound/grpc/mapper"
	"local-chain/internal/adapters/outbound/inMem"
	"local-chain/internal/pkg"
	"local-chain/internal/runners"
	"local-chain/internal/service"
	transport2 "local-chain/transport/gen/transport"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	transport "github.com/Jille/raft-grpc-transport"
	"github.com/gotidy/ptr"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	fsm "local-chain/internal/adapters/inbound/raft"
	leveldbpkg "local-chain/internal/adapters/outbound/leveldb"
)

const (
	dbPath = "./db"

	fsmDbPath      = dbPath + "/fsm"
	logDbPath      = dbPath + "/log"
	stableDbPath   = dbPath + "/stable"
	snapshotDbPath = dbPath + "/snapshot"
)

var (
	myAddr = flag.String("address", "localhost:8001", "TCP host+port for this node")
	raftId = flag.String("raft_id", "", "Node id used by Raft")

	raftBootstrap = flag.Bool("raft_bootstrap", false, "Whether to bootstrap the Raft cluster")
)

func main() {
	flag.Parse()

	logger := slog.Default()
	ctx := context.Background()

	db, err := leveldb.OpenFile(fsmDbPath, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() // nolint:errcheck

	store := leveldbpkg.New(db)
	fsmStore := fsm.New(store)

	logStore, err := raftboltdb.NewBoltStore(logDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer logStore.Close() // nolint:errcheck
	stableStore, err := raftboltdb.NewBoltStore(stableDbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer stableStore.Close() // nolint:errcheck

	snapshotStore, err := raft.NewFileSnapshotStore(snapshotDbPath, 3, os.Stderr)
	if err != nil {
		log.Fatal(err)
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
	_ = transport.New("localhost", []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	r, err := raft.NewRaft(
		raft.DefaultConfig(),
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
					ID:      raft.ServerID(*raftId),
					Address: raft.ServerAddress(*myAddr),
				},
			},
		})
		if err := configFuture.Error(); err != nil {
			log.Fatal(err)
		}
	}
	transactor := service.NewTransactor(store.Transaction())
	tm := mapper.NewTransactionMapper()
	txPool := inMem.NewTxPool()
	localChainManager := grpc2.NewLocalChainManager(r, raft.ServerID(*raftId), txPool, tm, transactor)

	grpcRunner := runners.New(9001, func(s *grpc.Server) {
		transport2.RegisterLocalChainManagerServer(s, localChainManager)
	}, *logger)

	blockchain := service.NewBlockchain(store.Blockchain(), store.Transaction())
	blockchainScheduler := runners.NewBlockchainScheduler(r, blockchain, txPool)

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
