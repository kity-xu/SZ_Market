package publish

import (
	"haina.com/share/lib"

	"ProtocolBuffer/format/redis/pbdef/securitytable"

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
		lib.WriteString(c, 40004, nil)
		return
	}
}

func (this *SecurityTable) GetJson(c *gin.Context) {
	securitytable, err := publish.NewSecurityTable().GetSecurityTable()
	if err != nil {
		lib.WriteString(c, 40002, nil)
		return
	}
	lib.WriteString(c, 200, securitytable)
}

func (this *SecurityTable) GetPB(c *gin.Context) {
	var reply securitytable.ReplySecurityCodeTable
	securitytable, err := publish.NewSecurityTable().GetSecurityTable()
	if err != nil {
		reply.Code = 40002
	} else {
		reply.Code = 200
		reply.Stable = securitytable
	}
	v, err := proto.Marshal(&reply)
	if err != nil {
		logging.Info("%v", err)
		return
	}
	lib.WriteData(c, v)
}
