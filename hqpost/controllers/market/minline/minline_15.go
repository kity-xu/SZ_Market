package minline

import (
	"ProtocolBuffer/format/kline"

	"haina.com/market/hqpost/models/redistore"

	//"haina.com/share/logging"
)

//生成历史15分钟线
func (this *MinKline) HMinLine_15() {
	rstore15 := redistore.NewHMinKLine(REDISKEY_SECURITY_HMIN15)

	for _, dmin := range *this.list.All { //个股当天数据
		//logging.Debug("SID:%v   今天15分钟时间----%v", dmin.Sid, dmin.Time_15)

		var tmps []*kline.KInfo
		for _, min15 := range *dmin.Time_15 { //当天的每个15分钟
			tmp := &kline.KInfo{}

			var (
				i          int
				min        int32
				AvgPxTotal uint32
			)
			for i, min = range min15 {
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
			tmp.NTime = dmin.Min[min15[len(min15)-1]].NTime //时间
			tmp.NOpenPx = dmin.Min[min15[0]].NOpenPx        //开盘价
			if len(tmps) > 0 {
				tmp.NPreCPx = tmps[len(tmps)-1].NLastPx //昨收价
			} else {
				tmp.NPreCPx = 0
			}
			tmp.NLastPx = dmin.Min[min15[i]].NLastPx //最新价
			tmp.NAvgPx = AvgPxTotal / uint32(i+1)    //平均价
			tmps = append(tmps, tmp)
		}

		//个股当天5分钟数据并入历史
		this.mergeMin(dmin.Sid, rstore15, &tmps)
	}
}
