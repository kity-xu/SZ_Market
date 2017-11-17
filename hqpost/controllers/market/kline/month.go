package kline

import (
	"ProtocolBuffer/projects/hqpost/go/protocol"

	"haina.com/market/hqpost/models/filestore"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

func HisMonthKline(sids *[]int32) {

	for _, sid := range *sids {
		base := new(BaseLine)
		// 获取当天快照
		today, err := GetIntradayKInfo(sid)
		if err != nil {
			logging.Error("严重错误！！！，程序被迫停止执行")
			return
		}

		monthName := filePath(cfg, cfg.File.Month, sid)
		if !lib.IsFileExist(monthName) { // 不存在或其他做第一从生成 TODO
			if err = base.CreateMonthLine(sid, monthName); err != nil { // 在此之前day已更新
				continue
			}
		} else { // 追加周线
			if err = filestore.UpdateMonthLineToFile(sid, monthName, today); err != nil {
				logging.Error("UpdateWeekLineToFile: %v", err)
			}
		}
	}
}

// 读day文件生成week
func (this *BaseLine) CreateMonthLine(sid int32, monthFile string) error {
	err := this.ReadHGSDayLines(sid)
	if err != nil {
		return err
	}
	this.getSecurityMonthDay()
	wTable := this.ProduceMonthprotocol()
	if err := filestore.WiteHainaFileStore(monthFile, wTable); err != nil {
		logging.Error("CreateMonthLine: WiteHainaFileStore error | %v", err)
	}
	return nil
}

// 月线日期
func (this *BaseLine) getSecurityMonthDay() {
	var lastyear int32 = 0

	var dates [][]int32
	var month []int32

	if len(this.date) < 1 {
		logging.Error("SID:%v---No historical data...", this.sid)
		return
	}
	for j, day := range this.date { // v.Date: 单个股票的所有时间
		if j == 0 {
			month = append(month, day)
			lastyear = day / 100
			continue
		}
		if lastyear == day/100 {
			month = append(month, day)
			if j == int(len(this.date)-1) { //执行到最后一个
				dates = append(dates, month)
			}
		} else {
			dates = append(dates, month)
			month = nil
			month = append(month, day)
		}
		lastyear = day / 100

	}
	this.kindDays = &dates
}

// 生成月线
func (this *BaseLine) ProduceMonthprotocol() *protocol.KInfoTable {
	var tmps []protocol.KInfo
	//PB
	var klist protocol.KInfoTable

	for _, month := range *(this.kindDays) { //每个月
		var (
			i          int
			day        int32
			AvgPxTotal uint32
			tmp        protocol.KInfo //pb类型
		)

		for i, day = range month { //每一天
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
		tmp.NTime = this.sigStock[month[0]].NTime     //时间
		tmp.NOpenPx = this.sigStock[month[0]].NOpenPx //开盘价
		tmp.NPreCPx = this.sigStock[month[0]].NPreCPx
		tmp.NLastPx = this.sigStock[month[i]].NLastPx //最新价
		tmp.NAvgPx = AvgPxTotal / uint32(i+1)         //平均价

		tmps = append(tmps, tmp)
		//logging.Debug("yue线是:%v", tmps)
		//入PB
		klist.List = append(klist.List, &tmp)
	}
	return &klist
}
