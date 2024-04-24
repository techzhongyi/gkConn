package gkCore

import (
	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
	comslibgo "github.com/techzhongyi/comlibgo"
)

var Confg *Config
var KafkaProducer sarama.AsyncProducer

const (
	LoginErr   = "-1" // 厂商登录失败
	MsgSuccess = "01" // 消息校验成功
	MsgFail    = "02" // 消息校验失败

	ReadHeadErr = "-1"   // 读消息头报错
	ProtocolErr = "-2"   // 协议错误
	MsgParseErr = "-3"   // 消息解析错误
	ReadBodyErr = "-4"   // 读消息体报错
	UnknowErr   = "-100" // 未知错误

)

type Config struct {
	Server   Server
	KafKa    KafKa
	Factorys []Factory
}

type KafKa struct {
	Uri   string
	Topic string
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
	DelayMinutes int64 // 延时函数执行时间
}

// ParseConfig 解析配置文件 分为两部分 解析conf.yaml 与 interface.yaml
func ParseConfig(pathConfig string) {
	var conf Config
	comslibgo.Parse4Yaml(pathConfig, &conf)
	Confg = &conf
}

func InitKakfaProducer() {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	_KafkaProducer, err := sarama.NewAsyncProducer([]string{Confg.KafKa.Uri}, config)
	if err != nil {
		return
	}
	KafkaProducer = _KafkaProducer
	go func() {
		for {
			select {
			case suc := <-KafkaProducer.Successes():
				log.WithFields(log.Fields{
					"Offset":    suc.Offset,
					"Partition": suc.Partition,
					"Timestamp": suc.Timestamp,
					"Key":       suc.Key,
					"Value":     suc.Value,
				}).Info("Kafka success info  ")
			case fail := <-KafkaProducer.Errors():
				log.Error("Failed to send message ", fail.Err.Error())
			}
		}
	}()
}
