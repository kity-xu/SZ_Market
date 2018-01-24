package fund

import (
	"haina.com/market/f9/services/fundService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetFundData(c *gin.Context) {
	scode := c.Query("scode")
	data, err := fundService.GetFundData(scode)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
