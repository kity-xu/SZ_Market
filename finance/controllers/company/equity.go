package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
)

type EquityInfo struct {
}

func NewEquityInfo() *EquityInfo {
	return &EquityInfo{}
}

/**
获取股东人数信息
*/
func (this *EquityInfo) GetShareholderJson(c *gin.Context) {

	lib.WriteString(c, 200, "test")
}

/**
获取十大流通股东信息
*/
func (this *EquityInfo) GetTop10Json(c *gin.Context) {

	lib.WriteString(c, 200, "test")
}

/**
获取机构持股信息
*/
func (this *EquityInfo) GetOrganizationJson(c *gin.Context) {

	lib.WriteString(c, 200, "test")
}
