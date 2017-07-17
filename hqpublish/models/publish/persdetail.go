// 个股详情
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/memdata"
	"haina.com/market/hqpublish/models/publish/security"
	"haina.com/share/logging"
)

type PerSDetailM struct {
}

func NewPerSDetailM() *PerSDetailM {
	return &PerSDetailM{}
}

// 个股详情
func (this *PerSDetailM) GetPerSDtail(req *protocol.RequestPerSDetail) (*protocol.PayloadPerSDetail, error) {

	info := memdata.HandleSwapInfo(int32(req.SID))
	if info == nil {
		logging.Error("select stock type error")
	}

	var statsd *protocol.PayloadSecurityStatic
	switch t := info.Type[1]; t {
	case 'S': // 股票
		// 调用静态数据处理
		var reqsta protocol.RequestSecurityStatic
		reqsta.SID = int32(req.SID)
		stat, err := security.NewSecurityStatic().GetSecurityStatic(&reqsta)
		if err != nil {
			logging.Info("调用静态数据处理 error %v", err)
		}
		statsd = stat
	case 'I': // 指数
		statsd.SSInfo.LlTotalShare = 0
		statsd.SSInfo.LlCircuShare = 0
		statsd.SSInfo.LlFlowProperty = 0
		statsd.SSInfo.SzStatus = ""
	}

	// 调用快照数据处理
	var psdb protocol.PerSDetailBase
	var reqsh protocol.RequestSnapshot
	reqsh.SID = int32(req.SID)
	stock, index, err := NewSnapshot().GetSnapshot(&reqsh)

	if err != nil {
		logging.Error("%v", err)

	}

	if stock != nil {
		psdb.NOpenPx = stock.SnapInfo.NOpenPx
		psdb.NHighPx = stock.SnapInfo.NHighPx
		psdb.NPreClosePx = stock.SnapInfo.NPreClosePx
		psdb.NLowPx = stock.SnapInfo.NLowPx
		psdb.NLowLimitPx = stock.SnapInfo.NLowLimitPx
		psdb.NPxAmplitude = stock.SnapInfo.NPxAmplitude
		psdb.NHighLimitPx = stock.SnapInfo.NHighLimitPx
		psdb.LlTotalMo = stock.SnapInfo.LlValue
		psdb.LlInnerVolume = stock.SnapInfo.LlInnerVolume
		psdb.LlOuterVolume = stock.SnapInfo.LlOuterVolume
		psdb.NTurnOver = stock.SnapInfo.NTurnOver
		psdb.NLiangbi = stock.SnapInfo.NLiangbi
		psdb.NWeibi = stock.SnapInfo.NWeibi
		psdb.LlToBidVol = stock.SnapInfo.LlToBidVol - stock.SnapInfo.LlToOfferVol
		psdb.NPB = stock.SnapInfo.NPB
		psdb.NPE = stock.SnapInfo.NPE
		psdb.NLastPx = stock.SnapInfo.NLastPx
		psdb.NTime = stock.SnapInfo.NTime
		psdb.LlValue = stock.SnapInfo.LlValue
		psdb.LlVolume = stock.SnapInfo.LlVolume
		psdb.LlTradeNum = stock.SnapInfo.LlTradeNum

	} else if index != nil {
		psdb.NOpenPx = index.SnapInfo.NOpenPx
		psdb.NHighPx = index.SnapInfo.NHighPx
		psdb.NPreClosePx = index.SnapInfo.NPreClosePx
		psdb.NLowPx = index.SnapInfo.NLowPx
		psdb.NLowLimitPx = index.SnapInfo.NLowLimitPx
		psdb.NPxAmplitude = index.SnapInfo.NPxAmplitude
		psdb.NHighLimitPx = index.SnapInfo.NHighLimitPx
		psdb.LlTotalMo = index.SnapInfo.LlValue

		psdb.NTurnOver = index.SnapInfo.NTurnOver
		psdb.NLiangbi = index.SnapInfo.NLiangbi
		psdb.NWeibi = index.SnapInfo.NWeibi
		psdb.LlToBidVol = 0
		psdb.NPB = index.SnapInfo.NPB
		psdb.NPE = index.SnapInfo.NPE
		psdb.NLastPx = index.SnapInfo.NLastPx
		psdb.NTime = index.SnapInfo.NTime
		psdb.LlValue = index.SnapInfo.LlValue
		psdb.LlVolume = index.SnapInfo.LlVolume
		psdb.LlTradeNum = index.SnapInfo.LlTradeNum
	}

	// 处理成交量和成交额可能为零
	var avgp int32
	if psdb.LlValue == 0 || psdb.LlVolume == 0 {
		avgp = 0
	} else {
		avgp = int32(psdb.LlValue / psdb.LlVolume)
	}

	pers := &protocol.PayloadPerSDetail{
		SID: req.SID,
		Psdb: &protocol.PerSDetailBase{
			NOpenPx:        psdb.NOpenPx,                         // 开盘价	(*10000)
			NHighPx:        psdb.NHighPx,                         // 最高价	(*10000)
			NAvePrice:      avgp,                                 // 均价	 	(*10000)
			NPreClosePx:    psdb.NPreClosePx,                     // 昨收价	(*10000)
			NLowPx:         psdb.NLowPx,                          // 最低价	(*10000)
			NLowLimitPx:    psdb.NLowLimitPx,                     // 跌停价格	(*10000)
			NPxAmplitude:   psdb.NPxAmplitude,                    // 振幅  	(*10000)
			NHighLimitPx:   psdb.NHighLimitPx,                    // 涨停价格	(*10000)
			LlTotalMo:      psdb.LlValue,                         // 总额	    (*10000)
			LlInnerVolume:  psdb.LlInnerVolume,                   // 内盘成交量(*10000)
			LlOuterVolume:  psdb.LlOuterVolume,                   // 外盘成交量(*10000)
			NTurnOver:      psdb.NTurnOver,                       // 换手率	(*10000)
			LlTotalShare:   statsd.SSInfo.LlTotalShare,           // 总股本(万股)(*10000)
			NLiangbi:       psdb.NLiangbi,                        // 量比  !   (*100)
			NWeibi:         psdb.NWeibi,                          // 委比		(*10000)
			LlCircuShare:   statsd.SSInfo.LlCircuShare,           // 流通股(万股)(*10000)
			LlToBidVol:     psdb.LlToBidVol,                      // 委差(总委买量llToBidVol-llToOfferVol总委卖量)
			NPB:            psdb.NPB,                             // 动态市净率(*10000)
			NPE:            psdb.NPE,                             // 动态市盈率(*10000)
			LlFlowProperty: statsd.SSInfo.LlFlowProperty,         // 流动资产
			Bid:            make([]*protocol.QuoteRecordp, 0, 5), // 买五档
			Offer:          make([]*protocol.QuoteRecordp, 0, 5), // 卖五档
			NLastPx:        psdb.NLastPx,                         // 最新价
			NTime:          psdb.NTime,                           // 时间 unix time
			SzStatus:       statsd.SSInfo.SzStatus,               // 证券状态
			LlVolume:       psdb.LlVolume,                        // 成交量
			LlTradeNum:     psdb.LlTradeNum,                      // 成交笔数
			LlValue:        psdb.LlValue,                         // 成交额(*10000)

		},
	}

	for _, v := range psdb.Bid {
		bid := &protocol.QuoteRecordp{
			NPx:      v.NPx,
			LlVolume: v.LlVolume,
		}
		pers.Psdb.Bid = append(pers.Psdb.Bid, bid)
	}

	for _, v := range psdb.Offer {
		offer := &protocol.QuoteRecordp{
			NPx:      v.NPx,
			LlVolume: v.LlVolume,
		}
		pers.Psdb.Offer = append(pers.Psdb.Offer, offer)
	}
	return pers, err
}
