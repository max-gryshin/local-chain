package runners

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"local-chain/internal/pkg"
)

type GrpcRunner struct {
	server *grpc.Server
	logger *slog.Logger

	port int
}

func New(port int, reg func(s *grpc.Server), logger slog.Logger) *GrpcRunner {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(panicRecoveryHandler(logger))),
		),
	)

	reflection.Register(server)

	reg(server)

	return &GrpcRunner{
		port:   port,
		server: server,

		logger: logger.With("source", "internal/runners/grpc"),
	}
}

func (r *GrpcRunner) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", r.port))
	if err != nil {
		return fmt.Errorf("err with listen grpc address: %v", err)
	}

	return pkg.Run(
		ctx,
		r.logger,
		pkg.Func(func(_ context.Context) error {
			if err := r.server.Serve(listener); err != nil {
				return fmt.Errorf("err with serve grpc: %v", err)
			}

			return nil
		}),
		pkg.Func(func(ctx context.Context) error {
			<-ctx.Done()
			r.server.GracefulStop()

			return nil
		}),
	)
}

func panicRecoveryHandler(l slog.Logger) func(any) error {
	return func(p any) error {
		l.Error("panic", p, string(debug.Stack()))
		return status.Errorf(codes.Internal, "%s", p)
	}
}
