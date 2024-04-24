// 消息处理

package gkCore

import (
	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// Save2Kafka 将数据存储到kafka
func Save2Kafka(factoryId, msg string) {
	if factoryId == "" || factoryId == LoginErr || msg == "" {
		log.Error("save2Kafka.......factoryId=", factoryId)
		return
	}
	head := msg[0:48]
	dataLen, _ := strconv.ParseInt(head[len(head)-4:], 16, 64)
	if dataLen == 0 {
		return
	}
	kMsg := &sarama.ProducerMessage{
		Topic: Confg.KafKa.Topic,
		Key:   sarama.StringEncoder(factoryId),
		Value: sarama.StringEncoder(msg),
	}
	KafkaProducer.Input() <- kMsg

}
