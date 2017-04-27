package market

import (
	"ProtocolBuffer/format/redis/pbdef/kline"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"
	"time"

	"haina.com/share/store/redis"

	"haina.com/share/lib"

	"github.com/golang/protobuf/proto"

	"haina.com/market/hqpost/config"
	tool "haina.com/market/hqpost/controllers"
	"haina.com/market/hqpost/models/tb_stokcode"
	"haina.com/share/logging"
)

var (
	cfg   *config.AppConfig
	codes *[]*tb_stokcode.Code
)

type Security struct {
	week WeekSecurityList
}

//单个股票
type WeekSecurity struct {
	Sid      int32                 //股票SID
	Date     []int32               //单个股票的历史日期
	SigStock map[int32]StockSingle //单个股票的历史数据
	WeekDays *[][]int32            //单个股票的周天
}

//所有股票
type WeekSecurityList struct {
	Securitys *[]WeekSecurity
}

//个股信息
type StockSingle struct {
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

func Update(fg *config.AppConfig) {
	cs, err := getSecurityTable()
	if err != nil {
		logging.Error("%v", err.Error())
		return
	}
	cfg = fg
	codes = cs
	security := new(Security)

	/*********************开始时间************************/
	start := time.Now()
	//--------------------------------------------------/
	security.KDayLine()
	security.GetAllSecurityDayList()

	security.KWeekLine()
	//--------------------------------------------------/
	end := time.Now()
	logging.Info("Update Kline historical data successed, and running time:%v", end.Sub(start))
	/*********************结束时间***********************/

}

func (this *Security) KDayLine() {
	var stock StockSingle
	var filename string
	STOCKSIZE := binary.Size(&stock)

	var seList []WeekSecurity

	/******************************沪深所有股票*************************************/
	for _, v := range *codes {
		var count int = 0   //股票计数器
		var exchange string //股票交易所

		//PB
		var klist kline.KInfoTable
		var kReply kline.ReplyKInfoTable

		//History of Single-Security
		var sigList WeekSecurity
		var date []int32
		week := make(map[int32]StockSingle)

		if v.SID/100000000 == 1 { //ascii 字符
			exchange = tool.SH
		} else if v.SID/100000000 == 2 {
			exchange = tool.SZ
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
			var kdata kline.KInfo //pb类型

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
		kReply.Code = 200
		kReply.Ktable = &klist

		data, err := proto.Marshal(&kReply)
		if err != nil {
			logging.Error("Encode protocbuf of day Line error...%v", err.Error())
			return
		}

		key := fmt.Sprintf(tool.REDISKEY_SECURITY_HDAY, v.SID)
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

	this.week.Securitys = &seList
	/*-----------------------------------------end----------------------------------*/
}

func (this *Security) KWeekLine() {
	securitys := *this.week.Securitys
	//logging.Debug("-kw-----%v", securitys[0].SigStock)

	for _, single := range securitys { // 以sid分类的单个股票
		//logging.Debug("%s:date:", single.Sid, *single.WeekDays)	//得到了该支股票的所有历史周天
		//logging.Debug("SID:%v", single.Sid)
		var tmps []StockSingle

		//PB
		var klist kline.KInfoTable
		var kReply kline.ReplyKInfoTable

		for _, week := range *single.WeekDays { //每一周
			//logging.Debug("Week:%v", week)
			tmp := StockSingle{}
			var mdata kline.KInfo //pb类型

			var (
				i          int
				day        int32
				AvgPxTotal uint32
			)
			for i, day = range week { //每一天
				//logging.Debug("day:%v---single.SigStock[day]:%v", day, single.SigStock[day])
				stockday := single.SigStock[day]
				if tmp.HighPx < stockday.HighPx || tmp.HighPx == 0 { //最高价
					tmp.HighPx = stockday.HighPx
				}
				if tmp.LowPx > stockday.LowPx || tmp.LowPx == 0 { //最低价
					tmp.LowPx = stockday.LowPx
				}
				tmp.Volume += stockday.Volume //成交量
				tmp.Value += stockday.Value   //成交额
				AvgPxTotal += stockday.AvgPx
			}
			tmp.SID = single.Sid
			tmp.Time = single.SigStock[week[0]].Time     //时间（取每周第一天）
			tmp.OpenPx = single.SigStock[week[0]].OpenPx //开盘价（每周第一天的开盘价）
			if len(tmps) > 0 {
				tmp.PreCPx = tmps[len(tmps)-1].LastPx //昨收价(上周的最新价)
			} else {
				tmp.PreCPx = 0
			}
			tmp.LastPx = single.SigStock[week[i]].LastPx //最新价
			if i > 0 {
				tmp.AvgPx = AvgPxTotal / uint32(i+1) //平均价
			} else if i == 0 {
				tmp.AvgPx = AvgPxTotal
			}
			//tmps = append(tmps, tmp)
			//logging.Debug("周线是:%v", tmps)

			//入PB
			mdata.NSID = tmp.SID
			mdata.NTime = tmp.Time
			mdata.NPreCPx = tmp.PreCPx
			mdata.NOpenPx = tmp.OpenPx
			mdata.NHighPx = tmp.HighPx
			mdata.NLowPx = tmp.LowPx
			mdata.NLastPx = tmp.LastPx
			mdata.LlVolume = tmp.Volume
			mdata.LlValue = tmp.Value
			mdata.NAvgPx = tmp.AvgPx

			klist.List = append(klist.List, &mdata)

		}
		//入PB
		kReply.Code = 200
		kReply.Ktable = &klist

		//入redis
		data, err := proto.Marshal(&kReply)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		key := fmt.Sprintf(tool.REDISKEY_SECURITY_HWEEK, single.Sid)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}

	}

}

func (this *Security) GetAllSecurityDayList() {
	secs := *this.week.Securitys

	for i, v := range secs {
		var wday [][]int32
		sat := DateAdd(int(v.Date[0])) //该股票第一个交易日所在周的周六

		var dates []int32
		for _, date := range v.Date {
			if IntToTime(int(date)).Before(sat) {
				dates = append(dates, date)
			} else {
				//logging.Debug("------一周的日期是：%v------", dates) //it's here
				wday = append(wday, dates)

				//				logging.Debug("------一周的日期完成------")
				//				logging.Debug("----------当前日期----%v---", date)
				sat = DateAdd(int(date))
				dates = nil
				dates = append(dates, date)
			}
		}
		//logging.Debug("------一周的日期完成-%v-----", wday)
		wday = append(wday, dates)
		secs[i].WeekDays = &wday

	}

	//logging.Debug("-----单个股票，所有周天：%v------", (*this.week.Securitys)[0].WeekDays)
	//logging.Debug("-----单个股票，secs[0].date：%v------", (*this.week.Securitys)[0].Date)

}

func IntToTime(date int) time.Time {
	swap := date % 10000
	year := date / 10000
	month := swap / 100
	day := swap % 100
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func DateAdd(date int) time.Time {
	var sat time.Time
	swap := date % 10000
	year := date / 10000
	month := swap / 100
	day := swap % 100

	//logging.Debug("%d-%d-%d", year, month, day)

	baseTime := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	weekday := baseTime.Weekday().String()

	var basedate string
	if strings.EqualFold(weekday, "Monday") {
		basedate = fmt.Sprintf("%d%s", 24*5, "h")

	} else if strings.EqualFold(weekday, "Tuesday") {
		basedate = fmt.Sprintf("%d%s", 24*4, "h")

	} else if strings.EqualFold(weekday, "Wednesday") {
		basedate = fmt.Sprintf("%d%s", 24*3, "h")

	} else if strings.EqualFold(weekday, "Thursday") {
		basedate = fmt.Sprintf("%d%s", 24*2, "h")

	} else if strings.EqualFold(weekday, "Friday") {
		basedate = fmt.Sprintf("%d%s", 24*1, "h")

	} else {
		logging.Error("Invalid trade date...")
		return sat
	}

	dd, _ := time.ParseDuration(basedate)
	sat = baseTime.Add(dd) //Saturday（星期六）
	return sat
}

func getSecurityTable() (*[]*tb_stokcode.Code, error) {
	return tb_stokcode.GetSecurityTableFromMG()
}
