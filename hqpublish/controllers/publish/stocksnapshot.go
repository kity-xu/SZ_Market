// 证券快照
package publish

import (
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	"haina.com/share/logging"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"github.com/gin-gonic/gin"
)

type StockSnapshot struct{}

func NewStockSnapshot() *StockSnapshot {
	return &StockSnapshot{}
}

func (this *StockSnapshot) POST(c *gin.Context) {
	replayfmt := c.Query(models.CONTEXT_FORMAT)
	if len(replayfmt) == 0 {
		replayfmt = "pb" // 默认
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

func (this *StockSnapshot) PostJson(c *gin.Context) {
	var req protocol.RequestSnapshot

	code, err := RecvAndUnmarshalJson(c, 1024, &req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request %+v", req)
	if req.SID == 0 {
		WriteJson(c, 40004, nil)
		return
	}

	data, err := publish.NewStockSnapshot().GetStockSnapshot(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, data)
}

func (this *StockSnapshot) PostPB(c *gin.Context) {
	var req protocol.RequestSnapshot

	code, err := RecvAndUnmarshalPB(c, 1024, &req)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}
	logging.Info("Request %+v", req)
	if req.SID == 0 {
		WriteDataErrCode(c, 40004)
		return
	}

	data, err := publish.NewStockSnapshot().GetStockSnapshot(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_SNAPSHOT, data)
}
