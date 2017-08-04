package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

type OptionalSids struct{}

func NewOptionalSids() *OptionalSids {
	return &OptionalSids{}
}

func (this *OptionalSids) GET(c *gin.Context) {
	replayfmt := c.Query(models.CONTEXT_FORMAT)
	if len(replayfmt) == 0 {
		replayfmt = "pb" // 默认
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

func (this *OptionalSids) POST(c *gin.Context) {
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

func (this *OptionalSids) PostJson(c *gin.Context) {
	access_token := c.Query(models.ACCESS_TOKEN)

	req := &protocol.RequestOptstockPut{}
	code, err := RecvAndUnmarshalJson(c, 1024, req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}

	err = publish.NewMOptSids().OperationStockSids(req, access_token)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, nil)
}

func (this *OptionalSids) PostPB(c *gin.Context) {
	access_token := c.Query(models.ACCESS_TOKEN)

	req := &protocol.RequestOptstockPut{}
	if code, err := RecvAndUnmarshalPB(c, 1024, req); err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}

	err := publish.NewMOptSids().OperationStockSids(req, access_token)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_OPTSTOCK_GET, nil)
}

func (this *OptionalSids) GetJson(c *gin.Context) {
	access_token := c.Query(models.ACCESS_TOKEN)

	sidList, err := publish.NewMOptSids().SelectAllSidsByAccessToken(access_token)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, sidList)
}

func (this *OptionalSids) GetPB(c *gin.Context) {
	access_token := c.Query(models.ACCESS_TOKEN)

	sidList, err := publish.NewMOptSids().SelectAllSidsByAccessToken(access_token)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_OPTSTOCK_GET, sidList)
}
