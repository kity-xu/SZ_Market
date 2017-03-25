// 资产负债表
package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"haina.com/market/finance/models/finchina"
)

type LiabilitiesInfo struct {
}

func NewLiabilitiesInfo() *LiabilitiesInfo {
	return &LiabilitiesInfo{}
}

func (this *LiabilitiesInfo) getJson(req *company.RequestParam) (*company.ResponseFinAnaJson, error) {
	sess := company.Session{}
	sess.Responser = finchina.NewLiabilitiesInfo()
	return sess.GetJson(req)
}

func (this *LiabilitiesInfo) GET(c *gin.Context) {
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
