// 涨跌统计接口
package publish2

import (
	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models/publish2"
	"haina.com/share/lib"
)

type CapStatistics struct {
}

func NewStatistics() *CapStatistics {
	return &CapStatistics{}
}

func (CapStatistics) GET(c *gin.Context) {
	res, err := publish2.NewMarketsStatistics().GetMarketsStatistics()
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	lib.WriteString(c, 200, res)
}
