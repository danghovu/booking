package postgresql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type ExecerContext interface {
	sqlx.ExecerContext
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
