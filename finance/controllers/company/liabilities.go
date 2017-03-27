// 资产负债表
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/market/finance/models/finchina"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type RequestLiabilitiesInfo struct {
}

func NewRequestLiabilitiesInfo() *RequestLiabilitiesInfo {
	return &RequestLiabilitiesInfo{}
}

func (this *RequestLiabilitiesInfo) getJson(req *finchina.RequestParam) (*finchina.ResponseFinAnaJson, error) {
	return company.NewLiabilities().GetJson(req)
}

func (this *RequestLiabilitiesInfo) GET(c *gin.Context) {
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
