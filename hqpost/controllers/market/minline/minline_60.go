package minline

import (
	"ProtocolBuffer/format/kline"

	//"haina.com/share/logging"
)

//生成历史60分钟线
func (this *MinKline) HMinLine_60() {
	for _, dmin := range *this.list.All { //个股当天数据
		//logging.Debug("SID:%v   今天60分钟时间----%v", dmin.Sid, dmin.Time_60)

		var tmps kline.HMinLineDay
		for _, min60 := range *dmin.Time_60 { //当天的每个60分钟
			tmp := &kline.KInfo{}

			var (
				i          int
				min        int32
				AvgPxTotal uint32
			)
			for i, min = range min60 {
				stockmin := dmin.Min[min]
				if tmp.NHighPx < stockmin.NHighPx || tmp.NHighPx == 0 { //最高价
					tmp.NHighPx = stockmin.NHighPx
				}
				if tmp.NLowPx > stockmin.NLowPx || tmp.NLowPx == 0 { //最低价
					tmp.NLowPx = stockmin.NLowPx
				}
				tmp.LlVolume += stockmin.LlVolume //成交量
				tmp.LlValue += stockmin.LlValue   //成交额
				AvgPxTotal += stockmin.NAvgPx
			}

			tmp.NSID = dmin.Sid
			tmp.NTime = dmin.Min[min60[len(min60)-1]].NTime //时间
			tmp.NOpenPx = dmin.Min[min60[0]].NOpenPx        //开盘价
			if len(tmps.List) > 0 {
				tmp.NPreCPx = tmps.List[len(tmps.List)-1].NLastPx //昨收价
			} else {
				tmp.NPreCPx = 0
			}
			tmp.NLastPx = dmin.Min[min60[i]].NLastPx //最新价
			tmp.NAvgPx = AvgPxTotal / uint32(i+1)    //平均价
			tmps.List = append(tmps.List, tmp)
		}

		tmps.Date = GetDateToday()
		//个股当天5分钟数据并入历史
		this.mergeMin(dmin.Sid, REDISKEY_SECURITY_HMIN60, tmps)
	}
}
