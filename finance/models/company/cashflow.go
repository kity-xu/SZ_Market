// 现金流量表
package company

import (
	"haina.com/market/finance/models/finchina"
)

type Cashflow struct {
}

func NewCashflow() *Cashflow {
	return &Cashflow{}
}

func (this *Cashflow) GetJson(req *finchina.RequestParam) (*finchina.ResponseFinAnaJson, error) {
	return finchina.NewCashflowInfo().GetJson(req)
}
