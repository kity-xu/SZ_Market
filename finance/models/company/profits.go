// 利润表
package company

import (
	"haina.com/market/finance/models/finchina"
)

type Profits struct {
}

func NewProfits() *Profits {
	return &Profits{}
}

func (this *Profits) GetJson(req *finchina.RequestParam) (*finchina.ResponseFinAnaJson, error) {
	return finchina.NewProfitsInfo().GetJson(req)
}
