package debug

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"local-chain/transport/gen/transport"

	"github.com/spf13/cobra"
)

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
