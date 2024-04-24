package main

import (
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	log "github.com/sirupsen/logrus"
	"github.com/techzhongyi/comlibgo/mhandler"
	gkCore "gkConn/src"
	"net"
	"path"
)

func main() {
	curDir, err := osext.ExecutableFolder()
	if err != nil {
		panic(errors.New(""))
	}
	gkCore.ParseConfig(path.Join(curDir, "config.yaml"))
	mhandler.InitLog(gkCore.Confg.Server.LogPath, gkCore.Confg.Server.Debug)
	log.Info(" 解析项目配置文件config + api... ")
	// 初始化kafka
	gkCore.InitKakfaProducer()
	// 监听地址和端口
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", gkCore.Confg.Server.Port))
	if err != nil {
		log.Error("!!!监听失败", err)
		return
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
		}
	}(listener)
	log.Info("服务器已启动，等待客户端连接.....")
	for {
		// 等待客户端连接
		conn, err := listener.Accept()
		if err != nil {
			log.Error("!!!客户端连接失败", err)
			continue
		}
		log.Info("客户端连接成功:", conn.RemoteAddr())
		// 启动一个goroutine处理客户端请求
		go gkCore.HandleConnection(conn)
	}

}
