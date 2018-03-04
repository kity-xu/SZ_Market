package summary

import (
	. "haina.com/market/f9/controllers"
	"haina.com/market/f9/services/fundService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetFundData(c *gin.Context) {
	sid, err := Param_Norm_Sid(c)
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	data, err := fundService.GetFundData(sid)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
