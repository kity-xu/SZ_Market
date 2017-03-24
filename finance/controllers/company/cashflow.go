// 现金流量表
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type CashflowInfo struct {
}

func NewCashflowInfo() *CashflowInfo {
	return &CashflowInfo{}
}

func (this *CashflowInfo) getJson(c *gin.Context) (*company.ResponseInfo, error) {
	return company.NewCashflow().GetJson(c)
}

func (this *CashflowInfo) GET(c *gin.Context) {
	data, err := this.getJson(c)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	lib.WriteString(c, 200, data)
}
