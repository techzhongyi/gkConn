package gkCore

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gkConn/src/helper"
	"net"
	"strconv"
	"time"
)

type MsgStat struct {
	msg  string
	stat string
}

var connStat = make(map[net.Conn]string)

// 延时函数创建连接后40分钟后触发，如果没有实现平台登录，关闭连接
func delayConn(conn net.Conn, ch chan MsgStat) {
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
func dealMsg(msg string, conn net.Conn, ch chan<- MsgStat) string {
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
				log.Info("<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< 登录 >>>>>>>>>>>>>>>>>>>>>>>>>", connStat[conn])
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

// 校验厂商的用户名和密码,如果成功返回厂商id.
func validFactory(msg string) string {
	name, password := helper.GetFactoryInfo(msg)
	log.Debug("factoryInfo ---> ", name, "--------", password)
	for _, v := range Confg.Factorys {
		log.Info("-------------- ", v.Name, "--------------", v.Password)
		if v.Name == name && v.Password == password {
			log.Info("--------------login success")
			return name
		}
	}
	log.Error("--------------login fail")
	return LoginErr

}

// 从ch读取消息发送给client。消息应答。如果应答标识是FE，需要应答；不是FE不需要应答
func write2Client(conn net.Conn, ch <-chan MsgStat) {
	for v := range ch {
		reply := helper.GetReply(v.msg)
		switch reply {
		case "FE":
			response := helper.GetResponseMsg(v.stat, v.msg)
			_, err := conn.Write([]byte(response))
			if err != nil {
				log.Error("$$$$$$$$$$$$ reply err", err)
			}
			log.Debug("@@@应答 reply success ", response)

		default:
			log.Debug("不是FE不需要应答")
		}
	}
	log.Debug("write2Client over.......")

}

// 从conn中解析出消息
func parseMsg(conn net.Conn) (_err error, _msg string) {
	defer func() {
		if r := recover(); r != nil {
			_err = fmt.Errorf("%s", r)
			_msg = "-4"
		}
	}()
	bufferHead := make([]byte, 48)
	_, err := conn.Read(bufferHead)
	if err != nil {
		return err, "-1"
	}
	dataLen, err := strconv.ParseInt(string(bufferHead)[len(string(bufferHead))-4:], 16, 64)
	if err != nil {
		return err, "-2"
	}
	bufferBody := make([]byte, dataLen+2)
	_, err = conn.Read(bufferBody)
	if err != nil {
		return err, "-3"
	}
	return nil, string(bufferHead) + string(bufferBody)
}

func HandleConnection(conn net.Conn) {
	log.Info("....................建立连接..................")
	defer func(conn net.Conn) {
		log.Info("@@@@@@@@@@@@@@@@@@@@ --- conn close --- @@@@@@@@@@@@@@@@@@@@")
		err := conn.Close()
		if err != nil {
		}
	}(conn)
	ch := make(chan MsgStat)
	time.AfterFunc(time.Duration(Confg.Server.DelayMinutes)*time.Minute, func() {
		log.Info("延时函数执行啦！")
		delayConn(conn, ch)
	})
	go write2Client(conn, ch)
	for {
		err, msg := parseMsg(conn)
		if err != nil {
			log.Error("------------- parseMsg err ", err, msg)
			connStat[conn] = ""
			return
		}
		log.Debug("##########################################################################################")
		log.Debug("@@@@received: ", len(msg), "----", msg)
		save2Redis(dealMsg(msg, conn, ch), msg)
	}
}
