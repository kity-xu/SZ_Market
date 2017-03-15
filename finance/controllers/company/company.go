package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
)

type CompanyInfo struct {
}

func NewCompanyInfo() *CompanyInfo {
	return &CompanyInfo{}
}

func (this *CompanyInfo) GetInfo(c *gin.Context) {
	lib.WriteString(c, 200, "test")
}
