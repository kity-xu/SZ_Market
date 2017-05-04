package security

import (
	"ProtocolBuffer/format/securitytable"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/security"
	"haina.com/share/lib"
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
	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}

	var request securitytable.RequestSingleSecurity
	if err := json.Unmarshal(buf, &request); err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}

	single, ok := getSingleSecurityBasicInfoFromRedis(request.SID)
	if !ok {
		logging.Error("Can't find the security information...")
		lib.WriteString(c, 200, "Can't find the security information...")
		return
	}
	reply := &securitytable.ReplySecuritySingle{
		Code: 200,
		Data: single,
	}

	c.JSON(http.StatusOK, reply)

}
func (this *SecurityInfo) PostPB(c *gin.Context) {
	var (
		request securitytable.RequestSingleSecurity
		reply   securitytable.ReplySecuritySingle
		replypb []byte
		err     error
	)

	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}

	if err := proto.Unmarshal(buf, &request); err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}
	logging.Info("Request Data: %+v", request)

	single, ok := getSingleSecurityBasicInfoFromRedis(request.SID)
	reply.Data = single
	if !ok {
		reply.Code = 40002
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}
		logging.Error("Can't find the security information...")
		return
	}

	reply.Code = 200

	//转PB
	replypb, err = proto.Marshal(&reply)
	if err != nil {
		reply := securitytable.ReplySecuritySingle{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}
		return
	}
	lib.WriteData(c, replypb)

}

func getSingleSecurityBasicInfoFromRedis(sid int32) (*securitytable.SecuritySingle, bool) {
	var (
		count int
		ok    bool
		reply securitytable.SecuritySingle
	)

	marketID := sid / 1000000 * 1000000
	table, _, err := security.NewSecurityNameTable().GetSecurityTable(marketID)
	if err != nil {
		logging.Error("%v", err)
		return &reply, ok
	}

	for _, v := range table.List {
		if v.NSID == sid {
			reply.SID = sid
			reply.Single = v
			break
		}
		count++
	}
	if count < len(table.List) {
		ok = true
	} else {
		ok = false
	}

	logging.Debug("count:%v", count)

	return &reply, ok
}
