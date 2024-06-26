package helper

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
)

func TestValidCode(t *testing.T) {
	oriMsg := "2323070167756f6b657a686979756e3030303030310100006c"
	ok := ValidCode(oriMsg)
	fmt.Println(ok)

}

func TestGetResponseMsg(t *testing.T) {
	msg := "232307fe67756f6b657a686979756e30303030303101000093"
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
	fmt.Println("vvvvv ", len(body))
	msg := "2323" + "05" + "FE" + uk + "01" + "0052" + body + "**"
	rr := GetBccChecksum(msg)
	code := strconv.FormatInt(int64(rr), 16)
	if len(code) == 1 {
		code = "0" + code
	}
	v := msg[:len(msg)-2] + code
	fmt.Println(len(v), "---", v)
	name, pass = GetFactoryInfo(v)
	fmt.Println(name, " -- ", pass)
}

func TestConver(t *testing.T) {
	// 将字符串编码成十六进制格式
	ss := "##\u0005�guokezhiyun00001\u0001\u0000\u0000(\u0002\u0004\u00111\u0000N�gkzydeepway1swddffgdr56u990hjkff\u0001" // 要转换的字符串
	hexStr := hex.EncodeToString([]byte(ss))
	fmt.Println("字符串转十六进制:", hexStr, " len=", len(hexStr), "---", len(ss))

	xx := "232305fe67756f6b657a686979756e30303030310100002802041131004efb676b7a796465657077617931737764646666676472353675393930686a6b66660104"
	//// 16进制转字符串
	bytes, err := hex.DecodeString(xx)
	if err != nil {
		fmt.Println("解码失败:", err)
		return
	}
	str := string(bytes)
	fmt.Println("十六进制转字符串:", str)
	hexStr2 := hex.EncodeToString(bytes)
	fmt.Println("ccccccc:", hexStr2)

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
	decimal := 78
	hex := fmt.Sprintf("%X", decimal)
	fmt.Println(hex)
	//02010e21354e7f676b7a796465657077617931737764646666676472353675393930686a6b666601
	// 16进制串转10进制 02010e21354e
	xx1, _ := strconv.ParseInt("02", 16, 64)
	fmt.Println(xx1)
	xx2, _ := strconv.ParseInt("01", 16, 64)
	fmt.Println(xx2)
	xx3, _ := strconv.ParseInt("0e", 16, 64)
	fmt.Println(xx3)
	xx4, _ := strconv.ParseInt("21", 16, 64)
	fmt.Println(xx4)
	xx5, _ := strconv.ParseInt("35", 16, 64)
	fmt.Println(xx5)
	xx6, _ := strconv.ParseInt("4E", 16, 64)
	fmt.Println(xx6)

	fmt.Println("-------------- ", xx1, xx2, xx3, xx4, xx5, xx6)

}

func TestMes(t *testing.T) {
	//msg := "232302fe4c4457544553543230323331313037313201002518041811370b01ffff010000ffffffff0000000000ff0f00000000050006cbba5b0157e54a9f"
	//fmt.Println(GetCommand(msg))
	head := "232302fe4c44575445535432303233313130373132010025"
	fmt.Println(len(head))
	dataLen, _ := strconv.ParseInt(head[len(head)-4:], 16, 64)
	fmt.Println(dataLen)

}
