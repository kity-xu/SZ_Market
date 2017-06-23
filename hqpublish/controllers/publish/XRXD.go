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

var (
	_ = fmt.Print
	_ = lib.WriteString
	_ = logging.Error
	_ = protocol.MarketStatus{}
	_ = publish.MarketStatus{}
	_ = strings.ToLower
	_ = proto.Marshal
	_ = json.Marshal
	_ = io.ReadFull
)

type XRXD struct{}

func NewXRXD() *XRXD {
	return &XRXD{}
}

func (this *XRXD) POST(c *gin.Context) {
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

func (this *XRXD) PostJson(c *gin.Context) {
	var req protocol.RequestXRXD
	code, err := RecvAndUnmarshalJson(c, 1024, &req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request %+v", req)
	if req.SID == 0 || req.TimeIndex < 0 || req.Num < 0 {
		WriteJson(c, 40004, nil)
		return
	}

	payload, err := publish.NewXRXD().GetXRXDObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, payload)
}

func (this *XRXD) PostPB(c *gin.Context) {
	var req protocol.RequestXRXD
	code, err := RecvAndUnmarshalPB(c, 1024, &req)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}
	logging.Info("Request %+v", req)
	if req.SID == 0 || req.TimeIndex < 0 || req.Num < 0 {
		WriteDataErrCode(c, 40004)
		return
	}

	payload, err := publish.NewXRXD().GetXRXDObj(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_XRXD, payload)

}
