package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
)

type DividendInfo struct {
}

func NewDividendInfo() *DividendInfo {
	return &DividendInfo{}
}

func (this *DividendInfo) GetDiv(c *gin.Context) {
	div := new(company.Dividend).GetDiv()
	lib.WriteString(c, 200, div)
}
func (this *DividendInfo) GetSEO(c *gin.Context) {
	seo := new(company.Dividend).GetSEO()
	lib.WriteString(c, 200, seo)
}
func (this *DividendInfo) GetRO(c *gin.Context) {
	ro := new(company.Dividend).GetRO()
	lib.WriteString(c, 200, ro)
}
