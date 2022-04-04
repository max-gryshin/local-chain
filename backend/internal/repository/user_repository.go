package repository

import (
	"strings"

	"github.com/ZmaximillianZ/local-chain/internal/logging"
	"github.com/ZmaximillianZ/local-chain/internal/models"
	"github.com/ZmaximillianZ/local-chain/internal/utils"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres" // need to import right dialect
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
)

const tagName = "db"
const maxItemsPerPage = 100

// UserRepository is repository implementation for models.Users
type UserRepository struct {
	BaseRepository
}

// NewUserRepository creates new instance of UserRepository
func NewUserRepository(db *sqlx.DB, queryBuilder goqu.DialectWrapper) *UserRepository {
	table := `user`
	fields := utils.GetTagValue(models.User{}, tagName)
	baseQuery := queryBuilder.From(table).Select(fields...).Prepared(true)

	return &UserRepository{BaseRepository{
		db:           db,
		table:        table,
		baseQuery:    baseQuery,
		queryBuilder: queryBuilder,
	}}
}

func (repo *UserRepository) GetByID(id int) (models.User, error) {
	user := models.User{}
	sql, params, err := repo.baseQuery.Where(exp.Ex{"id": id}).ToSQL()
	if err != nil {
		logging.Error(err)

		return user, err
	}
	err = repo.db.Get(&user, sql, params...)
	if err != nil {
		logging.Error(err)

		return user, err
	}

	return user, nil
}

func (repo *UserRepository) GetByEmail(email string) (models.User, error) {
	user := models.User{}
	sql, _, err := repo.
		baseQuery.
		Where(exp.Ex{"email": email}).
		ToSQL()
	if err != nil {
		logging.Error(err)
		return user, err
	}
	err = repo.db.Get(&user, sql, email)
	if err != nil {
		if strings.Contains(err.Error(), "no rows in result set") {
			return user, nil
		}
		logging.Error(err)
		return user, err
	}

	return user, nil
}

func (repo *UserRepository) GetAll() (models.Users, error) {
	var user = models.Users{}
	query := repo.baseQuery.Limit(maxItemsPerPage)
	sql, p, err := query.ToSQL()
	if err != nil {
		return user, err
	}

	err = repo.db.Select(&user, sql, p...)

	return user, err
}

func (repo *UserRepository) Create(user *models.User) error {
	query := repo.
		baseQuery.
		Insert().
		Into(`user`).
		Cols(
			"email",
			"password_hash",
			"status",
			"first_name",
			"last_name",
			"middle_name",
			"created_at",
			"updated_at",
			"created_by",
			"updated_by",
			"roles",
		).
		Vals(goqu.Vals{
			user.Email,
			user.Password,
			user.Status,
			user.FirstName,
			user.LastName,
			user.MiddleName,
			user.CreatedAt,
			user.UpdatedAt,
			user.CreatedBy,
			user.UpdatedBy,
			user.Roles,
		})

	return repo.execInsert(query)
}

func (repo *UserRepository) Update(user *models.User) error {
	userMap, err := utils.GetMapFromModel(user)
	if err != nil {
		return err
	}
	expr := repo.baseQuery.Update().Set(userMap).Where(exp.Ex{"id": user.ID})
	return repo.execUpdate(expr)
}
