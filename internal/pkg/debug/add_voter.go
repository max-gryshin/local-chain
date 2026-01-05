package debug

import (
	"context"
	"fmt"
	"local-chain/transport/gen/transport"

	"github.com/spf13/cobra"
)

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
