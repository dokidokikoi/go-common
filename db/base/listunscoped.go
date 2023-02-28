package base

import (
	"context"
	"fmt"
	"strings"

	meta "github.com/dokidokikoi/go-common/meta/option"
)

func (p *PgModel[T]) CountUnscoped(ctx context.Context, t *T, option *meta.GetOption) (int64, error) {
	var result int64
	if option != nil && len(option.Include) > 0 {
		var fields []any
		for _, i := range option.Include {
			fields = append(fields, i)
		}
		err := p.DB.Model(t).Where(t, fields...).Unscoped().Count(&result).Error
		return result, err
	}
	err := p.DB.Model(t).Where(t).Unscoped().Count(&result).Error
	return result, err
}

func (p *PgModel[T]) CountComplexUnscoped(ctx context.Context, example *T, condition *meta.WhereNode, option *meta.GetOption) (int64, error) {
	var result int64
	var t T
	var db = p.DB.Model(&t)
	if option != nil {
		if option.Include == nil {
			db = db.Where(example)
		} else {
			var fields []any
			for _, i := range option.Include {
				fields = append(fields, i)
			}
			db = db.Where(example, fields...)
		}

		for _, join := range option.Join {
			joinSQL := fmt.Sprintf("%s %s ON %s.%s = %s.%s", join.Method, join.JoinTable, join.Table, join.TableField, join.JoinTable, join.JoinTableField)
			var joinConditions []string
			var values []any
			joinConditions = append(joinConditions, joinSQL)
			for _, condition := range join.JoinTableCondition {
				joinConditions = append(joinConditions, fmt.Sprintf("%s.%s %s ?", join.JoinTable, condition.Field, condition.Operator))
				values = append(values, condition.Value)
			}
			db.Joins(strings.Join(joinConditions, " AND "), values...)
		}
	}
	err := CompositeQuery(db, condition).Unscoped().Count(&result).Error
	return result, err
}

func (p *PgModel[T]) ListUnscoped(ctx context.Context, t *T, option *meta.ListOption) ([]*T, error) {
	var tList []*T
	err := CommonDealList(p.DB, t, option).Unscoped().Find(&tList).Error
	return tList, err
}

func (p *PgModel[T]) ListComplexUnscoped(ctx context.Context, example *T, condition *meta.WhereNode, option *meta.ListOption) ([]*T, error) {
	var tList []*T
	err := CompositeQuery(CommonDealList(p.DB, example, option), condition).Unscoped().Find(&tList).Error
	return tList, err
}
