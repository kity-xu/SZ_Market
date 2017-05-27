//证券基本信息
package security

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	. "haina.com/market/hqpublish/controllers"

	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/security"
	"haina.com/share/logging"
)

type SecurityInfo struct {
}

func NewSecurityInfo() *SecurityInfo {
	return &SecurityInfo{}
}

func (this *SecurityInfo) POST(c *gin.Context) {
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
func (this *SecurityInfo) PostJson(c *gin.Context) {
	request := &protocol.RequestSingleSecurity{}
	code, err := RecvAndUnmarshalJson(c, 1024, request)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}

	if request.SID == 0 {
		WriteJson(c, 40004, nil)
		return
	}

	info, err := security.NewSecurityInfo().GetSecurityBasicInfo(request)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, 40002, nil)
		return
	}
	WriteJson(c, 200, info)

}
func (this *SecurityInfo) PostPB(c *gin.Context) {
	request := &protocol.RequestSingleSecurity{}

	code, err := RecvAndUnmarshalPB(c, 1024, request)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}

	info, err := security.NewSecurityInfo().GetSecurityBasicInfo(request)
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_SINGLE_SECURITY, info)
}
