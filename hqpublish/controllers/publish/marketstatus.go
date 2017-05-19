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
var _ = proto.Marshal
var _ = protocol.MarketStatus{}
var _ = publish.MarketStatus{}
var _ = strings.ToLower

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
		this.GetJson(c)
	case "pb":
		this.GetPB(c)
	default:
		return
	}
}

func (this *MarketStatus) RecvAndCheckJson(c *gin.Context) (*protocol.RequestMarketStatus, int, error) {
	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		return nil, 40004, err
	}

	var req protocol.RequestMarketStatus
	if err := json.Unmarshal(buf, &req); err != nil {
		return nil, 40004, err
	}
	logging.Info("Request Data: %+v", req)

	if int(req.Num) != len(req.MarketIDList) {
		logging.Error("Num %d, List len %d", req.Num, len(req.MarketIDList))
		return nil, 40002, ERROR_REQUEST_PARAM
	}
	return &req, 0, nil
}
func (this *MarketStatus) RecvAndCheckPB(c *gin.Context) (*protocol.RequestMarketStatus, int, error) {
	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		return nil, 40004, err
	}

	var req protocol.RequestMarketStatus
	if err := proto.Unmarshal(buf, &req); err != nil {
		return nil, 40004, err
	}
	logging.Info("Request Data: %+v", req)

	if int(req.Num) != len(req.MarketIDList) {
		logging.Error("Num %d, List len %d", req.Num, len(req.MarketIDList))
		return nil, 40002, ERROR_REQUEST_PARAM
	}
	return &req, 0, nil
}

func (this *MarketStatus) GetJson(c *gin.Context) {
	req, code, err := this.RecvAndCheckJson(c)
	if err != nil {
		logging.Error("GetJson %v", err)
		lib.WriteString(c, code, nil)
		return
	}
	res, err := publish.NewMarketStatus().GetPayloadObj(req)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	lib.WriteString(c, 200, res)
}

func (this *MarketStatus) GetPB(c *gin.Context) {
	req, code, err := this.RecvAndCheckPB(c)
	if req == nil {
		logging.Error("GetJson %v", err)
		data, err := MakeRespDataByBytes(code, 0, nil)
		if err != nil {
			return
		}
		lib.WriteData(c, data)
		return
	}
	res, err := publish.NewMarketStatus().GetPayloadPB(req)
	if err != nil {
		logging.Error("%v", err)
		data, err := MakeRespDataByBytes(40002, 0, nil)
		if err != nil {
			return
		}
		lib.WriteData(c, data)
		return
	}

	data, err := MakeRespDataByBytes(200, int(protocol.HAINA_PUBLISH_CMD_ACK_MARKET_STATUS), res)
	if err != nil {
		return
	}

	//  // 解码查看验证数据
	//	var a protocol.PayloadMarketStatus
	//	if err := proto.Unmarshal(res, &a); err != nil {
	//		logging.Error("%v", err)
	//	}
	//	logging.Info("%+v", a)

	lib.WriteData(c, data)
}
