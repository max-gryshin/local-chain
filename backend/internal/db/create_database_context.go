package db

import (
	"fmt"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/v4/stdlib" // need to connect with db
	"github.com/jmoiron/sqlx"
)

type ConnectionSettions struct {
	Database    string
	URL         string
	MaxIdleCons int
	MaxOpenCons int
}

type DatabaseContext struct {
	Connection   *sqlx.DB
	QueryBuilder goqu.DialectWrapper
}

var dialects = map[string]string{
	"mysql":     "mysql",
	"postgres":  "postgres",
	"sqlite3":   "sqlite3",
	"sqlserver": "sqlserver",
}

var drivers = map[string]string{
	"mysql":     "mysql",
	"postgres":  "pgx",
	"sqlite3":   "sqlite3",
	"sqlserver": "sqlserver",
}

func CreateDatabaseContext(setting ConnectionSettions) (DatabaseContext, error) {
	db, err := sqlx.Connect(getValue(drivers, setting.Database), setting.URL)
	if err != nil {
		return DatabaseContext{}, err
	}

	db.SetMaxIdleConns(setting.MaxIdleCons)
	db.SetMaxOpenConns(setting.MaxIdleCons)

	return DatabaseContext{
		Connection:   db,
		QueryBuilder: goqu.Dialect(getValue(dialects, setting.Database)),
	}, nil
}

func getValue(data map[string]string, database string) string {
	if value, ok := data[database]; ok {
		return value
	}

	panic(fmt.Sprintf("database %s is not supported", database))
}
