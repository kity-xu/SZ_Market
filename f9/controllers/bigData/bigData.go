package bigData

import (
	"haina.com/market/f9/services/bigService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetBigData(c *gin.Context) {
	scode := c.Query("scode")
	data, err := bigService.GetBigData(scode)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
