package contractions

import "github.com/ZmaximillianZ/local-chain/internal/models"

// WalletRepository is interface to communicate with wallet storage
type WalletRepository interface {
	GetByID(id int) (models.Wallet, error)
	GetAll() (models.Wallets, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	Delete(wallet *models.Wallet) error
}
