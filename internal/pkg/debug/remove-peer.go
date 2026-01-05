package debug

import (
	"context"
	"fmt"
	"local-chain/transport/gen/transport"

	"github.com/spf13/cobra"
)

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
