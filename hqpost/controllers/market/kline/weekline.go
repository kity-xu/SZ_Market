package kline

import (
	pbk "ProtocolBuffer/format/kline"
	"fmt"

	"haina.com/share/store/redis"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
)

func (this *Security) WeekLine() {
	this.GetAllSecurityDayList()
	securitys := *this.list.Securitys
	//logging.Debug("-kw-----%v", securitys[0].SigStock)

	for _, single := range securitys { // 以sid分类的单个股票
		//logging.Debug("%s:date:", single.Sid, *single.WeekDays)	//得到了该支股票的所有历史周天
		//logging.Debug("SID:%v", single.Sid)
		var tmps []StockSingle

		//PB
		var klist pbk.KInfoTable

		for _, week := range *single.WeekDays { //每一周
			//logging.Debug("Week:%v", week)
			tmp := StockSingle{}
			var mdata pbk.KInfo //pb类型

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
			tmp.AvgPx = AvgPxTotal / uint32(i+1)         //平均价
			tmps = append(tmps, tmp)
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

		//入PB 入redis
		data, err := proto.Marshal(&klist)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		key := fmt.Sprintf(REDISKEY_SECURITY_HWEEK, single.Sid)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}

	}

}

func (this *Security) GetAllSecurityDayList() {
	secs := *this.list.Securitys

	for i, v := range secs {
		defer func() {
			if err := recover(); err != nil {
				logging.Error("date:%v------nsid:%v", v.Date, v.Sid)
			}

		}()
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
