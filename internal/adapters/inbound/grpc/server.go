package grpc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"local-chain/internal/adapters/outbound/leveldb"

	grpcPkg "local-chain/transport/gen/transport"

	"local-chain/internal/types"

	"github.com/google/uuid"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/types/known/emptypb"
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
	TransactionToRpc(tx *types.Transaction) *grpcPkg.Transaction
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

type BlockchainStore interface {
	GetKeys() ([]uint64, error)
	GetByTimestamp(t uint64) (*types.Block, error)
}

type TransactionStore interface {
	Get(id uuid.UUID) (*types.Transaction, error)
}

type BlockMapper interface {
	BlockToRpc(block *types.Block) *grpcPkg.Block
	BlocksToRpc(blocks types.Blocks) []*grpcPkg.Block
}

type LocalChainServer struct {
	serverID raft.ServerID
	raftAPI  RaftAPI
	tm       transactionMapper
	grpcPkg.UnimplementedLocalChainServer
	transactor       Transactor
	user             User
	userMapper       UserMapper
	blockchainStore  BlockchainStore
	transactionStore TransactionStore
	blockMapper      BlockMapper
}

func NewLocalChain(
	serverID raft.ServerID,
	raftAPI RaftAPI,
	tm transactionMapper,
	transactor Transactor,
	user User,
	userMapper UserMapper,
	blockchainStore BlockchainStore,
	transactionStore TransactionStore,
	blockMapper BlockMapper,
) *LocalChainServer {
	return &LocalChainServer{
		serverID:         serverID,
		raftAPI:          raftAPI,
		tm:               tm,
		transactor:       transactor,
		user:             user,
		userMapper:       userMapper,
		blockchainStore:  blockchainStore,
		transactionStore: transactionStore,
		blockMapper:      blockMapper,
	}
}

func (s *LocalChainServer) AddPeer(ctx context.Context, req *grpcPkg.AddPeerRequest) (*grpcPkg.AddPeerResponse, error) {
	if req.GetId() == "" || req.GetAddress() == "" {
		return &grpcPkg.AddPeerResponse{Success: false}, errors.New("peer ID and address must be provided")
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
	future := s.raftAPI.AddVoter(raft.ServerID(req.GetId()), raft.ServerAddress(req.GetAddress()), 0, 10*time.Second)
	if err := future.Error(); err != nil {
		return &grpcPkg.AddVoterResponse{Success: false}, err
	}
	return &grpcPkg.AddVoterResponse{Success: true}, nil
}

func (s *LocalChainServer) AddTransaction(ctx context.Context, req *grpcPkg.AddTransactionRequest) (*grpcPkg.AddTransactionResponse, error) {
	txReq, err := s.tm.RpcToTransaction(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal add transaction request: %w", err)
	}
	tx, err := s.transactor.CreateTx(txReq)
	if err != nil {
		return nil, fmt.Errorf("transactor.CreateTx: %w", err)
	}
	// todo: validate req transaction* can skip it to speed up the implementation

	return &grpcPkg.AddTransactionResponse{Transaction: s.tm.TransactionToRpc(tx)}, nil
}

func (s *LocalChainServer) GetBalance(ctx context.Context, req *grpcPkg.GetBalanceRequest) (*grpcPkg.GetBalanceResponse, error) {
	resp := &grpcPkg.GetBalanceResponse{Amount: &grpcPkg.Amount{}}
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
	if err := s.user.AddUser(s.userMapper.RpcToUser(req)); err != nil {
		return &grpcPkg.AddUserResponse{Success: false}, fmt.Errorf("user.AddUser: %w", err)
	}
	return &grpcPkg.AddUserResponse{Success: true}, nil
}

func (s *LocalChainServer) GetUser(ctx context.Context, req *grpcPkg.GetUserRequest) (*grpcPkg.GetUserResponse, error) {
	if req.GetUsername() == "" {
		return nil, errors.New("username must be provided")
	}
	user, err := s.user.GetUser(req.GetUsername())
	if err != nil && !errors.Is(err, leveldb.ErrNotFound) {
		return nil, fmt.Errorf("user.GetUser: %w", err)
	}
	if user == nil {
		return nil, nil
	}
	return &grpcPkg.GetUserResponse{
		User: s.userMapper.UserToRpc(user),
	}, nil
}

func (s *LocalChainServer) ListUsers(ctx context.Context, req *emptypb.Empty) (*grpcPkg.ListUsersResponse, error) {
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

func (s *LocalChainServer) GetBlockKeys(ctx context.Context, req *emptypb.Empty) (*grpcPkg.GetBlockKeysResponse, error) {
	keys, err := s.blockchainStore.GetKeys()
	if err != nil {
		return nil, fmt.Errorf("blockchainStore.GetKeys: %w", err)
	}
	return &grpcPkg.GetBlockKeysResponse{Timestamp: keys}, nil
}

func (s *LocalChainServer) GetBlock(ctx context.Context, req *grpcPkg.GetBlockRequest) (*grpcPkg.GetBlockResponse, error) {
	if req.GetTimestamp() == 0 {
		return nil, errors.New("timestamp must be provided")
	}
	block, err := s.blockchainStore.GetByTimestamp(req.GetTimestamp())
	if err != nil {
		return nil, fmt.Errorf("blockchainStore.GetByTimestamp: %w", err)
	}
	if block == nil {
		return &grpcPkg.GetBlockResponse{Blocks: []*grpcPkg.Block{}}, nil
	}
	return &grpcPkg.GetBlockResponse{Blocks: []*grpcPkg.Block{s.blockMapper.BlockToRpc(block)}}, nil
}

func (s *LocalChainServer) GetTransaction(ctx context.Context, req *grpcPkg.GetTransactionRequest) (*grpcPkg.GetTransactionResponse, error) {
	if len(req.GetId()) == 0 {
		return nil, errors.New("transaction id must be provided")
	}
	txID, err := uuid.ParseBytes(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("invalid transaction id format: %w", err)
	}
	tx, err := s.transactionStore.Get(txID)
	if err != nil {
		return nil, fmt.Errorf("transactionStore.Get: %w", err)
	}
	if tx == nil {
		return nil, fmt.Errorf("transaction in not found or not added to store yet")
	}
	return &grpcPkg.GetTransactionResponse{Transaction: s.tm.TransactionToRpc(tx)}, nil
}

func (s *LocalChainServer) VerifyTransaction(ctx context.Context, req *grpcPkg.VerifyTransactionRequest) (*grpcPkg.VerifyTransactionResponse, error) {
	return &grpcPkg.VerifyTransactionResponse{}, nil
}
