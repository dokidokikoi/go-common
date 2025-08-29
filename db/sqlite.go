package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SLConfigs struct {
	maxIdleConnections    int
	maxOpenConnections    int
	maxConnectionLifeTime time.Duration
}

func NewSqlite(database string, options ...OptionFunc) (*gorm.DB, error) {
	ops := Options{
		maxConnectionLifeTime: defaultMaxConnectionLifeTime,
		maxOpenConnections:    defaultMaxOpenConnections,
		maxIdleConnections:    defaultMaxIdleConnections,
	}

	for _, o := range options {
		o(&ops)
	}

	configs := SLConfigs{
		maxIdleConnections:    ops.maxIdleConnections,
		maxOpenConnections:    ops.maxOpenConnections,
		maxConnectionLifeTime: ops.maxConnectionLifeTime,
	}

	db, err := gorm.Open(sqlite.Open(database), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}
	// 设定数据库最大连接数
	sqlDB.SetMaxOpenConns(configs.maxOpenConnections)
	// 设定数据库最长连接时长
	sqlDB.SetConnMaxLifetime(configs.maxConnectionLifeTime)
	// 设定数据库最大空闲数
	sqlDB.SetMaxIdleConns(configs.maxIdleConnections)

	return db, nil

}
