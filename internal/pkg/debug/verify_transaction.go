package debug

import (
	"context"
	"fmt"

	"local-chain/transport/gen/transport"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// verifyTransaction creates the verify-transaction command
func verifyTransaction() *cobra.Command {
	var txID string

	cmd := &cobra.Command{
		Use:   "verify-transaction",
		Short: "Verify a transaction by ID",
		Long:  "Call VerifyTransaction RPC and display whether the transaction is valid along with its details",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, closeConn, err := createClient()
			if err != nil {
				return err
			}
			defer closeConn()

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			txIDParsed := uuid.MustParse(txID)

			resp, err := client.VerifyTransaction(ctx, &transport.VerifyTransactionRequest{Id: []byte(txIDParsed.String())})
			if err != nil {
				return fmt.Errorf("failed to verify transaction: %w", err)
			}

			tx := resp.GetTransaction()
			if tx == nil {
				fmt.Printf("No transaction returned for id %s\n", txID)
				return nil
			}

			verdict := "❌ INVALID"
			if resp.GetIsValid() {
				verdict = "✅ VALID"
			}

			fmt.Printf("\n%s Transaction Verification\n\n", verdict)
			fmt.Printf("═══════════════════════════════════════════════════════\n")
			fmt.Printf("  ID:               %s\n", tx.GetId())
			fmt.Printf("  Timestamp:        %d\n", tx.GetTimestamp())
			fmt.Printf("  Hash:             %x\n", tx.GetHash())
			if tx.GetBlockTimestamp() > 0 {
				fmt.Printf("  Block Timestamp:  %x\n", tx.GetBlockTimestamp())
			}

			fmt.Printf("\n  INPUTS (%d):\n", len(tx.GetInputs()))
			for i, input := range tx.GetInputs() {
				fmt.Printf("    [%d] Public Key:  %x\n", i+1, input.GetPubKey())
				fmt.Printf("        SignatureS:   %x\n", input.GetSignatureS())
				fmt.Printf("        SignatureR:   %x\n", input.GetSignatureR())
			}

			fmt.Printf("\n  OUTPUTS (%d):\n", len(tx.GetOutputs()))
			for i, output := range tx.GetOutputs() {
				fmt.Printf("    [%d] Public Key:  %x\n", i+1, output.GetPubKey())
				fmt.Printf("        Amount:      %d (unit: %d)\n", output.GetAmount().GetValue(), output.GetAmount().GetUnit())
			}

			fmt.Printf("\n═══════════════════════════════════════════════════════\n\n")
			return nil
		},
	}

	cmd.Flags().StringVarP(&txID, "id", "i", "", "Transaction ID (hex) (required)")
	if err := cmd.MarkFlagRequired("id"); err != nil {
		panic(err)
	}

	return cmd
}
