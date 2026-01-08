package debug

import (
	"context"
	"fmt"

	"local-chain/transport/gen/transport"

	"github.com/spf13/cobra"
)

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
