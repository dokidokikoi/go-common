package meta

import "gorm.io/gorm/clause"

type CreateOption struct {
	Omit []string
}

type DoUpdates struct {
	Columns []clause.Column
	Updates clause.Set
}

type CreateCollectionOption struct {
	Omit      []string
	DoUpdates *DoUpdates
}
