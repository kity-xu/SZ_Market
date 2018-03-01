package basic

import (
	. "haina.com/market/f9/controllers"
	"haina.com/market/f9/services/newService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetNewData(c *gin.Context) {
	sid, err := Param_Norm_Sid(c)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	data, err := newService.GetBigData(sid)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
