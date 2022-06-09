package repository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/max-gryshin/local-chain/internal/logging"
	"github.com/max-gryshin/local-chain/internal/models"
	"github.com/max-gryshin/local-chain/internal/utils"
)

// WalletRepository is repository implementation for models.Wallets
type WalletRepository struct {
	BaseRepository
}

// NewWalletRepository creates new instance of WalletRepository
func NewWalletRepository(db *sqlx.DB, queryBuilder goqu.DialectWrapper) *WalletRepository {
	table := `wallet`
	fields := utils.GetTagValue(models.Wallet{}, tagName)
	baseQuery := queryBuilder.From(table).Select(fields...).Prepared(true)

	return &WalletRepository{BaseRepository{
		db:           db,
		table:        table,
		baseQuery:    baseQuery,
		queryBuilder: queryBuilder,
	}}
}

func (repo *WalletRepository) GetByID(id int) (models.Wallet, error) {
	wallet := models.Wallet{}
	sql, params, err := repo.baseQuery.Where(exp.Ex{"id": id}).ToSQL()
	if err != nil {
		logging.Error(err)

		return wallet, err
	}
	err = repo.db.Get(&wallet, sql, params...)
	if err != nil {
		logging.Error(err)

		return wallet, err
	}

	return wallet, nil
}

func (repo *WalletRepository) GetAll() (models.Wallets, error) {
	var wallets = models.Wallets{}
	query := repo.baseQuery.Limit(maxItemsPerPage)
	sql, p, err := query.ToSQL()
	if err != nil {
		return wallets, err
	}

	err = repo.db.Select(&wallets, sql, p...)

	return wallets, err
}

func (repo *WalletRepository) Create(wallet *models.Wallet) error {
	query := repo.
		baseQuery.
		Insert().
		Into(`wallet`).
		Cols(
			"status",
			"wallet_id",
			"private_key",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
			"account_id",
		).
		Vals(goqu.Vals{
			wallet.Status,
			wallet.WalletID,
			wallet.PrivateKey,
			wallet.CreatedAt,
			wallet.UpdatedAt,
			wallet.CreatedBy,
			wallet.UpdatedBy,
			wallet.AccountID,
		})

	return repo.execInsert(query)
}

func (repo *WalletRepository) Update(wallet *models.Wallet) error {
	expr := repo.baseQuery.Update().Set(wallet).Where(exp.Ex{"id": wallet.ID})
	return repo.execUpdate(expr)
}

func (repo *WalletRepository) Delete(wallet *models.Wallet) error {
	expr := repo.baseQuery.Delete().Where(exp.Ex{"id": wallet.ID})
	return repo.execDelete(expr)
}
