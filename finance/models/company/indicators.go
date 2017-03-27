// 关键指标
package company

import (
	"haina.com/market/finance/models/finchina"
)

type Indicators struct {
}

func NewIndicators() *Indicators {
	return &Indicators{}
}

func (this *Indicators) GetJson(req *finchina.RequestParam) (*finchina.ResponseFinAnaJson, error) {
	return finchina.NewIndicatorsInfo().GetJson(req)
}
