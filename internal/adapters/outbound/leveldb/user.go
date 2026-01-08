package leveldb

import (
	"errors"
	"fmt"

	"local-chain/internal/types"

	"github.com/ethereum/go-ethereum/rlp"
	leveldbErrors "github.com/syndtr/goleveldb/leveldb/errors"
)

type userS struct {
	db Database
}

func newUserStore(conn Database) *userS {
	return &userS{
		db: conn,
	}
}

func (s *userS) GetAll() ([]*types.User, error) {
	keys, err := s.getKeys()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}
	users := make([]*types.User, 0, len(keys))
	for _, key := range keys {
		raw, err := s.db.Get(key, nil)
		if err != nil && !errors.Is(err, leveldbErrors.ErrNotFound) {
			return nil, fmt.Errorf("UserStore.Get get user error: %w", err)
		}
		if raw == nil {
			return nil, ErrNotFound
		}
		var user *types.User
		if err = rlp.DecodeBytes(raw, &user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *userS) Get(username string) (*types.User, error) {
	raw, err := s.db.Get([]byte(username), nil)
	if err != nil && !errors.Is(err, leveldbErrors.ErrNotFound) {
		return nil, fmt.Errorf("UserStore.Get get user error: %w", err)
	}
	if raw == nil {
		return nil, ErrNotFound
	}
	var user *types.User
	if err = rlp.DecodeBytes(raw, &user); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return user, nil
}

func (s *userS) Put(user *types.User) error {
	existingUser, err := s.Get(user.Username)
	if err != nil &&
		!errors.Is(err, leveldbErrors.ErrNotFound) &&
		!errors.Is(err, ErrNotFound) {
		return fmt.Errorf("UserStore.Put get user error: %w", err)
	}
	if existingUser != nil {
		return errors.New("UserStore.Put already exists")
	}

	encoded, err := rlp.EncodeToBytes(user)
	if err != nil {
		return fmt.Errorf("failed to encode user: %w", err)
	}
	if err = s.db.Put([]byte(user.Username), encoded, nil); err != nil {
		return fmt.Errorf("failed to put new user: %w", err)
	}
	return nil
}

func (s *userS) getKeys() ([][]byte, error) {
	iterator := s.db.NewIterator(nil, nil)
	defer iterator.Release()

	var keys [][]byte
	for iterator.Next() {
		// Make a copy of the key since the iterator reuses the buffer
		key := make([]byte, len(iterator.Key()))
		copy(key, iterator.Key())
		keys = append(keys, key)
	}

	if err := iterator.Error(); err != nil {
		return nil, fmt.Errorf("failed to iterate over keys: %w", err)
	}
	return keys, nil
}
