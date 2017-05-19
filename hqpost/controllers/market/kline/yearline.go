package kline

import (
	pbk "ProtocolBuffer/format/kline"
	"fmt"

	"github.com/golang/protobuf/proto"
	"haina.com/share/logging"
	"haina.com/share/store/redis"
)

func (this *Security) YearLine() {
	this.GetYearDay()
	securitys := *this.list.Securitys

	for _, single := range securitys { //每支股票

		//PB
		var klist pbk.KInfoTable

		//logging.Debug("YearDays:%+v", single.YearDays)

		for _, year := range *single.YearDays { //每年
			var (
				i          int
				day        int32
				AvgPxTotal uint32
				tmp        pbk.KInfo //pb类型
			)

			for i, day = range year { //每一天
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
			tmp.NTime = single.SigStock[year[0]].NTime     //时间（取每周第一天）
			tmp.NOpenPx = single.SigStock[year[0]].NOpenPx //开盘价（每周第一天的开盘价）
			if len(klist.List) > 0 {
				tmp.NPreCPx = klist.List[len(klist.List)-1].NLastPx //昨收价(上周的最新价)
			} else {
				tmp.NPreCPx = 0
			}
			tmp.NLastPx = single.SigStock[year[i]].NLastPx //最新价
			tmp.NAvgPx = AvgPxTotal / uint32(i+1)          //平均价

			klist.List = append(klist.List, &tmp)
			//logging.Debug("year线是:%v", klist.List)
		}
		//logging.Debug("year line:%+v", klist)
		//入PB 入redis
		data, err := proto.Marshal(&klist)
		if err != nil {
			logging.Error("Encode protocbuf of week Line error...%v", err.Error())
			return
		}

		//入文件
		if e := KlineWriteFile(single.Sid, cfg.File.Year, &data); e != nil {
			return
		}

		//入redis
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
