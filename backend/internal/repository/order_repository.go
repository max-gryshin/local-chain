package repository

import (
	"github.com/ZmaximillianZ/local-chain/internal/logging"
	"github.com/ZmaximillianZ/local-chain/internal/models"
	"github.com/ZmaximillianZ/local-chain/internal/utils"
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

// OrderRepository is repository implementation for models.Orders
type OrderRepository struct {
	BaseRepository
}

// NewOrderRepository creates new instance of OrderRepository
func NewOrderRepository(db *sqlx.DB, queryBuilder goqu.DialectWrapper) *OrderRepository {
	table := `order`
	fields := utils.GetTagValue(models.Order{}, tagName)
	baseQuery := queryBuilder.From(table).Select(fields...).Prepared(true)

	return &OrderRepository{BaseRepository{
		db:           db,
		table:        table,
		baseQuery:    baseQuery,
		queryBuilder: queryBuilder,
	}}
}

func (repo *OrderRepository) GetByID(id int) (models.Order, error) {
	order := models.Order{}
	sql, params, err := repo.baseQuery.Where(exp.Ex{"id": id}).ToSQL()
	if err != nil {
		logging.Error(err)

		return order, err
	}
	err = repo.db.Get(&order, sql, params...)
	if err != nil {
		logging.Error(err)

		return order, err
	}

	return order, nil
}

func (repo *OrderRepository) GetAll() (models.Orders, error) {
	var orders = models.Orders{}
	query := repo.baseQuery.Limit(maxItemsPerPage)
	sql, p, err := query.ToSQL()
	if err != nil {
		return orders, err
	}

	err = repo.db.Select(&orders, sql, p...)

	return orders, err
}

func (repo *OrderRepository) Create(order *models.Order) error {
	query := repo.
		baseQuery.
		Insert().
		Into(`order`).
		Cols(
			"status",
			"amount",
			"wallet_id",
			"description",
			"request_reasons",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
		).
		Vals(goqu.Vals{
			order.Status,
			order.Amount,
			order.WalletID,
			order.Description,
			order.RequestReason,
			order.CreatedAt,
			order.UpdatedAt,
			order.CreatedBy,
			order.UpdatedBy,
		})

	return repo.execInsert(query)
}

func (repo *OrderRepository) Update(order *models.Order) error {
	expr := repo.baseQuery.Update().Set(order).Where(exp.Ex{"id": order.ID})
	return repo.execUpdate(expr)
}

func (repo *OrderRepository) Delete(order *models.Order) error {
	expr := repo.baseQuery.Delete().Where(exp.Ex{"id": order.ID})
	return repo.execDelete(expr)
}
