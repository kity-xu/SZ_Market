// 证券快照
package publish

import (
	"encoding/json"
	"io"
	"net/http"

	"haina.com/share/lib"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	"haina.com/share/logging"

	"ProtocolBuffer/format/snap"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
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
		//lib.WriteString(c, 40004, nil)
		return
	}
}

func (this *StockSnapshot) PostJson(c *gin.Context) {
	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}

	var request snap.RequestSnap
	if err := json.Unmarshal(buf, &request); err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)
	data, err := publish.NewStockSnapshot().GetStockSnapshot(&request)
	if err != nil {
		logging.Error("%v", err)
		reply := snap.ReplySnap{
			Code: 40002,
		}
		c.JSON(http.StatusOK, reply)
		return
	}

	reply := snap.ReplySnap{
		Code: 200,
		Data: data,
	}

	c.JSON(http.StatusOK, reply)
}

func (this *StockSnapshot) PostPB(c *gin.Context) {
	var (
		replypb []byte
		err     error
		request snap.RequestSnap
	)

	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}
	if err := proto.Unmarshal(buf, &request); err != nil {
		logging.Error("PB Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)
	data, err := publish.NewStockSnapshot().GetStockSnapshot(&request)
	if err != nil {
		reply := snap.ReplySnap{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}
	}
	reply := &snap.ReplySnap{
		Code: 200,
		Data: data,
	}
	replypb, err = proto.Marshal(reply)
	if err != nil {
		reply := snap.ReplySnap{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}

	}
	lib.WriteData(c, replypb)
}
