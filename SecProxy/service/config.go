package service

import (
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

const (
	ProductStatusNormal       = 0
	ProductSatusSaleOut       = 1
	ProdcutStatusForceSaleOut = 2
)

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type EtcdConf struct {
	EtcdAddr         string
	EtcdTimeout      int
	EtcdSecKeyPrefix string
	EtcdProductKey   string
}

type SecKillConf struct {
	RedisBlackConf       RedisConf
	RedisProxy2LayerConf RedisConf
	RedisLayer2ProxyConf RedisConf
	EtcdConf             EtcdConf
	SecProductInfoMap    map[int]*SecProductInfoConf
	RwSecProductLock     sync.RWMutex
	LogPath              string
	LogLevel             string
	CookieSecretKey      string

	AccLimitConf   AccessLimitConf
	ReferWhiteList []string
	IpBlackMap     map[string]bool
	IdBlackMap     map[int]bool
	RwBlackLock    sync.RWMutex

	RedisBlackPool               *redis.Pool
	RedisProxy2LayerPool         *redis.Pool
	RedisLayer2ProxyPool         *redis.Pool
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

	AccLimitMgr *AccessLimitMgr

	ReqChan     chan *SecRequest
	ReqChanSize int

	UserConnMap       map[string]chan *SecResult
	RwUserConnMapLock sync.RWMutex
}

type AccessLimitConf struct {
	UserSecAccessLimit int
	IpSecAccessLimit   int
	UserMinAccessLimit int
	IpMinAccessLimit   int
}

type SecProductInfoConf struct {
	ProductId int
	StartTime int64
	EndTime   int64
	Status    int
	Total     int
	Left      int
}

type SecRequest struct {
	ProductId       int
	Source          string
	AuthCode        string
	SecTime         string
	Nance           string
	UserId          int
	UserAuthSign    string
	AccessTime      time.Time
	ClientAddr      string
	ClientReference string
	CloseNotify     <-chan bool     `json:"-"`
	ResultChan      chan *SecResult `json:"-"`
}

type SecResult struct {
	ProductId int
	UserId    int
	Code      int
	Token     string
}
