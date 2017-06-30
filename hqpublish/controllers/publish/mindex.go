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

type MIndex struct{}

func NewMIndex() *MIndex {
	return &MIndex{}
}

func (this *MIndex) POST(c *gin.Context) {
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

func (this *MIndex) PostJson(c *gin.Context) {
	payload, err := publish.NewMIndex().GetMIndexObj()
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, payload)
}

func (this *MIndex) PostPB(c *gin.Context) {
	payload, err := publish.NewMIndex().GetMIndexObj()
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_MINDEX, payload)
}
