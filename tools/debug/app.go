package debug

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
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

	return &Debug{
		CMD: rootCmd,
	}
}

// send creates the Send command
func send() *cobra.Command {
	var (
		sender   string
		receiver string
		amount   uint64
		unit     uint32
	)

	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send money from sender to receiver",
		Long:  "Send a specified amount of money from the sender to the receiver",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// Check if sender already exists
			userSender, err := client.GetUser(ctx, &transport.GetUserRequest{Username: sender})
			if err != nil {
				return fmt.Errorf("failed to get user: %v", err)
			}
			// Check if receiver already exists
			userReceiver, err := client.GetUser(ctx, &transport.GetUserRequest{Username: receiver})
			if err != nil {
				return fmt.Errorf("failed to get user: %v", err)
			}

			resp, err := client.AddTransaction(ctx, &transport.AddTransactionRequest{
				Sender:   userSender.GetUser().GetPrivateKey(),
				Receiver: userReceiver.GetUser().GetPublicKey(),
				Amount:   &transport.Amount{Value: amount, Unit: unit},
			})
			if err != nil {
				return fmt.Errorf("failed to add transaction: %w", err)
			}

			fmt.Printf("✅ Transaction added successfully!\n")
			fmt.Printf("Response: %v\n", resp)
			return nil
		},
	}

	cmd.Flags().StringVarP(&sender, "sender", "s", "", "Sender username (required)")
	cmd.Flags().StringVarP(&receiver, "receiver", "r", "", "Receiver username (required)")
	cmd.Flags().Uint64VarP(&amount, "amount", "a", 0, "Amount to transfer (required)")
	cmd.Flags().Uint32VarP(&unit, "unit", "u", 100, "Unit/precision for the amount")

	if err := cmd.MarkFlagRequired("sender"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("receiver"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("amount"); err != nil {
		panic(err)
	}

	return cmd
}

// balance creates the balance command
func balance() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Get balance by user name",
		Long:  "Get the balance of a user by their username",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			user, err := client.GetUser(ctx, &transport.GetUserRequest{Username: name})
			if err != nil {
				return fmt.Errorf("failed to get user: %v", user)
			}

			resp, err := client.GetBalance(ctx, &transport.GetBalanceRequest{
				Sender: user.GetUser().PrivateKey,
			})
			if err != nil {
				return fmt.Errorf("failed to get balance: %w", err)
			}

			fmt.Printf("💰 Balance for user '%s': %d (unit: %d)\n", name, resp.Amount.Value, resp.Amount.Unit)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Username (required)")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

// addUser creates the add user command
func addUser() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add-user",
		Short: "Add a new user",
		Long:  "Add a new user with generated ECDSA key pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// Check if user already exists
			user, err := client.GetUser(ctx, &transport.GetUserRequest{Username: name})
			if err == nil {
				return fmt.Errorf("failed to get user: %v", err)
			}
			if user.GetUser() != nil {
				return fmt.Errorf("user already exists: %s", name)
			}

			privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				return fmt.Errorf("failed to generate private key: %w", err)
			}

			privBytes, err := x509.MarshalECPrivateKey(privateKey)
			if err != nil {
				return fmt.Errorf("failed to marshal private key: %w", err)
			}

			privPEM := pem.EncodeToMemory(&pem.Block{
				Type:  "EC PRIVATE KEY",
				Bytes: privBytes,
			})

			// Encode public key to PEM
			pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
			if err != nil {
				return fmt.Errorf("failed to marshal public key: %w", err)
			}

			pubPEM := pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: pubBytes,
			})

			if _, err = client.AddUser(ctx, &transport.AddUserRequest{
				User: &transport.User{Username: name, PrivateKey: privPEM, PublicKey: pubPEM},
			}); err != nil {
				return fmt.Errorf("failed to add user: %w", err)
			}

			fmt.Printf("✅ User '%s' added successfully!\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Username (required)")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		panic(err)
	}

	return cmd
}

// fullEmission creates the full emission command
func fullEmission() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "full-emission",
		Short: "Calculate total emission",
		Long:  "Calculate the sum of all users' balances to check for double spending",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			// Get all users
			usersResp, err := client.ListUsers(ctx, &transport.ListUsersRequest{})
			if err != nil {
				return fmt.Errorf("failed to list users: %w", err)
			}

			totalBalance := uint64(0)
			userBalances := make(map[string]uint64)

			fmt.Printf("📊 Calculating full emission for %d users...\n\n", len(usersResp.Users))

			// Get balance for each user
			for _, user := range usersResp.Users {
				ctx, cancel = context.WithTimeout(context.Background(), timeout)
				resp, err := client.GetBalance(ctx, &transport.GetBalanceRequest{
					Sender: user.PrivateKey,
				})
				cancel()
				if err != nil {
					log.Printf("Warning: failed to get balance for user %s: %v", user.Username, err)
					continue
				}

				balance := resp.Amount.Value
				userBalances[user.Username] = balance
				totalBalance += balance

				if balance > 0 {
					fmt.Printf("  💰 %-15s: %d\n", user.Username, balance)
				}
			}

			fmt.Printf("\n" + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("📈 Total Emission: %d\n", totalBalance)
			fmt.Printf("👥 Total Users: %d\n", len(usersResp.Users))
			fmt.Printf("💵 Users with Balance: %d\n", countNonZeroBalances(userBalances))
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

			return nil
		},
	}

	return cmd
}

// addPeer creates the add peer command
func addPeer() *cobra.Command {
	var (
		peerAddr string
		peerID   string
	)

	cmd := &cobra.Command{
		Use:   "add-peer",
		Short: "Add a new peer to the network",
		Long:  "Add a new peer node to the blockchain network",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			resp, err := client.AddPeer(ctx, &transport.AddPeerRequest{
				Id:      peerID,
				Address: peerAddr,
			})
			if err != nil {
				return fmt.Errorf("failed to add peer: %w", err)
			}

			fmt.Printf("✅ Peer added successfully!\n")
			fmt.Printf("  Address: %s\n", peerAddr)
			fmt.Printf("  ID: %s\n", peerID)
			fmt.Printf("Response: %v\n", resp)
			return nil
		},
	}

	cmd.Flags().StringVarP(&peerAddr, "address", "a", "", "Peer address (e.g., 127.0.0.1:9002) (required)")
	cmd.Flags().StringVarP(&peerID, "id", "i", "", "Peer ID (required)")

	if err := cmd.MarkFlagRequired("address"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}

	return cmd
}

// addVoter creates the add voter command
func addVoter() *cobra.Command {
	var (
		voterID   string
		voterAddr string
	)

	cmd := &cobra.Command{
		Use:   "add-voter",
		Short: "Add a new voter",
		Long:  "Add a new voter to the blockchain network for consensus participation",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			resp, err := client.AddVoter(ctx, &transport.AddVoterRequest{
				Id:      voterID,
				Address: voterAddr,
			})
			if err != nil {
				return fmt.Errorf("failed to add voter: %w", err)
			}

			fmt.Printf("✅ Voter '%s' added successfully!\n", voterID)
			fmt.Printf("Response: %v\n", resp)
			return nil
		},
	}

	cmd.Flags().StringVarP(&voterID, "id", "", "", "Voter id (required)")
	cmd.Flags().StringVarP(&voterAddr, "address", "", "", "Voter address (required)")

	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}
	if err := cmd.MarkFlagRequired("address"); err != nil {
		panic(err)
	}

	return cmd
}

// removePeer creates the remove peer command
func removePeer() *cobra.Command {
	var peerID string

	cmd := &cobra.Command{
		Use:   "remove-peer",
		Short: "Remove a peer from the network",
		Long:  "Remove an existing peer node from the blockchain network",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			resp, err := client.RemovePeer(ctx, &transport.RemovePeerRequest{
				Id: peerID,
			})
			if err != nil {
				return fmt.Errorf("failed to remove peer: %w", err)
			}

			fmt.Printf("✅ Peer removed successfully!\n")
			fmt.Printf("  ID: %s\n", peerID)
			fmt.Printf("Response: %v\n", resp)
			return nil
		},
	}

	cmd.Flags().StringVarP(&peerID, "id", "i", "", "Peer ID to remove (required)")

	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}

	return cmd
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
