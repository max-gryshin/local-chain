package debug

import (
	"context"
	"fmt"
	"local-chain/transport/gen/transport"
	"log"

	"github.com/spf13/cobra"
)

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
