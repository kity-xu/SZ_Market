package kline

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"

	"haina.com/market/hqpost/models/filestore"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

func HisYearKline(sids *[]int32) {
	for _, sid := range *sids {
		base := new(BaseLine)
		// 获取当天快照
		today, err := GetIntradayKInfo(sid)
		if err != nil {
			logging.Error("严重错误！！！，程序被迫停止执行")
			return
		}

		yearName := filePath(cfg, cfg.File.Year, sid)
		if !lib.IsFileExist(yearName) { // 不存在或其他做第一从生成 TODO
			if err = base.CreateYearLine(sid, yearName, today); err != nil {
				continue
			}
		} else { // 追加周线
			if err = filestore.UpdateYearLineToFile(sid, yearName, today); err != nil {
				logging.Error("UpdateWeekLineToFile: %v", err)
			}
		}
	}
}
func (this *BaseLine) CreateYearLine(sid int32, yearFile string, today *protocol.KInfo) error {
	if err := this.ReadHGSDayLines(sid); err != nil {
		return err
	}
	this.getSecurityYearDay()
	wTable := this.ProduceYearprotocol()
	filestore.MaybeBelongAYear(wTable, today) //第一次生成的时候 如果同属一周加入当天数据
	if err := filestore.WiteHainaFileStore(yearFile, wTable); err != nil {
		logging.Error("WiteHainaFileStore error | %v", err)
	}
	return nil
}

// 月线日期
func (this *BaseLine) getSecurityYearDay() {
	var lastyear int32 = 0
	var dates [][]int32
	var years []int32

	if len(this.date) < 1 {
		logging.Error("SID:%v---No historical data...", this.sid)
		return
	}
	for j, day := range this.date { // v.Date: 单个股票的所有时间
		if lastyear == 0 {
			years = append(years, day)
			lastyear = day / 10000
			continue
		}
		if lastyear == day/10000 {
			years = append(years, day)
			if j == int(len(this.date)-1) { //执行到最后一个
				dates = append(dates, years)
			}
		} else {
			dates = append(dates, years)
			years = nil
			years = append(years, day)
		}
		lastyear = day / 10000

	}
	this.kindDays = &dates
}

// 生成年线
func (this *BaseLine) ProduceYearprotocol() *protocol.KInfoTable {
	//PB
	var klist protocol.KInfoTable

	for _, year := range *(this.kindDays) { //每年
		var (
			i          int
			day        int32
			AvgPxTotal uint32
			tmp        protocol.KInfo //pb类型
		)

		for i, day = range year { //每一天
			stockday := this.sigStock[day]
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
		tmp.NSID = this.sid
		tmp.NTime = this.sigStock[year[0]].NTime     //时间（取每周第一天）
		tmp.NOpenPx = this.sigStock[year[0]].NOpenPx //开盘价（每周第一天的开盘价）
		tmp.NPreCPx = this.sigStock[year[0]].NPreCPx
		tmp.NLastPx = this.sigStock[year[i]].NLastPx //最新价
		tmp.NAvgPx = AvgPxTotal / uint32(i+1)        //平均价

		klist.List = append(klist.List, &tmp)
		//logging.Debug("year线是:%v", klist.List)
	}
	return &klist
}
