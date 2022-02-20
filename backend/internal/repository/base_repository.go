package repository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
)

type BaseRepository struct {
	db           *sqlx.DB
	table        string
	baseQuery    *goqu.SelectDataset
	queryBuilder goqu.DialectWrapper
}

func (repo *BaseRepository) execInsert(data *goqu.InsertDataset) error {
	sql, params, err := data.ToSQL()
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(sql, params...)

	return err
}

func (repo *BaseRepository) execUpdate(data *goqu.UpdateDataset) error {
	sql, params, err := data.ToSQL()
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(sql, params...)

	return err
}

func (repo *BaseRepository) execDelete(data *goqu.DeleteDataset) error {
	sql, params, err := data.ToSQL()
	if err != nil {
		return err
	}

	_, err = repo.db.Exec(sql, params...)

	return err
}
