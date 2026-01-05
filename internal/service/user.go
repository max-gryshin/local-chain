package service

import "local-chain/internal/types"

type UserStore interface {
	GetAll() ([]*types.User, error)
	Get(username string) (*types.User, error)
	Put(user *types.User) error
}

type User struct {
	userStore UserStore
}

func NewUserService(userStore UserStore) *User {
	return &User{
		userStore: userStore,
	}
}

func (s *User) GetAllUsers() ([]*types.User, error) {
	return s.userStore.GetAll()
}

func (s *User) GetUser(username string) (*types.User, error) {
	return s.userStore.Get(username)
}

func (s *User) AddUser(user *types.User) error {
	return s.userStore.Put(user)
}
