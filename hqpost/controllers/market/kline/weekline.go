package kline

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"

	"haina.com/market/hqpost/models"

	"haina.com/market/hqpost/models/filestore"
	"haina.com/market/hqpost/models/redistore"
	"haina.com/share/logging"
)

func (this *Security) WeekLine() {
	this.GetAllSecurityDayList()
	securitys := *this.list.Securitys

	rstore := redistore.NewHKLine(REDISKEY_SECURITY_HWEEK)

	for _, single := range securitys { // 以sid分类的单个股票
		var (
			klist *protocol.KInfoTable
			err   error
		)
		filepath, ok := filestore.CheckFileSoteDir(single.Sid, cfg.File.Path, cfg.File.Week)
		if !ok { //不存在，做第一次生成
			klist = produceWeeprotocol(&single)
			if single.today != nil {
				filestore.MaybeBelongAWeek(klist, single.today) //第一次生成的时候 如果同属一周加入当天数据
			}

			//1.入文件
			if err = filestore.WiteHainaFileStore(filepath, klist); err != nil {
				logging.Error("%v", err.Error())
			}

			//redis做第一次生成
			for _, v := range klist.List {
				if err := rstore.LPushHisKLine(single.Sid, v); err != nil {
					logging.Error("%v", err.Error())
					return
				}
			}
		} else {
			if single.today != nil { //今天的数据加入历史
				if err = filestore.UpdateWeekLineToFile(filepath, single.today); err != nil {
					logging.Error("%v", err.Error())
				}

				if err = rstore.UpdateWeekKLineToRedis(single.Sid, single.today); err != nil {
					if err != models.ERROR_REDIS_LIST_NULL {
						return
					} else {
						continue
					}
				}
			}
		}
	}

}

func produceWeeprotocol(single *SingleSecurity) *protocol.KInfoTable {
	var tmps []protocol.KInfo
	var klist protocol.KInfoTable

	for _, week := range *single.WeekDays { //每一周
		if single.Sid == 100000001 {
			logging.Debug("Week:%v", week)
		}
		tmp := protocol.KInfo{}

		var (
			i          int
			day        int32
			AvgPxTotal uint32
		)
		for i, day = range week { //每一天
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

		klist.List = append(klist.List, &tmp)

	}
	return &klist
}

func (this *Security) GetAllSecurityDayList() {
	secs := *this.list.Securitys

	for i, v := range secs {
		if len(v.Date) < 1 {
			logging.Error("SID:%v---No historical data...", v.Sid)
			continue
		}
		var wday [][]int32
		sat, _ := filestore.DateAdd(v.Date[0]) //该股票第一个交易日所在周的周日（周六可能会有交易）

		var dates []int32
		for j, date := range v.Date {
			if filestore.IntToTime(int(date)).Before(sat) {
				dates = append(dates, date)
				if j == int(len(v.Date)-1) { //执行到最后一个
					wday = append(wday, dates)
				}
			} else {
				wday = append(wday, dates)

				sat, _ = filestore.DateAdd(date)
				dates = nil
				dates = append(dates, date)
			}
		}
		secs[i].WeekDays = &wday
	}
}
