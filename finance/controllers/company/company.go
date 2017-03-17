package company

import (
	"strings"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
)

type CompanyInfo struct {
}

func NewCompanyInfo() *CompanyInfo {
	return &CompanyInfo{}
}

func (this *CompanyInfo) GetInfo(c *gin.Context) {
	para := c.Query("scode")

	if !strings.EqualFold(para, "600036.SH") {
		lib.WriteString(c, 300, "invalid scode..")
		return
	}
	var cominfo *company.Company
	cominfo = company.NewCompany()
	lib.WriteString(c, 200, cominfo)
}
