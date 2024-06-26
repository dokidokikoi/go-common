package base

import (
	"context"

	myErrors "github.com/dokidokikoi/go-common/errors"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *PgModel[T]) Create(ctx context.Context, t *T, option *meta.CreateOption) error {
	db := p.DB
	if option != nil && len(option.Omit) > 0 {
		db = db.Omit(option.Omit...)
	}
	err := db.Create(t).Error
	pgErr, ok := err.(*pgconn.PgError)
	if ok && pgErr.Code == "23505" {
		err = myErrors.ErrNameDuplicate
	}
	return err
}

// func (p *PgModel[T]) CreateMany2Many(ctx context.Context, t *T, ids interface{}, option *meta.CreateOption) error {
// 	if len(option.Omit) > 0 {
// 		err := p.DB.Omit(option.Omit...).Create(t).Error
// 		if err != nil {
// 			pgErr, ok := err.(*pgconn.PgError)
// 			if ok && pgErr.Code == "23505" {
// 				err = myErrors.ErrNameDuplicate
// 			}
// 			return err
// 		}
// 		association := p.DB.Model(t).Association(option.Omit...)
// 		return association.Append(ids)
// 	}

// 	return errors.New("未指定关联字段名")
// }

func (p *PgModel[T]) Creates(ctx context.Context, t []*T, option *meta.CreateCollectionOption) error {
	if len(t) < 1 {
		return nil
	}
	db := p.DB
	if option != nil && len(option.Omit) > 0 {
		db = db.Omit(option.Omit...)
	}
	err := db.Create(t).Error
	pgErr, ok := err.(*pgconn.PgError)
	if ok && pgErr.Code == "23505" {
		err = myErrors.ErrNameDuplicate
	}
	return err
}
