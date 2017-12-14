package publish2

import (
	"github.com/gin-gonic/gin"

	"haina.com/market/hqpublish/models"
)

type ReportStatistics struct {
}

func NewReportStatistics() *ReportStatistics {
	return &ReportStatistics{}
}

func (this *ReportStatistics) POST(c *gin.Context) {
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

func (this *ReportStatistics) PostJson(c *gin.Context) {

}
func (this *ReportStatistics) PostPB(c *gin.Context) {

}
