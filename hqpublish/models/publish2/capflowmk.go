// 市场资金流向
package publish2

import "haina.com/market/hqpublish/models/szdb"

type MkCapflow struct{}

func NewMkCapflow() *MkCapflow {
	return &MkCapflow{}
}

type MarketFundFlowJson struct {
	MARKETID        int32   `json:"marketid"`
	TRADEDATE       int32   `json:"tradedate"`
	HUGEBUYVALUE    float64 `json:"llHugeBuyValue"`
	BIGBUYVALUE     float64 `json:"llBigBuyValue"`
	MIDDLEBUYVALUE  float64 `json:"llMiddleBuyValue"`
	SMALLBUYVALUE   float64 `json:"llSmallBuyValue"`
	HUGESELLVALUE   float64 `json:"llHugeSellValue"`
	BIGSELLVALUE    float64 `json:"llBigSellValue"`
	MIDDLESELLVALUE float64 `json:"llMiddleSellValue"`
	SMALLSELLVALUE  float64 `json:"llSmallSellValue"`
	ENTRYDATE       string  `json:"entrydate"` //#更新日期
	ENTRYTIME       string  `json:"entrytime"` //#更新时间
}

type FundDays struct {
	Num   int32                  `json:"num"`
	Funds *[]*MarketFundFlowJson `json:"funds"`
}

func (*MkCapflow) GetMkCapflow(marketID int32) (*FundDays, error) {
	flows, err := szdb.NewSZ_HQ_MARKETFUNDFLOW().GetMarketFundFlow(60, criterionMarketID(marketID))
	if len(flows) == 0 || err != nil {
		return nil, err
	}
	var fundjson []*MarketFundFlowJson
	for _, v := range flows {
		fj := &MarketFundFlowJson{
			MARKETID:        marketID,
			TRADEDATE:       v.TRADEDATE,
			HUGEBUYVALUE:    v.HUGEBUYVALUE.Float64,
			BIGBUYVALUE:     v.BIGBUYVALUE.Float64,
			MIDDLEBUYVALUE:  v.MIDDLEBUYVALUE.Float64,
			SMALLBUYVALUE:   v.SMALLBUYVALUE.Float64,
			HUGESELLVALUE:   v.HUGESELLVALUE.Float64,
			BIGSELLVALUE:    v.BIGSELLVALUE.Float64,
			MIDDLESELLVALUE: v.MIDDLESELLVALUE.Float64,
			SMALLSELLVALUE:  v.SMALLSELLVALUE.Float64,
			ENTRYDATE:       v.ENTRYDATE.String,
			ENTRYTIME:       v.ENTRYTIME.String,
		}
		fundjson = append(fundjson, fj)
	}
	res := &FundDays{
		Num:   int32(len(flows)),
		Funds: &fundjson,
	}

	return res, nil
}

// 规范marketID
func criterionMarketID(marketID int32) int32 {
	switch marketID {
	case 100:
		marketID = 100000000
	case 200:
		marketID = 200000000
	case 300:
		marketID = 300000000
	case 400:
		marketID = 400000000
	}
	return marketID
}
