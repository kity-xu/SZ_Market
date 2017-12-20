package f10

import (
	"haina.com/share/lib"

	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish/f10"

	"haina.com/share/logging"
)

type HN_F10_Mobile struct {
}

func NewHN_F10_Mobile() *HN_F10_Mobile {
	return &HN_F10_Mobile{}
}

// F10 首页
func (this *HN_F10_Mobile) GetF10_Mobile(c *gin.Context) {

	var _param struct {
		Scode int `json:"sid" binding:"required"`
	}

	if err := c.BindJSON(&_param); err != nil {
		logging.Debug("Bind Json | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	f10, err := f10.F10Mobile(_param.Scode)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	lib.WriteString(c, 200, f10)
}
