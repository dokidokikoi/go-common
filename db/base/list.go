package base

import (
	"context"
	"fmt"
	"strings"

	meta "github.com/dokidokikoi/go-common/meta/option"
	"gorm.io/gorm"
)

func (p *PgModel[T]) CountDB(ctx context.Context, t *T, option *meta.GetOption) *gorm.DB {
	if option != nil && len(option.Include) > 0 {
		var fields []any
		for _, i := range option.Include {
			fields = append(fields, i)
		}
		return p.DB.Model(t).Where(t, fields...)
	}
	return p.DB.Model(t).Where(t)
}

func (p *PgModel[T]) Count(ctx context.Context, t *T, option *meta.GetOption) (int64, error) {
	var result int64
	db := p.CountDB(ctx, t, option)
	err := db.Count(&result).Error
	return result, err
}

func (p *PgModel[T]) CountComplexDB(ctx context.Context, example *T, condition *meta.WhereNode, option *meta.GetOption) *gorm.DB {
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
	return CompositeQuery(db, condition)
}

func (p *PgModel[T]) CountComplex(ctx context.Context, example *T, condition *meta.WhereNode, option *meta.GetOption) (int64, error) {
	var result int64
	db := p.CountComplexDB(ctx, example, condition, option)
	err := db.Count(&result).Error
	return result, err
}

func (p *PgModel[T]) ListDB(ctx context.Context, t *T, option *meta.ListOption) *gorm.DB {
	return CommonDealList(p.DB, t, option)
}

func (p *PgModel[T]) List(ctx context.Context, t *T, option *meta.ListOption) ([]*T, error) {
	var tList []*T
	err := p.ListDB(ctx, t, option).Find(&tList).Error
	return tList, err
}

func (p *PgModel[T]) ListComplexDB(ctx context.Context, example *T, condition *meta.WhereNode, option *meta.ListOption) *gorm.DB {
	return CompositeQuery(CommonDealList(p.DB, example, option), condition)
}

func (p *PgModel[T]) ListComplex(ctx context.Context, example *T, condition *meta.WhereNode, option *meta.ListOption) ([]*T, error) {
	var tList []*T
	err := p.ListComplexDB(ctx, example, condition, option).Find(&tList).Error
	return tList, err
}
