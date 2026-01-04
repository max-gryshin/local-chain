package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

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
	State() raft.RaftState
}

type Transactor interface {
	CreateTx(txReq *types.TransactionRequest) (*types.Transaction, error)
	GetBalance(req *types.BalanceRequest) (*types.Amount, error)
}

type transactionMapper interface {
	RpcToTransaction(req *grpcPkg.AddTransactionRequest) (*types.TransactionRequest, error)
	RpcToBalanceRequest(req *grpcPkg.GetBalanceRequest) (*types.BalanceRequest, error)
}

type User interface {
	GetAllUsers() ([]*types.User, error)
	GetUser(username string) (*types.User, error)
	AddUser(user *types.User) error
}

type UserMapper interface {
	RpcToUser(req *grpcPkg.AddUserRequest) *types.User
	UserToRpc(user *types.User) *grpcPkg.User
}

type LocalChainServer struct {
	serverID raft.ServerID
	raftAPI  RaftAPI
	tm       transactionMapper
	grpcPkg.UnimplementedLocalChainServer
	transactor Transactor
	user       User
	userMapper UserMapper
}

func NewLocalChain(
	serverID raft.ServerID,
	raftAPI RaftAPI,
	tm transactionMapper,
	transactor Transactor,
	user User,
	userMapper UserMapper,
) *LocalChainServer {
	return &LocalChainServer{
		serverID:   serverID,
		raftAPI:    raftAPI,
		tm:         tm,
		transactor: transactor,
		user:       user,
		userMapper: userMapper,
	}
}

func (s *LocalChainServer) AddPeer(ctx context.Context, req *grpcPkg.AddPeerRequest) (*grpcPkg.AddPeerResponse, error) {
	if req.GetId() == "" || req.GetAddress() == "" {
		return &grpcPkg.AddPeerResponse{Success: false}, errors.New("peer ID and address must be provided")
	}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return &grpcPkg.AddPeerResponse{Success: false}, errors.New("failed to connect to leader")
		}
		return client.AddPeer(ctx, req)
	}
	future := s.raftAPI.AddNonvoter(raft.ServerID(req.GetId()), raft.ServerAddress(req.GetAddress()), 0, 0)
	if err := future.Error(); err != nil {
		fmt.Println("AddNonvoter error", future.Error())
		return &grpcPkg.AddPeerResponse{Success: false}, err
	}

	return &grpcPkg.AddPeerResponse{Success: true}, nil
}

func (s *LocalChainServer) RemovePeer(ctx context.Context, req *grpcPkg.RemovePeerRequest) (*grpcPkg.RemovePeerResponse, error) {
	if req.GetId() == "" || req.GetAddress() == "" {
		return &grpcPkg.RemovePeerResponse{Success: false}, errors.New("peer ID and address must be provided")
	}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
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
	if req.GetId() == "" || req.GetAddress() == "" {
		return &grpcPkg.AddVoterResponse{Success: false}, errors.New("peer ID and address must be provided")
	}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return &grpcPkg.AddVoterResponse{Success: false}, errors.New("failed to connect to leader")
		}
		return client.AddVoter(ctx, req)
	}
	future := s.raftAPI.AddVoter(raft.ServerID(req.GetId()), raft.ServerAddress(req.GetAddress()), 0, 10*time.Second)
	if err := future.Error(); err != nil {
		return &grpcPkg.AddVoterResponse{Success: false}, err
	}
	return &grpcPkg.AddVoterResponse{Success: true}, nil
}

func (s *LocalChainServer) AddTransaction(ctx context.Context, req *grpcPkg.AddTransactionRequest) (*grpcPkg.AddTransactionResponse, error) {
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
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
	if _, err = s.transactor.CreateTx(txReq); err != nil {
		return &grpcPkg.AddTransactionResponse{Success: false}, fmt.Errorf("transactor.CreateTx: %w", err)
	}
	// todo: validate req transaction* can skip it to speed up the implementation

	return &grpcPkg.AddTransactionResponse{Success: true}, nil
}

func (s *LocalChainServer) GetBalance(ctx context.Context, req *grpcPkg.GetBalanceRequest) (*grpcPkg.GetBalanceResponse, error) {
	resp := &grpcPkg.GetBalanceResponse{Amount: &grpcPkg.Amount{}}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return resp, errors.New("failed to connect to leader")
		}
		return client.GetBalance(ctx, req)
	}
	balanceReq, err := s.tm.RpcToBalanceRequest(req)
	if err != nil {
		return &grpcPkg.GetBalanceResponse{}, fmt.Errorf("failed to marshal get balance request: %w", err)
	}
	amount, err := s.transactor.GetBalance(balanceReq)
	if err != nil {
		return resp, fmt.Errorf("transactor.GetBalance: %w", err)
	}
	resp.Amount = &grpcPkg.Amount{Value: amount.Value, Unit: amount.Unit}

	return resp, err
}

func (s *LocalChainServer) AddUser(ctx context.Context, req *grpcPkg.AddUserRequest) (*grpcPkg.AddUserResponse, error) {
	if req.GetUser().GetUsername() == "" || len(req.GetUser().GetPrivateKey()) == 0 || len(req.GetUser().GetPublicKey()) == 0 {
		return &grpcPkg.AddUserResponse{Success: false}, errors.New("username, private key and public key must be provided")
	}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return &grpcPkg.AddUserResponse{Success: false}, errors.New("failed to connect to leader")
		}
		return client.AddUser(ctx, req)
	}
	if err := s.user.AddUser(s.userMapper.RpcToUser(req)); err != nil {
		return &grpcPkg.AddUserResponse{Success: false}, fmt.Errorf("user.AddUser: %w", err)
	}
	return &grpcPkg.AddUserResponse{Success: true}, nil
}

func (s *LocalChainServer) GetUser(ctx context.Context, req *grpcPkg.GetUserRequest) (*grpcPkg.GetUserResponse, error) {
	if req.GetUsername() == "" {
		return nil, errors.New("username must be provided")
	}
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return nil, errors.New("failed to connect to leader")
		}
		return client.GetUser(ctx, req)
	}
	user, err := s.user.GetUser(req.GetUsername())
	if err != nil {
		return nil, fmt.Errorf("user.GetUser: %w", err)
	}
	return &grpcPkg.GetUserResponse{
		User: s.userMapper.UserToRpc(user),
	}, nil
}

func (s *LocalChainServer) ListUsers(ctx context.Context, req *grpcPkg.ListUsersRequest) (*grpcPkg.ListUsersResponse, error) {
	leaderServer, leaderID := s.raftAPI.LeaderWithID()
	if leaderID != s.serverID {
		client, err := s.leaderClient(string(leaderServer))
		if err != nil {
			return nil, errors.New("failed to connect to leader")
		}
		return client.ListUsers(ctx, req)
	}
	users, err := s.user.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("user.GetAllUsers: %w", err)
	}
	rpcUsers := make([]*grpcPkg.User, 0, len(users))
	for _, user := range users {
		rpcUsers = append(rpcUsers, s.userMapper.UserToRpc(user))
	}
	return &grpcPkg.ListUsersResponse{Users: rpcUsers}, nil
}

// leaderClient creates a gRPC client connected to the current leader.
// todo: keep the connection open instead of creating a new one each time? if yes - move to main func to correctly close the connection on shutdown
func (s *LocalChainServer) leaderClient(leaderAddr string) (grpcPkg.LocalChainClient, error) {
	host, _, err := net.SplitHostPort(leaderAddr)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(net.JoinHostPort(host, "9001"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return grpcPkg.NewLocalChainClient(conn), nil
}
