package helper

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
)

func TestValidCode(t *testing.T) {
	oriMsg := "23230501111111111111111111111111111111111101000611111111111112"
	ok := ValidCode(oriMsg)
	fmt.Println(ok)

}

func TestGetResponseMsg(t *testing.T) {
	msg := "232305FE11111111111111111111111111111111110100521111111111112222676b7a796465657077617931737764646666676472353675393930686a6b666601a7"
	reply := GetResponseMsg("01", msg)
	fmt.Println("应答消息: ", reply, len(reply))
	ok := ValidCode(reply)
	fmt.Println("验证：", ok)

}

func Test1(t *testing.T) {
	msg := "232305FE11111111111111111111111111111111110100521111111111112222676b7a796465657077617931737764646666676472353675393930686a6b666601a7"
	fmt.Println("起始符:", msg[0:4])
	fmt.Println("命令标识:", msg[4:6])
	fmt.Println("应答标志:", msg[6:8])
	fmt.Println("唯一识别码:", msg[8:42]) // 4C464E41344C4441324A41583033343234
	fmt.Println("加密方式:", msg[42:44])
	fmt.Println("数据单元长度:", msg[44:48])
	fmt.Println("数据单元:", msg[48:])

}

func TestGenerateLogin(t *testing.T) {
	uk := ""
	for i := 0; i < 34; i++ {
		uk += fmt.Sprint(1)
	}
	time := "111111111111"
	no := "2222"
	name := "676b7a796465657077617931"
	pass := "737764646666676472353675393930686a6b6666"
	rule := "01"
	body := time + no + name + pass + rule
	msg := "2323" + "05" + "FE" + uk + "01" + "0052" + body + "**"
	print("vvvvvvv", len("2323"+"05"+"FE"+uk+"01"+"52"))
	rr := GetBccChecksum(msg)
	fmt.Println("rr=", rr)
	code := strconv.FormatInt(int64(rr), 16)
	fmt.Println("code=", code)
	if len(code) == 1 {
		code = "0" + code
	}
	v := msg[:len(msg)-2] + code
	fmt.Println(v)
	fmt.Println(len(v))
	name, pass = GetFactoryInfo(v)
	fmt.Println(name, " -- ", pass)
}

func TestConver(t *testing.T) {
	// 将字符串编码成十六进制格式
	ss := "gkzydeepway1" // 要转换的字符串
	hexStr := hex.EncodeToString([]byte(ss))
	fmt.Println("字符串转十六进制:", hexStr, " len=", len(hexStr))

	// 16进制转字符串
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Println("解码失败:", err)
		return
	}
	str := string(bytes)
	fmt.Println("十六进制转字符串:", str)

	//// 单个16进制转字符串
	//hexStr2 := "34" // 十六进制字符串 "41"
	//// 将十六进制字符串转换为整数值
	//intValue, err := strconv.ParseInt(hexStr2, 16, 64)
	//if err != nil {
	//	fmt.Println("转换失败:", err)
	//	return
	//}
	//fmt.Println(intValue)
	//fmt.Println(fmt.Sprintf("%c", intValue))
}

func TestConver2(t *testing.T) {
	// 10进制转16进制串
	decimal := 82
	hex := fmt.Sprintf("%X", decimal)
	fmt.Println(hex)

	// 16进制串转10进制
	xx, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}
	fmt.Println(xx)
}
