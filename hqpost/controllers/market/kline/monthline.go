package kline

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"

	"haina.com/market/hqpost/models"
	"haina.com/share/logging"

	"haina.com/market/hqpost/models/filestore"
	"haina.com/market/hqpost/models/redistore"
)

func (this *Security) MonthLine() {
	this.GetMonthDay()
	securitys := *this.list.Securitys

	rstore := redistore.NewHKLine(REDISKEY_SECURITY_HMONTH)

	for _, single := range securitys { //每支股票
		var (
			klist *protocol.KInfoTable
			err   error
		)
		filepath, ok := filestore.CheckFileSoteDir(single.Sid, cfg.File.Path, cfg.File.Month)
		if !ok { //不存在，做第一次生成
			klist = produceMonthline(&single) //klist已包含当天的数据了

			//1.入文件
			filestore.WiteHainaFileStore(filepath, klist)

			//2.redis做第一次生成
			for _, v := range klist.List {
				if err := rstore.LPushHisKLine(single.Sid, v); err != nil {
					logging.Error("%v", err.Error())
					return
				}
			}
		} else {
			if single.today != nil {
				if err = filestore.UpdateMonthLineToFile(filepath, single.today); err != nil {
					logging.Error("%v", err.Error())
				}

				if err = rstore.UpdateMonthKLineToRedis(single.Sid, single.today); err != nil {
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

func produceMonthline(single *SingleSecurity) *protocol.KInfoTable {
	var tmps []protocol.KInfo
	//PB
	var klist protocol.KInfoTable

	for _, month := range *single.MonthDays { //每个月
		var (
			i          int
			day        int32
			AvgPxTotal uint32
			tmp        protocol.KInfo //pb类型
		)

		for i, day = range month { //每一天
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
		tmp.NTime = single.SigStock[month[0]].NTime     //时间（取每周第一天）
		tmp.NOpenPx = single.SigStock[month[0]].NOpenPx //开盘价（每周第一天的开盘价）
		if len(tmps) > 0 {
			tmp.NPreCPx = tmps[len(tmps)-1].NLastPx //昨收价(上周的最新价)
		} else {
			tmp.NPreCPx = 0
		}
		tmp.NLastPx = single.SigStock[month[i]].NLastPx //最新价
		tmp.NAvgPx = AvgPxTotal / uint32(i+1)           //平均价

		tmps = append(tmps, tmp)
		//logging.Debug("yue线是:%v", tmps)
		//入PB
		klist.List = append(klist.List, &tmp)
	}
	return &klist
}

func (this *Security) GetMonthDay() {
	securitys := *this.list.Securitys

	for i, v := range securitys { // v: 单个股票
		var yesterday int32 = 0

		var dates [][]int32
		var month []int32

		if len(v.Date) < 1 {
			logging.Error("SID:%v---No historical data...", v.Sid)
			continue
		}
		for j, day := range v.Date { // v.Date: 单个股票的所有时间

			if j == 0 {
				month = append(month, day)
				yesterday = day / 100
				continue
			}
			if yesterday == day/100 {
				month = append(month, day)
				if j == int(len(v.Date)-1) { //执行到最后一个
					dates = append(dates, month)
				}
			} else {
				dates = append(dates, month)
				month = nil
				month = append(month, day)
			}
			yesterday = day / 100

		}
		securitys[i].MonthDays = &dates
	}
}
