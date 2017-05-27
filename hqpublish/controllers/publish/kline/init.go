package kline

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models/publish"
	"haina.com/share/logging"
)

//沪深交易市场当前交易日
var (
	Trade_100 int32 = 0
	Trade_200 int32 = 0
)

func init() {
	mlist, err := getMarketStatus()
	if err != nil {
		logging.Error("%v", err)
		return
	}
	for _, v := range mlist.MSList {
		if v.NMarket == 100000000 {
			Trade_100 = v.NTradeDate
		}
		if v.NMarket == 200000000 {
			Trade_200 = v.NTradeDate
		}
	}
}

//市场状态获取当前交易日
func getMarketStatus() (*protocol.PayloadMarketStatus, error) {
	var req = protocol.RequestMarketStatus{
		Num:          2,
		MarketIDList: []int32{100000000, 200000000},
	}

	res, err := publish.NewMarketStatus().GetPayloadObj(&req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
