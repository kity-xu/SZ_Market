// 个股详情
package publish

import (
	"ProtocolBuffer/projects/hqpublish/go/protocol"
	"strconv"

	"haina.com/market/hqpublish/models/publish/security"
	"haina.com/share/logging"
)

type PerSDetailM struct {
}

func NewPerSDetailM() *PerSDetailM {
	return &PerSDetailM{}
}

// 个股详情
func (this *PerSDetailM) GetPerSDtail(req *protocol.RequestPerSDetail) (*protocol.PayloadPerSDetail, error) {

	// 调用静态数据处理
	var reqsta protocol.RequestSecurityStatic
	reqsta.SID = int32(req.SID)
	sinfo, err := security.NewSecurityStatic().GetSecurityStatic(&reqsta)

	hid := strconv.Itoa(int(req.SID))
	logging.Info("===============%v", hid, sinfo)
	//var err error
	return nil, err
}
