package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"local-chain/transport/gen/transport"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr string
	timeout    time.Duration
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "debug",
		Short: "Debug tool for local-chain blockchain",
		Long:  "A CLI tool for debugging and interacting with the local-chain blockchain",
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&serverAddr, "server", "127.0.0.1:9001", "gRPC server address")
	rootCmd.PersistentFlags().DurationVar(&timeout, "timeout", 10*time.Second, "Request timeout")

	// Add commands
	rootCmd.AddCommand(addTransactionCmd())
	rootCmd.AddCommand(getBalanceCmd())
	rootCmd.AddCommand(addUserCmd())
	rootCmd.AddCommand(fullEmissionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// addTransactionCmd creates the add transaction command
func addTransactionCmd() *cobra.Command {
	var (
		sender   string
		receiver string
		amount   uint64
		unit     uint32
	)

	cmd := &cobra.Command{
		Use:   "add-transaction",
		Short: "Add a new transaction",
		Long:  "Add a new transaction from sender to receiver with specified amount",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate users exist
			if !userExists(sender) {
				return fmt.Errorf("user not found: %s", sender)
			}
			if !userExists(receiver) {
				return fmt.Errorf("user not found: %s", receiver)
			}

			// Get keys
			senderPriv, _, err := getUserKeys(sender)
			if err != nil {
				return fmt.Errorf("failed to get sender keys: %w", err)
			}

			_, receiverPub, err := getUserKeys(receiver)
			if err != nil {
				return fmt.Errorf("failed to get receiver keys: %w", err)
			}

			// Create gRPC client
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			// Add transaction
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			resp, err := client.AddTransaction(ctx, &transport.AddTransactionRequest{
				Sender:   senderPriv,
				Receiver: receiverPub,
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

	cmd.MarkFlagRequired("sender")
	cmd.MarkFlagRequired("receiver")
	cmd.MarkFlagRequired("amount")

	return cmd
}

// getBalanceCmd creates the get balance command
func getBalanceCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "balance",
		Short: "Get balance by user name",
		Long:  "Get the balance of a user by their username",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validate user exists
			if !userExists(name) {
				return fmt.Errorf("user not found: %s", name)
			}

			// Get keys
			userPriv, _, err := getUserKeys(name)
			if err != nil {
				return fmt.Errorf("failed to get user keys: %w", err)
			}

			// Create gRPC client
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			// Get balance
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			resp, err := client.GetBalance(ctx, &transport.GetBalanceRequest{
				Sender: userPriv,
			})
			if err != nil {
				return fmt.Errorf("failed to get balance: %w", err)
			}

			fmt.Printf("💰 Balance for user '%s': %d (unit: %d)\n", name, resp.Amount.Value, resp.Amount.Unit)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Username (required)")
	cmd.MarkFlagRequired("name")

	return cmd
}

// addUserCmd creates the add user command
func addUserCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "add-user",
		Short: "Add a new user",
		Long:  "Add a new user with generated ECDSA key pair",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if user already exists
			if userExists(name) {
				return fmt.Errorf("user already exists: %s", name)
			}

			// Generate new ECDSA key pair
			privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				return fmt.Errorf("failed to generate private key: %w", err)
			}

			// Encode private key to PEM
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

			// Add user to dataset
			if err := addUser(name, privPEM, pubPEM); err != nil {
				return fmt.Errorf("failed to add user: %w", err)
			}

			fmt.Printf("✅ User '%s' added successfully!\n", name)
			fmt.Printf("Keys saved to:\n")
			fmt.Printf("  - cmd/tools/debug/keys/%s-priv.pem\n", name)
			fmt.Printf("  - cmd/tools/debug/keys/%s-pub.pem\n", name)
			return nil
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Username (required)")
	cmd.MarkFlagRequired("name")

	return cmd
}

// fullEmissionCmd creates the full emission command
func fullEmissionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "full-emission",
		Short: "Calculate total emission",
		Long:  "Calculate the sum of all users' balances to check for double spending",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create gRPC client
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			// Get all users
			users := getAllUsers()
			totalBalance := uint64(0)
			userBalances := make(map[string]uint64)

			fmt.Printf("📊 Calculating full emission for %d users...\n\n", len(users))

			// Get balance for each user
			for _, username := range users {
				userPriv, _, err := getUserKeys(username)
				if err != nil {
					log.Printf("Warning: failed to get keys for user %s: %v", username, err)
					continue
				}

				ctx, cancel := context.WithTimeout(context.Background(), timeout)
				resp, err := client.GetBalance(ctx, &transport.GetBalanceRequest{
					Sender: userPriv,
				})
				cancel()

				if err != nil {
					log.Printf("Warning: failed to get balance for user %s: %v", username, err)
					continue
				}

				balance := resp.Amount.Value
				userBalances[username] = balance
				totalBalance += balance

				if balance > 0 {
					fmt.Printf("  💰 %-15s: %d\n", username, balance)
				}
			}

			fmt.Printf("\n" + "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("📈 Total Emission: %d\n", totalBalance)
			fmt.Printf("👥 Total Users: %d\n", len(users))
			fmt.Printf("💵 Users with Balance: %d\n", countNonZeroBalances(userBalances))
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

			return nil
		},
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
