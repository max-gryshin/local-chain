package contractions

import "github.com/ZmaximillianZ/local-chain/internal/models"

// AccountRepository is interface to communicate with account storage
type AccountRepository interface {
	GetByID(id int) (models.Account, error)
	GetAll() (models.Accounts, error)
	Create(account *models.Account) error
	Update(account *models.Account) error
	Delete(account *models.Account) error
}
