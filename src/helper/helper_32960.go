package helper

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// Is32960 是否是32960协议.
func Is32960(msg string) bool {
	// 校验协议头
	if strings.HasPrefix(msg, "2323") == false {
		return false
	}
	return true
}

// ValidCode 校验校验码
func ValidCode(msg string) bool {
	code := msg[len(msg)-2:]
	iCode, _ := strconv.ParseUint("0x"+code, 0, 8)
	bcc := GetBccChecksum(msg)
	if bcc == byte(iCode) {
		return true
	} else {
		log.Error("valid code err ...", bcc, byte(iCode))
		return false
	}

}

// GetBccChecksum 根据原始串计算bcc校验码
func GetBccChecksum(msg string) byte {
	msg_ := msg[4 : len(msg)-2]
	l := len(msg_) / 2
	var arr []byte
	d := 0
	for i := 0; i < l; i++ {
		v, _ := strconv.ParseUint("0x"+msg_[d:d+2], 0, 8)
		arr = append(arr, byte(v))
		d += 2
	}

	bcc := BccChecksum(arr)
	return bcc
}

// IsLogin 是否是登录
func IsLogin(msg string) bool {
	return GetCommand(msg) == "05"
}

// IsSignOut 是否是登出
func IsSignOut(msg string) bool {
	return GetCommand(msg) == "06"
}

// GetCommand 获取命令标识
func GetCommand(msg string) string {
	defer func() {
		if r := recover(); r != nil {
			log.Error("GetCommand ", r)
		}
	}()
	return msg[4:6]
}

// GetReply 获取应答标识
func GetReply(msg string) string {
	return strings.ToUpper(msg[6:8])
}

// GetFactoryInfo 解析获取用户名和密码
func GetFactoryInfo(msg string) (string, string) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("GetFactoryInfo", r)
		}
	}()
	if IsLogin(msg) == false {
		return "", ""
	}
	body := msg[48 : len(msg)-2]
	name, _ := hex.DecodeString(body[len(body)-66 : len(body)-42])
	password, _ := hex.DecodeString(body[len(body)-42 : len(body)-2])
	return string(name), string(password)

}

// GetResponseMsg 返回应答消息
func GetResponseMsg(responseStat string, msg string) string {
	//          起始符2 + 命令标识1 + 应答标志1 + 唯一标识码17 + 加密方式1 + 数据单元长度2 + 数据单元 + 校验码1
	head := msg[0:48]
	dataLen, _ := strconv.ParseInt(head[len(head)-4:], 16, 64)
	newMsg := ""
	// 如果数据单元长度是0 是心跳包 回复也是数据长度为0
	if dataLen == 0 {
		newMsg = msg[0:4] + msg[4:6] + responseStat + msg[8:42] + msg[42:44] + "0000" + "**"
	} else {
		newMsg = msg[0:4] + msg[4:6] + responseStat + msg[8:42] + msg[42:44] + "0006" + msg[48:60] + "**"
	}
	rr := GetBccChecksum(newMsg)
	code := strconv.FormatUint(uint64(rr), 16)
	if len(code) == 1 {
		code = "0" + code
	}
	return newMsg[:len(newMsg)-2] + code

}
