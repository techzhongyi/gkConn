// 消息处理

package gkCore

import (
	log "github.com/sirupsen/logrus"
)

// save2Redis 将数据存储到redis
func save2Redis(factoryId, msg string) {
	if factoryId == "" || factoryId == LoginErr {
		return
	}
	if RedCacheDbs.LLen(factoryId).Val() > Confg.Redis.MaxLen {
		log.Error("Length exceeds maximum.......", RedCacheDbs.LLen(factoryId).Val())
		return
	}
	RedCacheDbs.LPush(factoryId, msg)
}
