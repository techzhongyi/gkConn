package helper

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
)

func TestValidCode(t *testing.T) {
	oriMsg := "232302FE4C464E41344C4441324A41583033343234010025170C1D0B261E0101030101260022DAD014AA29546301211388FFFF0500062EE5B90226890E01"
	ok := ValidCode(oriMsg)
	fmt.Println(ok)

}

func TestGetResponseMsg(t *testing.T) {
	msg := "232302FE4C464E41344C4441324A41583033343234010025170C1D0B261E0101030101260022DAD014AA29546301211388FFFF0500062EE5B90226890E01"
	reply := GetResponseMsg("01", msg)
	fmt.Println("应答消息: ", reply, len(reply))
	ok := ValidCode(reply)
	fmt.Println("验证：", ok)

}

func Test1(t *testing.T) {
	msg := "232302FE4C464E41344C4441324A41583033343234010025170C1D0B261E0101030101260022DAD014AA29546301211388FFFF0500062EE5B90226890E01"
	fmt.Println("起始符:", msg[0:4])
	fmt.Println("命令标识:", msg[4:6])
	fmt.Println("应答标志:", msg[6:8])
	fmt.Println("唯一识别码:", msg[8:42]) // 4C464E41344C4441324A41583033343234
	fmt.Println("加密方式:", msg[42:44])
	fmt.Println("数据单元长度:", msg[44:48])
	fmt.Println("数据单元:", msg[48:])

}

func Test2(t *testing.T) {
	// 16进制转字符串
	hexStr := "4C464E41344C4441324A41583033343234" // 十六进制字符串 "48656c6c6f"
	// 将十六进制字符串解码为字节切片
	bytes, err := hex.DecodeString(hexStr)
	if err != nil {
		fmt.Println("解码失败:", err)
		return
	}
	// 将字节切片转换为字符串
	str := string(bytes)
	fmt.Println(str)

	// 单个16进制转字符串
	hexStr2 := "34" // 十六进制字符串 "41"
	// 将十六进制字符串转换为整数值
	intValue, err := strconv.ParseInt(hexStr2, 16, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}
	fmt.Println(intValue)
	fmt.Println(fmt.Sprintf("%c", intValue))
}

func TestConv(t *testing.T) {
	hexStr := "0025" // 十六进制字符串 "1A"
	// 将十六进制字符串转换为十进制整数
	decimalValue, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		fmt.Println("转换失败:", err)
		return
	}
	fmt.Printf("16进制转化为10进制 %s -> %d\n", hexStr, decimalValue)
	// 十 进制转十六进制
	decimalValue = 6 // 十六进制字符串 "1A"
	s := strconv.FormatInt(int64(decimalValue), 16)
	fmt.Printf("10进制转化为16进制 %d -> %s\n", decimalValue, s)
}

func TestGetFactoryInfo(t *testing.T) {
	a, b := GetFactoryInfo("232302FE4C464E41344C4441324A41583033343234010025170C1D0B261E0101030101260022DAD014AA29546301211388FFFF0500062EE5B90226890E01")
	fmt.Printf(" a=%s, b=%s\n", a, b)
	fmt.Println(" xxxxx", a == "")
	fmt.Println(" yyyyy", b == "")

}
