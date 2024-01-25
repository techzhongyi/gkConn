package gkCore

import (
	"fmt"
	"github.com/go-redis/redis"
	comslibgo "github.com/techzhongyi/comlibgo"
)

var RedCacheDbs *redis.Client
var Confg *Config

const (
	LoginErr   = "-1" // 厂商登录失败
	MsgSuccess = "01" // 消息校验成功
	MsgFail    = "02" // 消息校验失败
)

type Config struct {
	Server   Server
	Redis    *Redis
	Factorys []Factory
}

type Redis struct {
	Host     string
	Port     string
	Password string
	Db       map[string]int
}

type Factory struct {
	Name     string
	Password string
}

type Server struct {
	Name         string
	Debug        bool
	Port         string
	LogPath      string
	MaxLen       int64 // 消息最大长度
	DelayMinutes int64 // 延时函数执行时间
}

// String 目的很单一，打印Config 的时候 展示对应信息，而不是地址x
func (redis *Redis) String() string {
	return fmt.Sprintf("Redis{Host:%s, Port:%s, Password:%s, Db:%v}",
		redis.Host, redis.Port, redis.Password, redis.Db)

}

// ParseConfig 解析配置文件 分为两部分 解析conf.yaml 与 interface.yaml
func ParseConfig(pathConfig string) {
	var conf Config
	comslibgo.Parse4Yaml(pathConfig, &conf)
	Confg = &conf
}
