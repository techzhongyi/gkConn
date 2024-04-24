// 消息处理

package gkCore

import (
	"github.com/IBM/sarama"
	log "github.com/sirupsen/logrus"
)

// Save2Kafka 将数据存储到kafka
func Save2Kafka(factoryId, msg string) {
	if factoryId == "" || factoryId == LoginErr {
		log.Error("save2Kafka.......factoryId=", factoryId)
		return
	}
	kMsg := &sarama.ProducerMessage{
		Topic: Confg.KafKa.Topic,
		Key:   sarama.StringEncoder(factoryId),
		Value: sarama.StringEncoder(msg),
	}
	KafkaProducer.Input() <- kMsg

}
