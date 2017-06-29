// 分价成交
package publish

import (
	. "haina.com/share/models"
)

type TradePriceRecordM struct {
	Model `db:"-"`
}

func NewTradePriceRecordM() *TradePriceRecordM {
	return &TradePriceRecordM{
		Model: Model{
			CacheKey: REDISKEY_TRADE_PRICE,
		},
	}
}

//// 获取板块列表
//func (this *TradePriceRecordM) GetStockBlockBase(req *protocol.RequestStockBlockBase) (*protocol.PayloadStockBlockSet, error) {

//	var psb protocol.PayloadStockBlockBase
//}