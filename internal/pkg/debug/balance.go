package debug

import (
	"context"
	"fmt"
	"local-chain/transport/gen/transport"

	"github.com/spf13/cobra"
)

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
