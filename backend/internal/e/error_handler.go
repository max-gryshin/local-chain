package e

import (
	"net/http"

	"github.com/max-gryshin/local-chain/internal/logging"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
)

type ErrorHandler struct {
}

func (eh *ErrorHandler) Handle(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		e := next(c)
		if e == nil {
			return nil
		}
		if pg := handlePgError(e); pg != nil {
			return echo.NewHTTPError(
				http.StatusUnprocessableEntity,
				map[string]string{
					"detail":  pg.Detail,
					"message": pg.Message,
					"code":    pg.Code,
				},
			)
		}
		logging.Error(e)
		return echo.NewHTTPError(http.StatusBadRequest, e.Error())
	}
}

func handlePgError(e error) *pgconn.PgError {
	if err, ok := e.(*pgconn.PgError); ok {
		return err
	}
	return nil
}
