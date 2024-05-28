package meta

import "gorm.io/gorm"

const (
	INNER_JOIN  = "INNER JOIN"
	LEFT_JOIN   = "LEFT JOIN"
	RIGHRT_JOIN = "RIGHT JOIN"
)

type GetOption struct {
	Include []string
	Preload []string
	Select  []string
	Join    []*Join
	Group   string
}

type Join struct {
	Method             string
	Table              string
	JoinTable          string
	InnerQueryAlias    string
	InnerQuery         *gorm.DB
	TableField         string
	JoinTableField     string
	JoinTableCondition []Condition
}
