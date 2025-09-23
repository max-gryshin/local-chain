package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"local-chain/internal/pkg"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcPkg "local-chain/transport/gen/transport"

	"local-chain/internal/types"
)

type RaftAPI interface {
	AddNonvoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	RemoveServer(id raft.ServerID, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	AddVoter(id raft.ServerID, address raft.ServerAddress, prevIndex uint64, timeout time.Duration) raft.IndexFuture
	LeaderWithID() (raft.ServerAddress, raft.ServerID)
	Apply(cmd []byte, timeout time.Duration) raft.ApplyFuture
}

type txPool interface {
	AddTx(tx *types.Transaction)
}

type Transactor interface {
	CreateTx(txReq *types.TransactionRequest) (*types.Transaction, error)
}

type transactionMapper interface {
	RpcToTransaction(req *grpcPkg.AddTransactionRequest) (*types.TransactionRequest, error)
}

type LocalChainServer struct {
	raftAPI RaftAPI
	txPool  txPool
	tm      transactionMapper
	grpcPkg.UnimplementedLocalChainServer
	transactor Transactor
}

func NewLocalChain(
	raftAPI RaftAPI,
	txPool txPool,
	tm transactionMapper,
	transactor Transactor,
) *LocalChainServer {
	return &LocalChainServer{
		raftAPI:    raftAPI,
		txPool:     txPool,
		tm:         tm,
		transactor: transactor,
	}
}

func (s *LocalChainServer) AddPeer(ctx context.Context, req *grpcPkg.AddPeerRequest) (*grpcPkg.AddPeerResponse, error) {
	if req.GetId() == "" || req.GetAddress() == "" {
		return &grpcPkg.AddPeerResponse{Success: false}, errors.New("peer ID and address must be provided")
	}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != pkg.ServerIDFromContext(ctx) {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return &grpcPkg.AddPeerResponse{Success: false}, errors.New("failed to connect to leader")
		}
		return client.AddPeer(ctx, req)
	}
	future := s.raftAPI.AddNonvoter(raft.ServerID(req.GetId()), raft.ServerAddress(req.GetAddress()), 0, 0)
	if err := future.Error(); err != nil {
		return &grpcPkg.AddPeerResponse{Success: false}, err
	}

	return &grpcPkg.AddPeerResponse{Success: true}, nil
}

func (s *LocalChainServer) RemovePeer(ctx context.Context, req *grpcPkg.RemovePeerRequest) (*grpcPkg.RemovePeerResponse, error) {
	if req.GetId() == "" || req.GetAddress() == "" {
		return &grpcPkg.RemovePeerResponse{Success: false}, errors.New("peer ID and address must be provided")
	}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != pkg.ServerIDFromContext(ctx) {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return &grpcPkg.RemovePeerResponse{Success: false}, errors.New("failed to connect to leader")
		}
		return client.RemovePeer(ctx, req)
	}
	future := s.raftAPI.RemoveServer(raft.ServerID(req.GetId()), 0, 0)
	if err := future.Error(); err != nil {
		return &grpcPkg.RemovePeerResponse{Success: false}, err
	}
	return &grpcPkg.RemovePeerResponse{Success: true}, nil
}

func (s *LocalChainServer) AddVoter(ctx context.Context, req *grpcPkg.AddVoterRequest) (*grpcPkg.AddVoterResponse, error) {
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != pkg.ServerIDFromContext(ctx) {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return &grpcPkg.AddVoterResponse{Success: false}, errors.New("failed to connect to leader")
		}
		return client.AddVoter(ctx, req)
	}
	if req.GetId() == "" || req.GetAddress() == "" {
		return &grpcPkg.AddVoterResponse{Success: false}, errors.New("peer ID and address must be provided")
	}
	future := s.raftAPI.AddVoter(raft.ServerID(req.GetId()), raft.ServerAddress(req.GetAddress()), 0, 0)
	if err := future.Error(); err != nil {
		return &grpcPkg.AddVoterResponse{Success: false}, err
	}
	return &grpcPkg.AddVoterResponse{Success: true}, nil
}

func (s *LocalChainServer) AddTransaction(ctx context.Context, req *grpcPkg.AddTransactionRequest) (*grpcPkg.AddTransactionResponse, error) {
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != pkg.ServerIDFromContext(ctx) {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return &grpcPkg.AddTransactionResponse{Success: false}, errors.New("failed to connect to leader")
		}
		return client.AddTransaction(ctx, req)
	}
	txReq, err := s.tm.RpcToTransaction(req)
	if err != nil {
		return &grpcPkg.AddTransactionResponse{Success: false}, fmt.Errorf("failed to marshal add transaction request: %w", err)
	}
	tx, err := s.transactor.CreateTx(txReq)
	if err != nil {
		return nil, fmt.Errorf("transactor.CreateTx: %w", err)
	}
	s.txPool.AddTx(tx)
	// todo:
	// validate req transaction* can skip it to speed up the implementation
	// transactor.CreateTx
	// txPool.AddTx
	// implement scheduler which will start a process of creating a new block
	// scheduler calls blockchain.CreateBlock with the pool of transactions
	// and broadcast it to all peers by raft.Apply

	return &grpcPkg.AddTransactionResponse{Success: true}, nil
}

func (s *LocalChainServer) leaderClient(leaderAddr string) (grpcPkg.LocalChainClient, error) {
	conn, err := grpc.NewClient(leaderAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer func(conn *grpc.ClientConn) {
		err = conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	return grpcPkg.NewLocalChainClient(conn), nil
}
