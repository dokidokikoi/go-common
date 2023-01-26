package config

import (
	"gorm.io/datatypes"
)

type StoreInfo struct {
	PGInfo    *PGConfig
	RedisInfo *RedisConfig
	MongoInfo *MongoConfig
	KafkaInfo *KafkaConfig
	OtherInfo datatypes.JSONMap
}
