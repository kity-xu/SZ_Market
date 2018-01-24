package fundModel

import (
	"haina.com/share/logging"
	"haina.com/share/models"
)

type FundFlow struct {
	models.Model      `db:"-"`
	LlHugeBuyValue    float64 `db:"llHugeBuyValue"`
	LlBigBuyValue     float64 `db:"llBigBuyValue"`
	LlMiddleBuyValue  float64 `db:"llMiddleBuyValue"`
	LlSmallBuyValue   float64 `db:"llSmallBuyValue"`
	LlHugeSellValue   float64 `db:"llHugeSellValue"`
	LlBigSellValue    float64 `db:"llBigSellValue"`
	LlMiddleSellValue float64 `db:"llMiddleSellValue"`
	LlSmallSellValue  float64 `db:"llSmallSellValue"`
	NTime             string  `db:"nTime"`
}

func NewFundFlow() *FundFlow {
	return &FundFlow{
		Model: models.Model{
			TableName: "sz_fundflow",
			Db:        models.Db,
		},
	}
}

func (this *FundFlow) GetFundData(symbol string, num uint64) ([]*FundFlow, error) {
	//logging.Info(scode)
	var data []*FundFlow
	exps := map[string]interface{}{
		"SYMBOL=?": symbol,
	}
	builder := this.Db.Select("llHugeBuyValue,llBigBuyValue,llMiddleBuyValue,llSmallBuyValue,llHugeSellValue,llBigSellValue,llMiddleSellValue,llSmallSellValue,nTime").
		From(this.TableName).OrderBy("nTime desc").Limit(num)
	_, err := this.SelectWhere(builder, exps).LoadStructs(&data)
	if err != nil {
		logging.Debug("%v", err)
		return data, err
	}
	return data, err
}
func (this *FundFlow) IndustryDateG(swlevelcode string, num uint64) ([]*FundFlow, error) {
	var data []*FundFlow
	builder := this.Db.SelectBySql(`(SELECT SUM(llHugeBuyValue) AS llHugeBuyValue, SUM(llBigBuyValue) AS llBigBuyValue, SUM(llMiddleBuyValue) AS llMiddleBuyValue,
	 SUM(llSmallBuyValue) AS llSmallBuyValue, SUM(llHugeSellValue) AS llHugeSellValue, SUM(llBigSellValue) AS llBigSellValue, SUM(llMiddleSellValue) AS llMiddleSellValue,
	 SUM(llSmallSellValue) AS llSmallSellValue, nTime
     FROM sz_fundflow WHERE swlevel1code = ? GROUP BY nTime ORDER BY nTime DESC LIMIT ?)`, swlevelcode, num)
	_, err := this.SelectWhere(builder, map[string]interface{}{}).LoadStructs(&data)

	if err != nil {
		logging.Info(err.Error())
		return data, err
	}
	return data, err
}
