package base

import (
	"context"

	myErrors "github.com/dokidokikoi/go-common/errors"
	meta "github.com/dokidokikoi/go-common/meta/option"

	"github.com/jackc/pgx/v5/pgconn"
)

func (p *PgModel[T]) Update(ctx context.Context, t *T, option *meta.UpdateOption) error {
	db := p.DB
	if option != nil {
		if len(option.Select) > 0 {
			var params []any
			for _, s := range option.Select {
				params = append(params, s)
			}
			if len(params) > 0 {
				first := params[0]
				params = params[1:]
				db = p.DB.Select(first, params...)
			}
		}
		if len(option.Omit) > 0 {
			db = db.Omit(option.Omit...)
		}
	}
	result := db.Updates(t)
	err := result.Error
	pgErr, ok := err.(*pgconn.PgError)
	if ok && pgErr.Code == "23505" {
		err = myErrors.ErrNameDuplicate
	}
	row := result.RowsAffected
	if row == 0 {
		err = myErrors.ErrNoUpdateRows
	}
	return err
}

func (p *PgModel[T]) UpdateByWhere(ctx context.Context, node *meta.WhereNode, example *T, option *meta.UpdateOption) error {
	db := p.DB
	if option != nil {
		if len(option.Select) > 0 {
			var params []any
			for _, s := range option.Select {
				params = append(params, s)
			}
			if len(params) > 0 {
				first := params[0]
				params = params[1:]
				db = p.DB.Select(first, params...)
			}
		}
		if len(option.Omit) > 0 {
			db = db.Omit(option.Omit...)
		}
	}
	result := CompositeQuery(db, node).Updates(*example)
	err := result.Error
	pgErr, ok := err.(*pgconn.PgError)
	if ok && pgErr.Code == "23505" {
		err = myErrors.ErrNameDuplicate
	}
	row := result.RowsAffected
	if row == 0 {
		err = myErrors.ErrNoUpdateRows
	}
	return err
}

func (p *PgModel[T]) UpdateCollection(ctx context.Context, t []*T, option *meta.UpdateCollectionOption) []error {
	db := p.DB
	if option != nil {
		if len(option.Select) > 0 {
			var params []any
			for _, s := range option.Select {
				params = append(params, s)
			}
			if len(params) > 0 {
				first := params[0]
				params = params[1:]
				db = p.DB.Select(first, params...)
			}
		}
		if len(option.Omit) > 0 {
			db = db.Omit(option.Omit...)
		}
	}
	var errors []error
	for _, up := range t {
		if e := db.Updates(up).Error; e != nil {
			errors = append(errors, e)
		}
	}
	return errors
}

func (p *PgModel[T]) Save(ctx context.Context, t *T, option *meta.UpdateOption) error {
	db := p.DB
	if option != nil {
		if len(option.Select) > 0 {
			var params []any
			for _, s := range option.Select {
				params = append(params, s)
			}
			if len(params) > 0 {
				first := params[0]
				params = params[1:]
				db = p.DB.Select(first, params...)
			}
		}
		if len(option.Omit) > 0 {
			db = db.Omit(option.Omit...)
		}
	}
	return db.Save(t).Error
}
