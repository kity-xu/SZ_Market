package isNew

import (
	"haina.com/market/f9/services/newService"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

func GetNewData(c *gin.Context) {
	scode := c.Query("scode")
	data, err := newService.GetBigData(scode)
	if err != nil {
		logging.Info(err.Error())
	}
	lib.WriteString(c, 200, data)
}
