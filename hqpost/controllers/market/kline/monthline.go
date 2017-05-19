package kline

import (
	pbk "ProtocolBuffer/format/kline"
	"fmt"

	"haina.com/share/store/redis"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
)

func (this *Security) MonthLine() {
	this.GetMonthDay()
	securitys := *this.list.Securitys

	for _, single := range securitys { //每支股票

		var tmps []pbk.KInfo
		//PB
		var klist pbk.KInfoTable

		for _, month := range *single.MonthDays { //每个月
			var (
				i          int
				day        int32
				AvgPxTotal uint32
				tmp        pbk.KInfo //pb类型
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
		//入PB 入redis
		data, err := proto.Marshal(&klist)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		//入文件
		if e := KlineWriteFile(single.Sid, cfg.File.Month, &data); e != nil {
			return
		}
		//入redis
		key := fmt.Sprintf(REDISKEY_SECURITY_HMONTH, single.Sid)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}
	}

}

func (this *Security) GetMonthDay() {
	securitys := *this.list.Securitys

	for i, v := range securitys { // v: 单个股票
		var yesterday int32 = 0

		var dates [][]int32
		var month []int32
		for j, day := range v.Date { // v.Date: 单个股票的所有时间

			if j == 0 {
				month = append(month, day)
				yesterday = day / 100
				continue
			}
			if yesterday == day/100 {
				month = append(month, day)
			} else {
				dates = append(dates, month)
				month = nil
				month = append(month, day)
			}
			yesterday = day / 100

		}
		//logging.Debug("------day:%v", v.Date)
		//logging.Debug("-month---dates:%+v", dates)
		securitys[i].MonthDays = &dates
	}
}
