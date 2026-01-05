package debug

import (
	"context"
	"fmt"
	"local-chain/transport/gen/transport"

	"github.com/spf13/cobra"
)

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
