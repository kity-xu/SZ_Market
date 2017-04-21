package market

import (
	"ProtocolBuffer/format/redis/pbdef/kline"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"
	"time"

	"haina.com/share/store/redis"

	"haina.com/share/lib"

	"github.com/golang/protobuf/proto"

	"haina.com/market/hqpost/config"
	tool "haina.com/market/hqpost/controllers"
	"haina.com/market/hqpost/models/tb_stokcode"
	"haina.com/share/logging"
)

type Security struct {
}

type Stock struct {
	SID    int32  // 证券ID
	Time   int32  // 时间 unix time
	PreCPx int32  // 昨收价 * 10000
	OpenPx int32  // 开盘价 * 10000
	HighPx int32  // 最高价 * 10000
	LowPx  int32  // 最低价 * 10000
	LastPx int32  // 最新价 * 10000
	Volume int64  // 成交量
	Value  int64  // 成交额 * 10000
	AvgPx  uint32 // 平均价 * 10000

}

var STOCKSIZE int

//主
func (this *Security) UpdateMarket(cfg *config.AppConfig) {

	codes, err := this.getSecurityTable()
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}
	ReadFilesAndToPB(cfg, codes)
}

func ReadFilesAndToPB(cfg *config.AppConfig, codes []tb_stokcode.Code) {
	//开始时间
	start := time.Now()

	var stock Stock
	var filename string
	STOCKSIZE = binary.Size(&stock)

	/******************************沪深所有股票*************************************/
	for _, v := range codes {
		var count int = 0 //股票计数器
		var klist kline.DKInfoTable
		var kReply kline.ReplyDKInfoTable

		sid := strconv.Itoa(int(v.SID))

		if sid[0] == '1' { //ascii 字符
			filename = cfg.File.Path + tool.SH + sid + "/" + cfg.File.DKName
		} else if sid[0] == '2' {
			filename = cfg.File.Path + tool.SZ + sid + "/" + cfg.File.DKName
		}

		if !lib.IsFileExist(filename) {
			logging.Debug("File does not exist...%s", filename)
			continue
		}

		file, err := tool.OpenFile(filename)
		if err != nil {
			return
		}

		/*************************每只股票的历史信息（日K线）*****************************/
		for {
			var kdata kline.DKInfo //pb类型
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
			binary.Read(buffer, binary.LittleEndian, &stock)

			//stock 转pb格式
			kdata.NSID = stock.SID
			kdata.NTime = stock.Time
			kdata.NPreCPx = stock.PreCPx
			kdata.NOpenPx = stock.OpenPx
			kdata.NHighPx = stock.HighPx
			kdata.NLowPx = stock.LowPx
			kdata.NLastPx = stock.LastPx
			kdata.LlVolume = stock.Volume
			kdata.LlValue = stock.Value
			kdata.NAvgPx = stock.AvgPx
			//logging.Debug("------------stock:%v-----------", stock)

			klist.List = append(klist.List, &kdata)
			count++
		}
		file.Close()

		//入PB
		kReply.Code = 200
		kReply.Dktable = &klist

		data, err := proto.Marshal(&kReply)
		if err != nil {
			logging.Error("Encode protoc buf error...%v", err.Error())
			return
		}

		key := fmt.Sprintf(tool.KEY_KLINE, sid)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}

		/*-------------------------------------end------------------------------------*/

		//logging.Debug("The historical data of each stock number:%v", count)
	}
	/*-----------------------------------------end----------------------------------*/
	//结束时间
	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))
}

func (this *Security) getSecurityTable() ([]tb_stokcode.Code, error) {
	return tb_stokcode.GetSecurityTableFromMG()
}
