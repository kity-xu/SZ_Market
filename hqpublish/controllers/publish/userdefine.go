//自选股
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
)

type UserDefine struct{}

func NewUserDefine() *UserDefine {
	return &UserDefine{}
}

func (this *UserDefine) POST(c *gin.Context) {
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

func (this *UserDefine) PostJson(c *gin.Context) {
	var req protocol.RequestUserdef
	code, err := RecvAndUnmarshalJson(c, 1024, &req)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	logging.Info("Request %+v", req)

	reply, err := publish.NewUserDefine(publish.REDISKEY_SORT_KDAY_H).GetSecurityUserdefine(&req)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, reply)
}

func (this *UserDefine) PostPB(c *gin.Context) {
	var req protocol.RequestUserdef
	code, err := RecvAndUnmarshalPB(c, 1024, &req)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}
	logging.Info("Request %+v", req)

	reply, err := publish.NewUserDefine(publish.REDISKEY_SORT_KDAY_H).GetSecurityUserdefine(&req)
	if err != nil {
		logging.Error("%v", err.Error())
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_USERDEFSEC, reply)
}
