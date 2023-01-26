package base

import (
	"context"
	"fmt"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (p *PgModel[T]) Get(ctx context.Context, t *T, option *meta.GetOption) (*T, error) {
	var result T
	db := p.DB
	var err error
	if option != nil {
		for _, s := range option.Preload {
			db = db.Preload(fmt.Sprintf("%s", s[0]), s[1:]...)
		}
		if option.Include != nil {
			db = db.Where(t, option.Include)
		}
		err = db.Where(t).First(&result).Error
	} else {
		err = db.Where(t).First(&result).Error
	}

	return &result, err
}
