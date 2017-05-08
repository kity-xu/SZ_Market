package kline

import (
	pbk "ProtocolBuffer/format/kline"
	"fmt"

	"haina.com/share/store/redis"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
)

func (this *Security) YearLine() {
	this.GetYearDay()
	securitys := *this.list.Securitys

	for _, single := range securitys { //每支股票

		var tmps []StockSingle
		//PB
		var klist pbk.KInfoTable

		//logging.Debug("YearDays:%+v", single.YearDays)

		for _, year := range *single.YearDays { //每年

			tmp := StockSingle{}
			var mdata pbk.KInfo //pb类型

			var (
				i          int
				day        int32
				AvgPxTotal uint32
			)

			for i, day = range year { //每一天
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
			tmp.Time = single.SigStock[year[0]].Time     //时间（取每周第一天）
			tmp.OpenPx = single.SigStock[year[0]].OpenPx //开盘价（每周第一天的开盘价）
			if len(tmps) > 0 {
				tmp.PreCPx = tmps[len(tmps)-1].LastPx //昨收价(上周的最新价)
			} else {
				tmp.PreCPx = 0
			}
			tmp.LastPx = single.SigStock[year[i]].LastPx //最新价
			tmp.AvgPx = AvgPxTotal / uint32(i+1)         //平均价

			tmps = append(tmps, tmp)
			//logging.Debug("year线是:%v", tmps)
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
		//logging.Debug("year line:%+v", klist)
		//入PB 入redis
		data, err := proto.Marshal(&klist)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		key := fmt.Sprintf(REDISKEY_SECURITY_HYEAR, single.Sid)
		if err := redis.Set(key, data); err != nil {
			logging.Fatal("%v", err)
		}
	}

}

func (this *Security) GetYearDay() {
	securitys := *this.list.Securitys

	for i, v := range securitys { // v: 单个股票
		var lsatyear int32 = 0

		var dates [][]int32
		var years []int32
		for _, day := range v.Date { // v.Date: 单个股票的所有时间
			if lsatyear == 0 {
				years = append(years, day)
				lsatyear = day / 10000
				continue
			}
			if lsatyear == day/10000 {
				years = append(years, day)
			} else {
				dates = append(dates, years)
				years = nil
				years = append(years, day)
			}
			lsatyear = day / 10000

		}
		//logging.Debug("------day:%v", v.Date)
		//logging.Debug("-year---dates:%+v", dates)
		securitys[i].YearDays = &dates
	}
}
