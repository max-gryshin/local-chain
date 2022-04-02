package contractions

import "github.com/ZmaximillianZ/local-chain/internal/models"

// UserRepository is interface to communicate with user storage
type UserRepository interface {
	GetByID(id int) (models.User, error)
	GetByEmail(email string) (models.User, error)
	GetAll() (models.Users, error)
	Create(user *models.User) error
	Update(user *models.User) error
}
