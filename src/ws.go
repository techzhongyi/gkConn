package gkCore

import (
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"gkConn/src/helper"
	"net/http"
	"time"
	"unsafe"
)

type MsgStat struct {
	msg  string
	stat string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var connStat = make(map[*websocket.Conn]string)

// 延时函数创建连接后40分钟后触发，如果没有实现平台登录，关闭连接
func delayConn(conn *websocket.Conn, ch chan MsgStat) {
	log.Debug("#######################################  delayConn")
	if connStat[conn] == "" || connStat[conn] == LoginErr {
		log.Info("....................断开连接..................")
		closeChan(ch)
		err := conn.Close()
		if err != nil {
			log.Error("....................断开连接异常..................", err)
		}
	}
}

func closeChan(ch chan MsgStat) {
	defer func() {
		if r := recover(); r != nil {
			log.Debug("delayConn recover", r)
		}
	}()
	close(ch)

}

// 消息处理
func dealMsg(msg string, conn *websocket.Conn, ch chan<- MsgStat) string {
	// 1 如果不是32960协议不处理
	if helper.Is32960(msg) == false {
		return ""
	}
	// 2 处理登录登出重置状态
	responseStat := MsgSuccess
	if helper.ValidCode(msg) == false {
		responseStat = MsgFail
	} else {
		if connStat[conn] == "" || connStat[conn] == LoginErr {
			// 平台首次登录
			if helper.IsLogin(msg) {
				connStat[conn] = validFactory(msg)
				log.Info("######################################### 登录 #########################################", connStat[conn])
				if connStat[conn] == LoginErr {
					responseStat = MsgFail
				}
			}
		} else {
			// 登出操作重置连接状态
			if helper.IsSignOut(msg) {
				log.Info(".######################################### 登出 #########################################", connStat[conn])
				// 连接状态依然有效，客户端如果主动关闭连接再置状态
				//connStat[conn] = ""
			}
		}
	}
	// 3 发送消息写数据到客户端
	ch <- MsgStat{
		msg:  msg,
		stat: responseStat,
	}
	return connStat[conn]
}

// 校验厂商的用户名和密码,如果成功返回厂商id
func validFactory(msg string) string {
	name, password := helper.GetFactoryInfo(msg)
	for _, v := range Confg.Factorys {
		if v.Name == name && v.Password == password {
			return name
		}
	}
	return LoginErr

}

// 从ch读取消息发送给client。消息应答。如果应答标识是FE，需要应答；不是FE不需要应答
func write2Client(conn *websocket.Conn, ch <-chan MsgStat) {
	for v := range ch {
		reply := helper.GetReply(v.msg)
		switch reply {
		case "FE":
			response := helper.GetResponseMsg(v.stat, v.msg)
			err := conn.WriteMessage(1, *(*[]byte)(unsafe.Pointer(&response)))
			if err != nil {
				log.Error("$$$$$$$$$$$$ conn.WriteMessage err", err)
			}
			log.Debug("...................需要应答 write success ", response)

		default:
			log.Debug("不是FE不需要应答")
		}
	}
	log.Debug("write2Client over.......")

}

func HandlerWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Info(err)
		return
	}
	log.Info("....................建立连接..................")
	ch := make(chan MsgStat)
	defer func(conn *websocket.Conn) {
		closeChan(ch)
		connStat[conn] = ""
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	time.AfterFunc(time.Duration(Confg.Server.DelayMinutes)*time.Minute, func() {
		log.Info("延时函数执行啦！")
		delayConn(conn, ch)
	})
	go write2Client(conn, ch)
	for {
		messageType, p, err := conn.ReadMessage()
		log.Debug("@@@@@ ", messageType, *(*string)(unsafe.Pointer(&p)))
		if err != nil {
			log.Error("--conn.ReadMessage err ", err)
			connStat[conn] = ""
			err := conn.Close()
			if err != nil {
			}
			return
		}
		msg := *(*string)(unsafe.Pointer(&p))
		log.Debug("===========================================================================================")
		log.Debug(".............received: ", msg)
		save2Redis(dealMsg(msg, conn, ch), msg)
	}

}
