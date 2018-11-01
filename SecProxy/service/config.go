package service

import "sync"

const (
	ProductStatusNormal = 0
	ProductSatusSaleOut = 1
	ProdcutStatusForceSaleOut = 2
)

type RedisConf struct {
	RedisAddr string
	RedisMaxIdle int
	RedisMaxActive int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr string
	EtcdTimeout int
	EtcdSecKeyPrefix string
	EtcdProductKey string
}

type SecKillConf struct {
	RedisConf RedisConf
	EtcdConf EtcdConf
	SecProductInfoMap map[int]*SecProductInfoConf
	RwSecProductLock sync.RWMutex
	LogPath string
	LogLevel string
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int64
	EndTime int64
	Status int
	Total int
	Left int
}

