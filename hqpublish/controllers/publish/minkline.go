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

type MinKLine struct{}

func NewMinKLine() *MinKLine {
	return &MinKLine{}
}

func (this *MinKLine) POST(c *gin.Context) {
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

func (this *MinKLine) PostJson(c *gin.Context) {
	var request protocol.RequestMinK
	code, err := RecvAndUnmarshalJson(c, 1024, &request)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request %+v", request)

	if request.BeginTime > 1502 {
		logging.Error("%v %d", ERROR_KLINE_BEGIN_TIME, request.BeginTime)
		WriteJson(c, 40004, nil)
		return
	}

	js, err := publish.NewMinKLine().GetMinKJson(&request)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteDataJson(c, js)
}

func (this *MinKLine) PostPB(c *gin.Context) {
	var request protocol.RequestMinK
	code, err := RecvAndUnmarshalPB(c, 1024, &request)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}
	logging.Info("Request %+v", request)

	if request.BeginTime > 1502 {
		logging.Error("%v", ERROR_KLINE_BEGIN_TIME)
		WriteJson(c, 40004, nil)
		return
	}

	data, err := publish.NewMinKLine().GetMinKPB(&request)
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataBinary(c, data)
}
