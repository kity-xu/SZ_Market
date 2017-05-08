package kline

import (
	pbk "ProtocolBuffer/format/kline"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"haina.com/market/hqpost/config"

	tool "haina.com/market/hqpost/controllers"
	"haina.com/share/lib"
	"haina.com/share/store/redis"

	"github.com/golang/protobuf/proto"
	"haina.com/market/hqpost/models/tb_security"
	"haina.com/share/logging"
)

func (this *Security) DayLine(cfg *config.AppConfig, codes *[]*tb_security.SecurityCode) {
	var stock StockSingle
	var filename string
	STOCKSIZE := binary.Size(&stock)

	var seList []SingleSecurity

	/******************************沪深所有股票*************************************/
	for _, v := range *codes {
		var count int = 0   //股票计数器
		var exchange string //股票交易所

		//PB
		var klist pbk.KInfoTable

		//History of Single-Security
		var sigList SingleSecurity
		var date []int32
		week := make(map[int32]StockSingle)

		if v.SID/100000000 == 1 { //ascii 字符
			exchange = SH
		} else if v.SID/100000000 == 2 {
			exchange = SZ
		} else {
			logging.Error("%s", "Invalid file name...")
			return
		}
		filename = fmt.Sprintf("%s%s%d/%s", cfg.File.Path, exchange, v.SID, cfg.File.DKName)

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
			var kdata pbk.KInfo //pb类型

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
				logging.Error("StockSingle struct size error... or hqtools write file error")
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
			date = append(date, stock.Time)
			week[stock.Time] = stock

			klist.List = append(klist.List, &kdata)
			count++
		}
		file.Close()

		//入PB
		data, err := proto.Marshal(&klist)
		if err != nil {
			logging.Error("Encode protocbuf of day Line error...%v", err.Error())
			return
		}

		key := fmt.Sprintf(REDISKEY_SECURITY_HDAY, v.SID)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}

		sigList.Sid = v.SID
		sigList.Date = date
		sigList.SigStock = week

		seList = append(seList, sigList)

		/*-------------------------------------end------------------------------------*/

		//logging.Debug("The historical data of each stock number:%v", count)
	}

	this.list.Securitys = &seList
	/*-----------------------------------------end----------------------------------*/
}
