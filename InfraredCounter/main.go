package main

import (
	"log"
	"os"
	"encoding/csv"
	"fmt"
	"golang.org/x/text/transform"
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io/ioutil"
)

// 50mm*30mm标签打印模板
var (
	TAG5030 = "SIZE 50 mm,30 mm\r\n" +
		"DENSITY 15\r\n" +
		"DIRECTION 1\r\n" +
		"CLS\r\n" +
		"TEXT 12,40,\"TSS24.BF2\",0,1,1,\" 线路宝\"\r\n" +
		"TEXT 0,80,\"TSS24.BF2\",0,1,1,\"PCBBAO.COM\"\r\n" +
		"TEXT 280,40,\"TSS24.BF2\",0,1,1,\"物联网平台\"\r\n" +
		"QRCODE 120,20,M,3,A,0,\"https://pcbapi.pcbbao.com/mesdemo/index.php/v1/default/device-detail?product_key=%v&device_name=%v\"\r\n" +
		"TEXT 20,180,\"TSS24.BF2\",0,1,1,\"设备名称: %v\"\r\n" +
		"TEXT 20,210,\"TSS24.BF2\",0,1,1,\"设备编号: %v\"\r\n" +
		"PRINT 1,1\r\n" +
		"\r\n"
)

// 日期: 2019-03-25
// 程序说明: 用于打印红外计数器的标签
// 作者: @zjt
func main() {
	// 1.遍历CSV文档
	file, err := os.Open("./InfraredCounter.csv")
	defer file.Close()
	CheckErr(err)
	reader := csv.NewReader(file)
	// 2.遍历item
	strDatas, err := reader.ReadAll()
	CheckErr(err)
	var printResult string
	if len(strDatas) > 0 {
		for key, value := range strDatas {
			if key == 0 {
				continue
			}
			device := Device{
				ProductKey:  value[1],
				ChineseName: value[0],
				EnglishName: value[2],
			}
			tmp := fmt.Sprintf(TAG5030, device.ProductKey, device.EnglishName, device.ChineseName, device.EnglishName)
			printResult += tmp
		}
	}
	// 3.生成内容
	printfile, err := os.Create("printer.prn")
	defer printfile.Close()
	CheckErr(err)
	// 转码:从utf8转到GBK
	encodeStr, err := Utf8ToGbk([]byte(printResult))
	CheckErr(err)
	lenth, err := printfile.WriteString(string(encodeStr))
	CheckErr(err)
	fmt.Printf("打印长度:%v\r\n", lenth)
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func CheckErr(err error) {
	if nil != err {
		log.Fatal(err)
	}
}

type Device struct {
	ProductKey  string // 产品名称
	ChineseName string // 中文名
	EnglishName string // 英文名
}
