package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
	"strconv"

	tool "haina.com/market/hqtools/printfiletools/util"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type DKLine struct {
	NSID     int32  // 证券ID
	NTime    int32  // 时间 unix time
	NPreCPx  int32  // 昨收价 * 10000
	NOpenPx  int32  // 开盘价 * 10000
	NHighPx  int32  // 最高价 * 10000
	NLowPx   int32  // 最低价 * 10000
	NLastPx  int32  // 最新价 * 10000
	LlVolume int64  // 成交量
	LlValue  int64  // 成交额 * 10000
	NAvgPx   uint32 // 平均价 * 10000

}

var STOCKSIZE int

func main() {

	args := os.Args //获取用户输入的所有参数
	if args == nil {
		logging.Info("args is nil")
		return
	}
	input := args[1] //获取输入的第一个参数
	var output = ""
	if len(args) > 2 {
		output = args[2] //获取输入的第二个参数
	}

	// 根据第一个参数打开文件
	file, err := tool.OpenFile(input + "/dk.dat")
	if err != nil {
		return
	}
	var dkline DKLine
	STOCKSIZE = binary.Size(&dkline)
	var dklines []DKLine
	for {
		des := make([]byte, STOCKSIZE)
		num, err := tool.ReadFiles(file, des)
		if err != nil {
			if err == io.EOF { //读到了文件末尾
				break
			}
			logging.Error("Read file error...%v", err.Error())
			return
		}
		if num < STOCKSIZE && 0 < num {
			logging.Error("Stock struct size error... or hqtools write file error")
			return
		}
		//todoing		des
		buffer := bytes.NewBuffer(des)
		binary.Read(buffer, binary.LittleEndian, &dkline)
		dklines = append(dklines, dkline)
	}

	file.Close()
	// 把取到的值打印到txt文件中
	var fileadd = "" // 输出地址
	if output != "" {
		// 检查目录如果没有创建
		lib.CheckDir(output)
		fileadd = output + "/dk.txt"
	} else {
		fileadd = input + "/dk.txt"
	}
	// 循环追加打印到文件
	file, err = os.OpenFile(fileadd, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
	if err != nil {
		logging.Info(" %v", err)
	}
	file.WriteString("证券ID\t\t时间 unix time\t昨收价 * 10000\t开盘价 * 10000\t最高价 * 10000\t最低价 * 10000\t最新价 * 10000\t成交量\t\t\t成交额 * 10000\t\t平均价 * 10000\n")
	for _, item := range dklines {
		file, err = os.OpenFile(fileadd, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0664)
		if err != nil {
			logging.Info(" %v", err)
		}
		var dklineinfo = ""
		dklineinfo += Int32tostr(item.NSID) + "\t" + Int32tostr(item.NTime) + "\t" + Int32tostr(item.NPreCPx) + "\t\t"
		dklineinfo += Int32tostr(item.NOpenPx) + "\t\t" + Int32tostr(item.NHighPx) + "\t\t" + Int32tostr(item.NLowPx) + "\t\t"
		dklineinfo += Int32tostr(item.NLastPx) + "\t\t" + Int64tostr(item.LlVolume) + "\t\t" + Int64tostr(item.LlValue) + "\t\t" + UInt32tostr(item.NAvgPx) + "\n"
		file.WriteString(dklineinfo)
	}
	file.Close()
	logging.Info("%v", file)
}

// int 32 转string
func Int32tostr(istr int32) string {
	return strconv.Itoa(int(istr))
}

// int 64 转string
func Int64tostr(istr int64) string {
	return strconv.Itoa(int(istr))
}

// uint 32 转string
func UInt32tostr(istr uint32) string {
	return strconv.Itoa(int(istr))
}
