package f10

import (
	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type ShareholdersTop10 struct {
}

func NewShareholdersTop10() *ShareholdersTop10 {
	return &ShareholdersTop10{}
}

// 历史股本变动
func (this *ShareholdersTop10) GetShareholdersTop10(c *gin.Context) {

	var _param struct {
		Scode string `json:"sid" binding:"required"`
		Count int    `json:"count"`
		HType int    `json:"htype" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	//
}
