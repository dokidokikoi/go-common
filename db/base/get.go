package base

import (
	"context"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (p *PgModel[T]) Get(ctx context.Context, t *T, option *meta.GetOption) (*T, error) {
	var result T
	db := p.DB
	var err error
	if option != nil {
		for _, s := range option.Preload {
			db = db.Preload(s)
		}
		if option.Include != nil {
			db = db.Where(t, option.Include)
		}
		if len(option.Select) > 0 {
			db = db.Select(option.Select)
		}
	}
	err = db.Where(t).First(&result).Error
	if err != nil {
		return nil, errorHandle(err)
	}

	return &result, nil
}
