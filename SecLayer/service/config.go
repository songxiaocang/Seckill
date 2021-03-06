package service

import (
	"github.com/garyburd/redigo/redis"
	"go.etcd.io/etcd/clientv3"
	"sync"
	"time"
)

var secLayerContext = &SecLayerContext{}

type RedisConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
	RedisQueueName   string
}

type EtcdConf struct {
	EtcdAddr          string
	EtcdTimeout       int
	EtcdSecKillPrefix string
	EtcdProductKey    string
}

type SecProductInfoConf struct {
	ProductId         int
	Start             int
	End               int
	Status            int
	Total             int
	Left              int
	OnePersonBuyLimit int
	BuyRate           float64
	SoldMaxLimit      int
	Seclimit          *SecLimit
}

type SecRequest struct {
	ProductId       int
	Source          string
	AuthCode        string
	AccessTime      time.Time
	Nance           string
	UserId          int
	UserAuthSign    string
	ClientAddr      string
	ClientReference string
	//CloseNotify string
	//ResChan *SecResponse
}

type SecResponse struct {
	ProductId int
	UserId    int
	Code      int
	Token     string
	TokenTime int64
}

type SecLayerConf struct {
	Proxy2LayerRedisConf RedisConf
	Layer2ProxyRedisConf RedisConf
	EtcdConfig           EtcdConf

	LogPath  string
	LogLevel string

	WriteGoroutineNum      int
	ReadGoroutineNum       int
	UserHandleGoroutineNum int
	Read2HandleChanSize    int
	Handle2WriteChanSize   int
	MaxRequestTimeout      int

	Send2WriteTimeout  int
	Send2HandleTimeout int

	SecProductInfoMap map[int]*SecProductInfoConf
	TokenPasswd       string
}

type SecLayerContext struct {
	Proxy2LayerRedisPool *redis.Pool
	Layer2ProxyRedisPool *redis.Pool
	EtcdClient           *clientv3.Client
	RwSecProductLock     sync.RWMutex

	WaitedGroup    sync.WaitGroup
	SecLayerConfig *SecLayerConf

	Read2HandleChan  chan *SecRequest
	Handle2WriteChan chan *SecResponse

	//,,,,
	UserHistoryMap    map[int]*UserHistory
	RwUserHistoryLock sync.RWMutex
	ProductcountMgr   *ProductCountMgr
}
