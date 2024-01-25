package main

import (
	"errors"
	"fmt"
	"github.com/kardianos/osext"
	log "github.com/sirupsen/logrus"
	"github.com/techzhongyi/comlibgo/mhandler"
	"github.com/techzhongyi/comlibgo/util"
	gkCore "gk/src"
	"net/http"
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
	// 初始化redis
	gkCore.RedCacheDbs = util.InitRedis(gkCore.Confg.Redis.Host+":"+gkCore.Confg.Redis.Port,
		gkCore.Confg.Redis.Password, gkCore.Confg.Redis.Db["cache"])
	http.HandleFunc("/gkws", gkCore.HandlerWebsocket)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", gkCore.Confg.Server.Port), nil))
}
