// 个股详情
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

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

	// 调用静态数据处理
	var reqsta protocol.RequestSecurityStatic
	reqsta.SID = int32(req.SID)
	stat, err := security.NewSecurityStatic().GetSecurityStatic(&reqsta)
	if err != nil {
		logging.Info("调用静态数据处理 error %v", err)
	}

	// 调用快照数据处理
	var reqsnap protocol.RequestSnapshot
	reqsnap.SID = int32(req.SID)
	snap, err := NewStockSnapshot().GetStockSnapshot(&reqsnap)
	if err != nil {
		logging.Info("调用快照数据处理 error %v", err)
	}

	// 处理成交量和成交额可能为零
	var avgp int32
	if snap.SnapInfo.LlValue == 0 || snap.SnapInfo.LlVolume == 0 {
		avgp = 0
	} else {
		avgp = int32(snap.SnapInfo.LlValue / snap.SnapInfo.LlVolume)
	}

	pers := &protocol.PayloadPerSDetail{
		SID: req.SID,
		Psdb: &protocol.PerSDetailBase{
			NOpenPx:        snap.SnapInfo.NOpenPx,                                 // 开盘价	(*10000)
			NHighPx:        snap.SnapInfo.NHighPx,                                 // 最高价	(*10000)
			NAvePrice:      avgp,                                                  // 均价	 	(*10000)
			NPreClosePx:    snap.SnapInfo.NPreClosePx,                             // 昨收价	(*10000)
			NLowPx:         snap.SnapInfo.NLowPx,                                  // 最低价	(*10000)
			NLowLimitPx:    snap.SnapInfo.NLowLimitPx,                             // 跌停价格	(*10000)
			NPxAmplitude:   snap.SnapInfo.NPxAmplitude,                            // 振幅  	(*10000)
			NHighLimitPx:   snap.SnapInfo.NHighLimitPx,                            // 涨停价格	(*10000)
			LlTotalMo:      snap.SnapInfo.LlValue,                                 // 总额	    (*10000)
			LlInnerVolume:  snap.SnapInfo.LlInnerVolume,                           // 内盘成交量(*10000)
			LlOuterVolume:  snap.SnapInfo.LlOuterVolume,                           // 外盘成交量(*10000)
			NTurnOver:      snap.SnapInfo.NTurnOver,                               // 换手率	(*10000)
			LlTotalShare:   stat.SSInfo.LlTotalShare,                              // 总股本(万股)(*10000)
			NLiangbi:       snap.SnapInfo.NLiangbi,                                // 量比  !   (*100)
			NWeibi:         snap.SnapInfo.NWeibi,                                  // 委比		(*10000)
			LlCircuShare:   stat.SSInfo.LlCircuShare,                              // 流通股(万股)(*10000)
			LlToBidVol:     snap.SnapInfo.LlToBidVol - snap.SnapInfo.LlToOfferVol, // 委差(总委买量llToBidVol-llToOfferVol总委卖量)
			NPB:            snap.SnapInfo.NPB,                                     // 动态市净率(*10000)
			NPE:            snap.SnapInfo.NPE,                                     // 动态市盈率(*10000)
			LlFlowProperty: stat.SSInfo.LlFlowProperty,                            // 流动资产
			Bid:            make([]*protocol.QuoteRecordp, 0, 5),
			Offer:          make([]*protocol.QuoteRecordp, 0, 5),
		},
	}

	for _, v := range snap.SnapInfo.Bid {
		bid := &protocol.QuoteRecordp{
			NPx:      v.NPx,
			LlVolume: v.LlVolume,
		}
		pers.Psdb.Bid = append(pers.Psdb.Bid, bid)
	}
	for _, v := range snap.SnapInfo.Offer {
		offer := &protocol.QuoteRecordp{
			NPx:      v.NPx,
			LlVolume: v.LlVolume,
		}
		pers.Psdb.Offer = append(pers.Psdb.Offer, offer)
	}
	return pers, err
}
