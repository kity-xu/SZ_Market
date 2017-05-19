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
		var tmps []pbk.KInfo

		//PB
		var klist pbk.KInfoTable

		for _, week := range *single.WeekDays { //每一周
			//logging.Debug("Week:%v", week)
			tmp := pbk.KInfo{}

			var (
				i          int
				day        int32
				AvgPxTotal uint32
			)
			for i, day = range week { //每一天
				//logging.Debug("day:%v---single.SigStock[day]:%v", day, single.SigStock[day])
				stockday := single.SigStock[day]
				if tmp.NHighPx < stockday.NHighPx || tmp.NHighPx == 0 { //最高价
					tmp.NHighPx = stockday.NHighPx
				}
				if tmp.NLowPx > stockday.NLowPx || tmp.NLowPx == 0 { //最低价
					tmp.NLowPx = stockday.NLowPx
				}
				tmp.LlVolume += stockday.LlVolume //成交量
				tmp.LlValue += stockday.LlValue   //成交额
				AvgPxTotal += stockday.NAvgPx
			}

			tmp.NSID = single.Sid
			tmp.NTime = single.SigStock[week[0]].NTime     //时间（取每周第一天）
			tmp.NOpenPx = single.SigStock[week[0]].NOpenPx //开盘价（每周第一天的开盘价）
			if len(tmps) > 0 {
				tmp.NPreCPx = tmps[len(tmps)-1].NLastPx //昨收价(上周的最新价)
			} else {
				tmp.NPreCPx = 0
			}
			tmp.NLastPx = single.SigStock[week[i]].NLastPx //最新价
			tmp.NAvgPx = AvgPxTotal / uint32(i+1)          //平均价
			tmps = append(tmps, tmp)
			//logging.Debug("周线是:%v", tmps)

			//入PB
			klist.List = append(klist.List, &tmp)

		}

		//PB
		data, err := proto.Marshal(&klist)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		//入文件
		if e := KlineWriteFile(single.Sid, cfg.File.Week, &data); e != nil {
			return
		}

		//入redis
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
		sat, _ := DateAdd(v.Sid, int(v.Date[0])) //该股票第一个交易日所在周的周日（周六可能会有交易）

		var dates []int32
		for _, date := range v.Date {
			if IntToTime(int(date)).Before(sat) {
				dates = append(dates, date)
			} else {
				//logging.Debug("------一周的日期是：%v------", dates) //it's here
				wday = append(wday, dates)

				sat, _ = DateAdd(v.Sid, int(date))
				dates = nil
				dates = append(dates, date)
			}
		}

		wday = append(wday, dates)
		secs[i].WeekDays = &wday

	}
	//logging.Debug("-----单个股票，所有周天：%v------", (*this.week.Securitys)[0].WeekDays)
	//logging.Debug("-----单个股票，secs[0].date：%v------", (*this.week.Securitys)[0].Date)
}
