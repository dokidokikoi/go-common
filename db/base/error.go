package base

import (
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
)

func errorHandle(err error) error {
	switch e := err.(type) {
	case *pgconn.PgError:
		if e.Code == "23505" {
			return gorm.ErrDuplicatedKey
		}
	case sqlite3.Error:
		if e.Code == sqlite3.ErrConstraint {
			return gorm.ErrDuplicatedKey
		}
	}
	return err
}
