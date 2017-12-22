//板块
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

type StockBlock struct{}

func NewStockBlock() *StockBlock {
	return &StockBlock{}
}

func (this *StockBlock) POST(c *gin.Context) {
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

func (this *StockBlock) PostJson(c *gin.Context) {
	var req protocol.RequestBlock
	code, err := RecvAndUnmarshalJson(c, 1024, &req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request %+v", req)

	reply, err := publish.NewBlock(publish.REDIS_KEY_CACHE_BLOCK).GetBlockReplyByRequest(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (this *StockBlock) PostPB(c *gin.Context) {
	var req protocol.RequestBlock
	code, err := RecvAndUnmarshalPB(c, 1024, &req)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}
	logging.Info("Request %+v", req)

	reply, err := publish.NewBlock(publish.REDIS_KEY_CACHE_BLOCK).GetBlockReplyByRequest(&req)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_SET, reply)
}

// -------------- hq2.2 附加 上证、深证、创业、中小板块-----------//
func (this *StockBlock) MulPOST(c *gin.Context) {
	var req protocol.RequestBlock
	code, err := RecvAndUnmarshalJson(c, 1024, &req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}

	if req.TypeID > 4 { // 5:上证主板 6:深证主板  7:中小板 8:创业板
		res, err := this.sort(c, &req)
		if err != nil {
			WriteJson(c, 40002, nil)
			return
		}
		WriteJson(c, 200, res)
		return
	}
	reply, err := publish.NewBlock(publish.REDIS_KEY_CACHE_BLOCK).GetBlockReplyByRequest(&req)
	if err != nil {
		logging.Error("板块排序 |%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (s *StockBlock) sort(c *gin.Context, o *protocol.RequestBlock) (interface{}, error) {
	req := &protocol.RequestSort{
		SetID:   o.TypeID,
		FieldID: o.FieldID,
		Begin:   o.Begin,
		Num:     o.Num,
	}

	var res struct {
		TypeID  int32
		FieldID int32
		Total   int32
		Num     int32
		List    []*protocol.TagStockSortInfo
	}

	reply, err := publish.NewSort(publish.REDISKEY_SORT_KDAY_H).GetPayloadSort(req)
	if err != nil {
		logging.Error("指数排序 |%v", err)
		return nil, err
	}

	res.TypeID = reply.SetID
	res.FieldID = reply.FieldID
	res.Total = reply.Total
	res.Num = reply.Num
	res.List = reply.List

	return &res, nil
}
