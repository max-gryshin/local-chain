package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	userDataSetKeys  = make(map[string]map[string][]byte)
	userDataSetNames = []string{
		"admin",
		"alice", "bob", "charlie", "dave", "eve", "frank", "grace", "hank", "irene", "jack",
		"karen", "leo", "mia", "nick", "olivia", "paul", "quinn", "rachel", "steve", "tina",
		"uma", "victor", "wendy", "xavier", "yvonne", "zack", "abby", "brad", "chris",
		"diana", "elena", "felix", "gina", "harry", "isabel", "jim", "kevin", "lara",
		"mike", "nora", "oscar", "peter", "queen", "ron", "sara", "tom", "ursula",
		"vince", "will", "xena", "yasmin", "zane", "aaron", "beth", "cody", "dana",
		"edgar", "fiona", "gary", "holly", "ivan", "jill", "kyle", "liam", "maggie",
		"neil", "opal", "priya", "quincy", "rose", "sean", "troy", "una", "vera",
		"walter", "xia", "yang", "zoe", "albert", "bella", "carl", "denise",
		"eric", "faith", "gabriel", "heidi", "ian", "jen", "ken", "luis", "mandy",
		"nate", "owen", "paula", "qadir", "rita", "sam", "tyler", "ugo", "val",
		"wayne", "xiao", "yara", "ziad",
	}
)

const keysDir = "cmd/tools/debug/keys"

func init() {
	for _, name := range userDataSetNames {
		userDataSetKeys[name] = map[string][]byte{}

		privPath := filepath.Join(keysDir, fmt.Sprintf("%s-priv.pem", name))
		pubPath := filepath.Join(keysDir, fmt.Sprintf("%s-pub.pem", name))

		priv, err := os.ReadFile(privPath)
		if err != nil {
			log.Printf("user_data_set: skip %s, can't read private key %q: %v", name, privPath, err)
			continue
		}
		pub, err := os.ReadFile(pubPath)
		if err != nil {
			log.Printf("user_data_set: skip %s, can't read public key %q: %v", name, pubPath, err)
			continue
		}
		userDataSetKeys[name][string(pub)] = priv
	}
}

// getUserKeys returns the private and public keys for a given username
func getUserKeys(name string) (privKey []byte, pubKey []byte, err error) {
	keys, exists := userDataSetKeys[name]
	if !exists {
		return nil, nil, fmt.Errorf("user not found: %s", name)
	}

	// Get the first (and only) key pair for this user
	for pub, priv := range keys {
		return priv, []byte(pub), nil
	}

	return nil, nil, fmt.Errorf("no keys found for user: %s", name)
}

// userExists checks if a user exists in the dataset
func userExists(name string) bool {
	_, exists := userDataSetKeys[name]
	return exists
}

// addUser adds a new user with generated keys
func addUser(name string, privKey, pubKey []byte) error {
	if userExists(name) {
		return fmt.Errorf("user already exists: %s", name)
	}

	userDataSetKeys[name] = map[string][]byte{
		string(pubKey): privKey,
	}

	// Save keys to files
	privPath := filepath.Join(keysDir, fmt.Sprintf("%s-priv.pem", name))
	pubPath := filepath.Join(keysDir, fmt.Sprintf("%s-pub.pem", name))

	if err := os.WriteFile(privPath, privKey, 0600); err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	if err := os.WriteFile(pubPath, pubKey, 0644); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}

// getAllUsers returns a list of all user names
func getAllUsers() []string {
	users := make([]string, 0, len(userDataSetKeys))
	for name := range userDataSetKeys {
		users = append(users, name)
	}
	return users
}
