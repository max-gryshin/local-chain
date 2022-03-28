package contractions

import "github.com/ZmaximillianZ/local-chain/internal/models"

// OrderRepository is interface to communicate with order storage
type OrderRepository interface {
	GetByID(id int) (models.Order, error)
	GetAll() (models.Orders, error)
	Create(order *models.Order) error
	Update(order *models.Order) error
	Delete(order *models.Order) error
}
