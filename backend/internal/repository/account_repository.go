package repository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/max-gryshin/local-chain/internal/logging"
	"github.com/max-gryshin/local-chain/internal/models"
	"github.com/max-gryshin/local-chain/internal/utils"
)

// AccountRepository is repository implementation for models.Accounts
type AccountRepository struct {
	BaseRepository
}

// NewAccountRepository creates new instance of AccountRepository
func NewAccountRepository(db *sqlx.DB, queryBuilder goqu.DialectWrapper) *AccountRepository {
	table := `account`
	fields := utils.GetTagValue(models.Account{}, tagName)
	baseQuery := queryBuilder.From(table).Select(fields...).Prepared(true)

	return &AccountRepository{BaseRepository{
		db:           db,
		table:        table,
		baseQuery:    baseQuery,
		queryBuilder: queryBuilder,
	}}
}

func (repo *AccountRepository) GetByID(id int) (models.Account, error) {
	account := models.Account{}
	sql, params, err := repo.baseQuery.Where(exp.Ex{"id": id}).ToSQL()
	if err != nil {
		logging.Error(err)

		return account, err
	}
	err = repo.db.Get(&account, sql, params...)
	if err != nil {
		logging.Error(err)

		return account, err
	}

	return account, nil
}

func (repo *AccountRepository) GetAll() (models.Accounts, error) {
	var accounts = models.Accounts{}
	query := repo.baseQuery.Limit(maxItemsPerPage)
	sql, p, err := query.ToSQL()
	if err != nil {
		return accounts, err
	}

	err = repo.db.Select(&accounts, sql, p...)

	return accounts, err
}

func (repo *AccountRepository) Create(account *models.Account) error {
	query := repo.
		baseQuery.
		Insert().
		Into(`account`).
		Cols("status", "phone", "created_at", "updated_at", "created_by", "updated_by", "user_id").
		Vals(goqu.Vals{
			account.Status,
			account.Phone,
			account.CreatedAt,
			account.UpdatedAt,
			account.CreatedBy,
			account.UpdatedBy,
			account.UserID,
		})

	return repo.execInsert(query)
}

func (repo *AccountRepository) Update(account *models.Account) error {
	expr := repo.baseQuery.Update().Set(account).Where(exp.Ex{"id": account.ID})
	return repo.execUpdate(expr)
}

func (repo *AccountRepository) Delete(account *models.Account) error {
	expr := repo.baseQuery.Delete().Where(exp.Ex{"id": account.ID})
	return repo.execDelete(expr)
}
