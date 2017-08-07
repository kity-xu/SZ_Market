//市场代码表
package security

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"

	"bytes"
	"compress/zlib"
	"encoding/json"

	"github.com/gin-gonic/gin"
	. "haina.com/market/hqpublish/controllers"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/publish/security"
	"haina.com/share/logging"
)

type SecurityTable struct {
}

func NewSecurityTable() *SecurityTable {
	return &SecurityTable{}
}

//A 股市场代码表
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
	table, err := security.NewSecurityNameTable().GetSecurityTableAStock()
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, 40002, nil)
		return
	}

	// 压缩处理
	byt, err := json.Marshal(table)
	if err != nil {
		logging.Error("date zib error:%v", err)
		WriteJson(c, 40002, nil)
		return
	}
	var in bytes.Buffer
	w := zlib.NewWriter(&in)
	w.Write(byt)
	w.Close()

	data := in.Bytes()
	WriteJson(c, 200, &data)
}

func (this *SecurityTable) GetPB(c *gin.Context) {
	table, err := security.NewSecurityNameTable().GetSecurityTableAStock()
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_MARKET_SECURITY_ASTOCK, table)
}

//单市场股票代码表
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
		return
	}
}

func (this *SecurityTable) PostJson(c *gin.Context) {
	request := &protocol.RequestMarketSecurityNameTable{}
	code, err := RecvAndUnmarshalJson(c, 1024, request)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, code, nil)
		return
	}
	if request.MarketID == 0 {
		WriteJson(c, 40004, nil)
	}

	table, err := security.NewSecurityNameTable().GetSecurityTable(request)
	if err != nil {
		logging.Error("post json %v", err)
		WriteJson(c, 40002, nil)
		return
	}
	//	// 压缩处理
	//	byt, err := json.Marshal(table)
	//	if err != nil {
	//		logging.Error("date zib error:%v", err)
	//		WriteJson(c, 40002, nil)
	//		return
	//	}
	//	var in bytes.Buffer
	//	w := zlib.NewWriter(&in)
	//	w.Write(byt)
	//	w.Close()

	//	data := in.Bytes()
	//	WriteJson(c, 200, &data)
	WriteJson(c, 200, table)
}
func (this *SecurityTable) PostPB(c *gin.Context) {
	request := &protocol.RequestMarketSecurityNameTable{}

	code, err := RecvAndUnmarshalPB(c, 1024, request)
	if err != nil {
		logging.Error("post pb %v", err)
		WriteDataErrCode(c, code)
		return
	}

	if request.MarketID == 0 {
		WriteJson(c, 40004, nil)
	}

	table, err := security.NewSecurityNameTable().GetSecurityTable(request)
	if err != nil {
		logging.Error("%v", err)
		WriteDataErrCode(c, 40002)
		return
	}
	WriteDataPB(c, protocol.HAINA_PUBLISH_CMD_ACK_MARKET_SECURITY, table)
}
