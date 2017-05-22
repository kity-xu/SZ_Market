package publish

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

var _ = fmt.Print
var _ = lib.WriteString
var _ = logging.Error
var _ = protocol.MarketStatus{}
var _ = publish.MarketStatus{}
var _ = strings.ToLower
var _ = proto.Marshal
var _ = json.Marshal
var _ = io.ReadFull

type MarketStatus struct{}

func NewMarketStatus() *MarketStatus {
	return &MarketStatus{}
}

func (this *MarketStatus) POST(c *gin.Context) {
	replayfmt := c.Query(models.CONTEXT_FORMAT)

	if len(replayfmt) == 0 {
		replayfmt = "pb"
	} else {
		replayfmt = strings.ToLower(replayfmt)
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

func (this *MarketStatus) PostJson(c *gin.Context) {
	var req protocol.RequestMarketStatus
	code, err := RecvAndUnmarshalJson(c, 1024, &req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	if int(req.Num) != len(req.MarketIDList) {
		logging.Error("Num %d, List len %d", req.Num, len(req.MarketIDList))
		WriteJson(c, 40002, nil)
		return
	}
	res, err := publish.NewMarketStatus().GetPayloadObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}

	WriteJson(c, 200, res)
}

func (this *MarketStatus) PostPB(c *gin.Context) {
	var req protocol.RequestMarketStatus
	code, err := RecvAndUnmarshalPB(c, 1024, &req)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}

	if int(req.Num) != len(req.MarketIDList) {
		logging.Error("Num %d, List len %d", req.Num, len(req.MarketIDList))
		WriteDataErrCode(c, 40002)
		return
	}

	res, err := publish.NewMarketStatus().GetPayloadPB(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}

	//	// 解码查看 验证数据
	//	{
	//		var a protocol.PayloadMarketStatus
	//		if err := proto.Unmarshal(res, &a); err != nil {
	//			logging.Error("%v", err)
	//		}
	//		logging.Info("%+v", a)
	//	}

	WriteDataBytes(c, protocol.HAINA_PUBLISH_CMD_ACK_MARKET_STATUS, res)
}
