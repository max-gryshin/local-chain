package leveldb

import (
	"errors"
	"fmt"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
)

const UsersKey = "users"

type UserStore struct {
	db Database
}

func NewUserStore(conn Database) *UserStore {
	return &UserStore{
		db: conn,
	}
}

func (s *UserStore) GetAll() ([]*types.User, error) {
	raw, err := s.db.Get([]byte(UsersKey), nil)
	if err != nil && !errors.As(err, &leveldbErrors.ErrNotFound) { // nolint:govet
		return nil, fmt.Errorf("UserStore.Get get user error: %w", err)
	}
	if raw == nil {
		return nil, nil
	}

	var users []*types.User
	if err = rlp.DecodeBytes(raw, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}
	return users, nil
}

func (s *UserStore) Get(username string) (*types.User, error) {
	users, err := s.GetAll()
	if err != nil { // nolint:govet
		return nil, fmt.Errorf("UserStore.Get get user error: %w", err)
	}
	for _, user := range users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, fmt.Errorf("UserStore.Get get user error: %w", err)
}

func (s *UserStore) Put(user *types.User) error {
	users, err := s.GetAll()
	if err != nil && !errors.Is(err, leveldbErrors.ErrNotFound) {
		return fmt.Errorf("UserStore.Put get user error: %w", err)
	}

	users = append(users, user) // nolint:ineffassign
	encoded, err := rlp.EncodeToBytes(users)
	if err != nil {
		return fmt.Errorf("failed to encode user: %w", err)
	}
	if err = s.db.Put([]byte(UsersKey), encoded, nil); err != nil {
		return fmt.Errorf("failed to put new user: %w", err)
	}
	return nil
}
