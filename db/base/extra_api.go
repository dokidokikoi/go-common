package base

import (
	"fmt"
	"strings"

	meta "github.com/dokidokikoi/go-common/meta/option"

	"gorm.io/gorm"
)

func CompositeQuery(db *gorm.DB, chainRoot *meta.WhereNode) (tx *gorm.DB) {
	if chainRoot == nil {
		tx = db
		return
	}
	s, v := getWhereSql(chainRoot)
	return db.Where(s, v...)
}

func getWhereSql(node *meta.WhereNode) (string, []any) {
	if node == nil {
		return "", []any{}
	}
	var orArray []string
	var values []any
	result := strings.Builder{}
	for _, v := range node.Conditions {
		orArray = append(orArray, fmt.Sprintf("%s %s ?", v.Field, v.Operator))
		values = append(values, v.Value)
	}
	for _, v := range node.Nodes {
		s, vs := getWhereSql(v)
		orArray = append(orArray, s)
		values = append(values, vs...)
	}
	next, nextV := getWhereSql(node.Next)
	if len(orArray) > 0 {
		result.WriteString("(")
		result.WriteString(strings.Join(orArray, " OR "))
		result.WriteString(")")
	}
	if next != "" {
		if result.Len() > 0 {
			result.WriteString(" AND ")
		}
		result.WriteString(next)
		values = append(values, nextV...)
	}
	return result.String(), values

}

func CommonDeal(db *gorm.DB, example interface{}, option *meta.GetOption) (tx *gorm.DB) {
	if option != nil {
		if option.Preload != nil {
			for _, s := range option.Preload {
				db = db.Preload(s)
			}
		}
		for _, join := range option.Join {
			joinSQL := fmt.Sprintf("%s %s ON %s.%s = %s.%s", join.Method, join.JoinTable, join.Table, join.TableField, join.JoinTable, join.JoinTableField)
			var values []any
			if join.InnerQueryAlias != "" {
				joinSQL = fmt.Sprintf("%s (?) %s ON %s.%s = %s.%s",
					join.Method,
					join.InnerQueryAlias,
					join.Table,
					join.TableField,
					join.InnerQueryAlias,
					join.JoinTableField,
				)
				values = append(values, join.InnerQuery)
			}
			joinConditions := []string{joinSQL}
			for _, condition := range join.JoinTableCondition {
				joinConditions = append(joinConditions, fmt.Sprintf("%s.%s %s ?", join.JoinTable, condition.Field, condition.Operator))
				values = append(values, condition.Value)
			}
			db = db.Joins(strings.Join(joinConditions, " AND "), values...)
		}
		if option.Include != nil {
			var fields []any
			for _, i := range option.Include {
				fields = append(fields, i)
			}
			db = db.Where(example, fields...)
		}
		if option.Group != "" {
			db.Group(option.Group)
		}
		if len(option.Select) > 0 {
			db.Select(option.Select)
		}
		tx = db
		return
	}
	tx = db.Where(example)
	return
}
