// 资产负债表
package company

import (
	"haina.com/market/finance/models/finchina"
)

type Liabilities struct {
}

func NewLiabilities() *Liabilities {
	return &Liabilities{}
}

func (this *Liabilities) GetJson(req *finchina.RequestParam) (*finchina.ResponseFinAnaJson, error) {
	return finchina.NewLiabilitiesInfo().GetJson(req)
}
