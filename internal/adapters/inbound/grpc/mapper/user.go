package mapper

import (
	grpcPkg "local-chain/transport/gen/transport"

	"local-chain/internal/types"
)

type UserMapper struct{}

func NewUserMapper() *UserMapper {
	return &UserMapper{}
}

func (u *UserMapper) RpcToUser(req *grpcPkg.AddUserRequest) *types.User {
	return &types.User{
		Username:   req.GetUser().GetUsername(),
		PublicKey:  req.GetUser().GetPublicKey(),
		PrivateKey: req.GetUser().GetPrivateKey(),
	}
}

func (u *UserMapper) UserToRpc(user *types.User) *grpcPkg.User {
	return &grpcPkg.User{
		Username:   user.Username,
		PublicKey:  user.PublicKey,
		PrivateKey: user.PrivateKey,
	}
}
