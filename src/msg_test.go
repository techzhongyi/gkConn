// 消息处理

package gkCore

import (
	"testing"
)

// save2Redis 将数据存储到kafka
func TestSave2Kafka(t *testing.T) {
	InitKakfaProducer()
	Save2Kafka("aaaa", "bbbb")

}
