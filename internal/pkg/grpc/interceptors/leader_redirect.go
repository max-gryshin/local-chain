package interceptors

import (
	"context"
	"net"

	"local-chain/internal/service"

	"github.com/hashicorp/raft"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	grpcPkg "local-chain/transport/gen/transport"
)

type contextKey string

const (
	// skipLeaderRedirectKey is used to mark requests that should skip leader redirection
	skipLeaderRedirectKey contextKey = "skipLeaderRedirect"
	leaderPort            string     = "9001"

	grpcSrvPrefix            string = "/LocalChain/"
	grpcMethodAddPeer               = grpcSrvPrefix + "AddPeer"
	grpcMethodRemovePeer            = grpcSrvPrefix + "RemovePeer"
	grpcMethodAddVoter              = grpcSrvPrefix + "AddVoter"
	grpcMethodAddTransaction        = grpcSrvPrefix + "AddTransaction"
	grpcMethodGetBalance            = grpcSrvPrefix + "GetBalance"
	grpcMethodAddUser               = grpcSrvPrefix + "AddUser"
	grpcMethodGetUser               = grpcSrvPrefix + "GetUser"
	grpcMethodListUsers             = grpcSrvPrefix + "ListUsers"
)

// LeaderRedirectInterceptor redirects requests to the leader node if the current node is not the leader.
type LeaderRedirectInterceptor struct {
	serverID raft.ServerID
	raftAPI  service.RaftAPI
	grpcPort string
}

// NewLeaderRedirectInterceptor creates a new leader redirect interceptor.
func NewLeaderRedirectInterceptor(serverID raft.ServerID, raftAPI service.RaftAPI) *LeaderRedirectInterceptor {
	return &LeaderRedirectInterceptor{
		serverID: serverID,
		raftAPI:  raftAPI,
		grpcPort: leaderPort,
	}
}

// UnaryInterceptor returns a gRPC unary interceptor that redirects requests to the leader.
func (i *LeaderRedirectInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Check is leader
		leaderServer, leaderID := i.raftAPI.LeaderWithID()
		if leaderID == i.serverID {
			return handler(ctx, req)
		}

		client, err := i.createLeaderClient(string(leaderServer))
		if err != nil {
			return nil, err
		}

		// redirect to the leader
		return i.forwardToLeader(ctx, client, info.FullMethod, req)
	}
}

// createLeaderClient creates a gRPC client connected to the leader.
func (i *LeaderRedirectInterceptor) createLeaderClient(leaderAddr string) (grpcPkg.LocalChainClient, error) {
	host, _, err := net.SplitHostPort(leaderAddr)
	if err != nil {
		return nil, err
	}
	conn, err := grpc.NewClient(
		net.JoinHostPort(host, i.grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return grpcPkg.NewLocalChainClient(conn), nil
}

// forwardToLeader forwards the request to the leader based on the method name.
func (i *LeaderRedirectInterceptor) forwardToLeader(
	ctx context.Context,
	client grpcPkg.LocalChainClient,
	method string,
	req interface{},
) (interface{}, error) {
	switch method {
	case grpcMethodAddPeer:
		return client.AddPeer(ctx, req.(*grpcPkg.AddPeerRequest))
	case grpcMethodRemovePeer:
		return client.RemovePeer(ctx, req.(*grpcPkg.RemovePeerRequest))
	case grpcMethodAddVoter:
		return client.AddVoter(ctx, req.(*grpcPkg.AddVoterRequest))
	case grpcMethodAddTransaction:
		return client.AddTransaction(ctx, req.(*grpcPkg.AddTransactionRequest))
	case grpcMethodGetBalance:
		return client.GetBalance(ctx, req.(*grpcPkg.GetBalanceRequest))
	case grpcMethodAddUser:
		return client.AddUser(ctx, req.(*grpcPkg.AddUserRequest))
	case grpcMethodGetUser:
		return client.GetUser(ctx, req.(*grpcPkg.GetUserRequest))
	case grpcMethodListUsers:
		return client.ListUsers(ctx, req.(*emptypb.Empty))
	default:
		// If method is not recognized, return an error (shouldn't happen in practice)
		return nil, grpc.ErrServerStopped
	}
}
