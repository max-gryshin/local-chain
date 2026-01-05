package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"local-chain/internal/adapters/outbound/leveldb"

	"local-chain/internal/types"
)

const (
	keysDir       = "/app/keys"
	superUserName = "admin"
)

func initSuperUser(us *leveldb.UserStore) *types.User {
	privPath := filepath.Join(keysDir, fmt.Sprintf("%s-priv.pem", superUserName))
	pubPath := filepath.Join(keysDir, fmt.Sprintf("%s-pub.pem", superUserName))

	privKey, err := os.ReadFile(privPath)
	if err != nil {
		log.Printf("user_data_set: skip %s, can't read private key %q: %v", superUserName, privPath, err)
	}
	pubKey, err := os.ReadFile(pubPath)
	if err != nil {
		log.Printf("user_data_set: skip %s, can't read public key %q: %v", superUserName, pubPath, err)
	}
	user := &types.User{
		Username:   superUserName,
		PublicKey:  pubKey,
		PrivateKey: privKey,
	}
	if err = us.Put(user); err != nil {
		log.Printf("user_data_set: failed to add super user %s: %v", superUserName, err)
	}
	return user
}
