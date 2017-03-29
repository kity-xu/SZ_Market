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

func (this *CashflowInfo) getJson(req *RequestParam) (*company.RespFinAnaJson, error) {
	return company.NewCashflow().GetJson(req.SCode, req.Type, req.PerPage, req.Page)
}

func (this *CashflowInfo) GET(c *gin.Context) {
	scode := c.Query("scode")
	stype := c.Query("type")
	spage := c.Query("page")
	perPage := c.Query("perpage")

	req := CheckAndNewRequestParam(scode, stype, perPage, spage)
	if req == nil {
		lib.WriteString(c, 40004, nil)
		return
	}

	data, err := this.getJson(req)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	lib.WriteString(c, 200, data)
}
