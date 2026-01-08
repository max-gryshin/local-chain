package debug

import (
	"fmt"
	"log"
	"time"

	"local-chain/transport/gen/transport"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr string
	timeout    time.Duration
)

type Debug struct {
	CMD *cobra.Command
}

func NewDebug() *Debug {
	rootCmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug tool for local-chain blockchain",
		Long:  "A CLI tool for debugging and interacting with the local-chain blockchain",
	}
	// global flags
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", "127.0.0.1:9001", "gRPC server address")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "Request timeout")

	rootCmd.AddCommand(send())
	rootCmd.AddCommand(balance())
	rootCmd.AddCommand(addUser())
	rootCmd.AddCommand(fullEmission())
	rootCmd.AddCommand(addPeer())
	rootCmd.AddCommand(addVoter())
	rootCmd.AddCommand(removePeer())
	rootCmd.AddCommand(verifyTransaction())

	return &Debug{
		CMD: rootCmd,
	}
}

// createClient creates a gRPC client connection
func createClient() (transport.LocalChainClient, func(), error) {
	conn, err := grpc.NewClient(serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := transport.NewLocalChainClient(conn)
	closeFunc := func() {
		if err := conn.Close(); err != nil {
			log.Printf("Warning: failed to close connection: %v", err)
		}
	}

	return client, closeFunc, nil
}

// countNonZeroBalances counts users with non-zero balances
func countNonZeroBalances(balances map[string]uint64) int {
	count := 0
	for _, balance := range balances {
		if balance > 0 {
			count++
		}
	}
	return count
}
