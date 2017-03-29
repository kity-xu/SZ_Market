// 关键指标
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type IndicatorsInfo struct {
}

func NewIndicatorsInfo() *IndicatorsInfo {
	return &IndicatorsInfo{}
}

func (this *IndicatorsInfo) getJson(req *RequestParam) (*company.RespFinAnaJson, error) {
	return company.NewIndicators().GetJson(req.SCode, req.Type, req.PerPage, req.Page)
}

func (this *IndicatorsInfo) GET(c *gin.Context) {
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
