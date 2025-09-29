package base

import (
	"context"

	meta "github.com/dokidokikoi/go-common/meta/option"
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

	return errorHandle(db.Updates(t).Error)
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

	return errorHandle(result.Error)
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
	return errorHandle(db.Save(t).Error)
}
