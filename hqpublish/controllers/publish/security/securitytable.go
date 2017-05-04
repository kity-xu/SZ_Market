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

type SecurityTable struct {
}

func NewSecurityTable() *SecurityTable {
	return &SecurityTable{}
}

func (this *SecurityTable) POST(c *gin.Context) {
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

func (this *SecurityTable) PostJson(c *gin.Context) {
	var (
		request securitytable.RequestMarketSecurityCodeTable
	)
	buf, err := getRequestData(c)
	if err != nil && err != io.EOF {
		logging.Error("%v", err)
		return
	}

	if err := json.Unmarshal(buf, &request); err != nil {
		logging.Error("Json Request Unmarshal: %v", err)
		return
	}

	table, _, err := security.NewSecurityNameTable().GetSecurityTable(request.MarketID)
	if err != nil {
		logging.Error("%v", err)
		lib.WriteString(c, 40002, err.Error())
		return
	}

	reply := &securitytable.ReplyMarketSecurityCodeTable{
		Code: 200,
		Data: table,
	}

	c.JSON(http.StatusOK, reply)
}
func (this *SecurityTable) PostPB(c *gin.Context) {
	var (
		replypb []byte
		err     error
		request securitytable.RequestMarketSecurityCodeTable
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

	table, _, err := security.NewSecurityNameTable().GetSecurityTable(request.MarketID)
	if err != nil {
		logging.Error("%v", err)
		return
	}

	reply := &securitytable.ReplyMarketSecurityCodeTable{
		Code: 200,
		Data: table,
	}

	replypb, err = proto.Marshal(reply)
	if err != nil {
		reply := securitytable.ReplyMarketSecurityCodeTable{
			Code: 40002,
		}
		replypb, err = proto.Marshal(&reply)
		if err != nil {
			logging.Error("pb marshal error: %v", err)
		}

	}
	lib.WriteData(c, replypb)

}

/*******************************************************************/
func getRequestData(c *gin.Context) ([]byte, error) {
	temp := make([]byte, 1024)
	n, err := c.Request.Body.Read(temp)
	if err != nil && err != io.EOF {
		logging.Error("Body Read: %v", err)
		return nil, err
	}
	//logging.Info("\nBody len %d\n%s", n, temp[:n])
	logging.Info("Body len %d", n)
	return temp[:n], nil

}
