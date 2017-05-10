package publish

import (
	"haina.com/share/lib"

	"ProtocolBuffer/format/securitytable"

	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish"

	"haina.com/share/logging"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
)

type SecurityTable struct{}

func NewSecurityTable() *SecurityTable {
	return &SecurityTable{}
}

func (this *SecurityTable) GET(c *gin.Context) {
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

func (this *SecurityTable) GetJson(c *gin.Context) {
	securitytable, err := publish.NewSecurityTable().GetSecurityTable()
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}
	lib.WriteString(c, 200, securitytable)
}

func (this *SecurityTable) GetPB(c *gin.Context) {
	var (
		replypb []byte
		err     error
	)
	replypb, err = publish.NewSecurityTable().GetSecurityTableReplyBytes()
	if err != nil {
		logging.Error("%v", err)
		reply := securitytable.ReplySecurityCodeTable{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}
	}
	lib.WriteData(c, replypb)
}
