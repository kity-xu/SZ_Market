// 利润表
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type ProfitsInfo struct {
}

func NewProfitsInfo() *ProfitsInfo {
	return &ProfitsInfo{}
}

func (this *ProfitsInfo) getJson(req *RequestParam) (*company.RespFinAnaJson, error) {
	return company.NewProfits().GetJson(req.SCode, req.Type, req.PerPage, req.Page)
}

func (this *ProfitsInfo) GET(c *gin.Context) {
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
