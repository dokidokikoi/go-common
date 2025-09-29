package base

import (
	"context"

	"gorm.io/gorm/clause"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (p *PgModel[T]) Create(ctx context.Context, t *T, option *meta.CreateOption) error {
	db := p.DB
	if option != nil && len(option.Omit) > 0 {
		db = db.Omit(option.Omit...)
	}

	return errorHandle(db.Create(t).Error)
}

func (p *PgModel[T]) Creates(ctx context.Context, ts []*T, option *meta.CreateCollectionOption) error {
	if len(ts) < 1 {
		return nil
	}
	db := p.DB
	if option != nil {
		if len(option.Omit) > 0 {
			db = db.Omit(option.Omit...)
		}
		if option.DoUpdates != nil {
			db = db.Clauses(handleOnUpdateOpeion(option.DoUpdates))
		}
	}

	return db.CreateInBatches(ts, 1000).Error
}

func handleOnUpdateOpeion(u *meta.DoUpdates) clause.OnConflict {
	noConflict := clause.OnConflict{
		DoNothing: u.DoNothing,
		UpdateAll: u.UpdateAll,
	}
	if len(u.Columns) > 0 {
		for _, column := range u.Columns {
			noConflict.Columns = append(noConflict.Columns, clause.Column{Name: column})
		}
	}
	if len(u.Updates) > 0 {
		if len(u.Values) == 0 {
			noConflict.DoUpdates = clause.AssignmentColumns(u.Updates)
		} else if len(u.Values) == len(u.Updates) {
			m := make(map[string]any)
			n := min(len(u.Updates), len(u.Values))
			for i := range n {
				m[u.Updates[i]] = u.Values[i]
			}
			noConflict.DoUpdates = clause.Assignments(m)
		}
	}
	return noConflict
}
