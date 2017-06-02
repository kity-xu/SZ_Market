//股票静态信息
package security

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	. "haina.com/market/hqpublish/controllers"

	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/security"
	"haina.com/share/logging"
)

type SecurityStatic struct {
}

func NewSecurityStatic() *SecurityStatic {
	return &SecurityStatic{}
}

func (this *SecurityStatic) POST(c *gin.Context) {
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
func (this *SecurityStatic) PostJson(c *gin.Context) {
	request := &protocol.RequestSecurityStatic{}
	code, err := RecvAndUnmarshalJson(c, 1024, request)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}

	sinfo, err := security.NewSecurityStatic().GetSecurityStatic(request)
	if err != nil {
		logging.Error("%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, sinfo)

}
func (this *SecurityStatic) PostPB(c *gin.Context) {
	request := &protocol.RequestSecurityStatic{}

	code, err := RecvAndUnmarshalPB(c, 1024, request)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}

	sinfo, err := security.NewSecurityStatic().GetSecurityStatic(request)
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_SECURITYSTATIC, sinfo)
}
