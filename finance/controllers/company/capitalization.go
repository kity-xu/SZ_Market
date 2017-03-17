package company

import (
	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
)

type CapitalizationInfo struct {
}

func NewCapitalizationInfo() *CapitalizationInfo {
	return &CapitalizationInfo{}
}

/**
获取股本结构信息
*/
func (this *CapitalizationInfo) GetStructureJson(c *gin.Context) {

	lib.WriteString(c, 200, "test")
}

/**
获取股本变动信息
*/
func (this *CapitalizationInfo) GetChangesJson(c *gin.Context) {

	lib.WriteString(c, 200, "test")
}
