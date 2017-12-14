package publish2

import (
	"github.com/gin-gonic/gin"

	"haina.com/market/hqpublish/models"
)

type ReportForecast struct {
}

func NewReportForecast() *ReportForecast {
	return &ReportForecast{}
}

func (this *ReportForecast) POST(c *gin.Context) {
	replayfmt := c.Query(models.CONTEXT_FORMAT)
	if len(replayfmt) == 0 {
		replayfmt = "json" // 默认
	}

	switch replayfmt {
	case "json":
		this.PostJson(c)
	case "pb":
		this.PostPB(c)
	default:
		return
	}
}

func (this *ReportForecast) PostJson(c *gin.Context) {

}
func (this *ReportForecast) PostPB(c *gin.Context) {

}
