package newService

import (
	"haina.com/market/f9/services"
	"haina.com/share/logging"
)

var (
	symbolParam   string //证券代码（带字母）
	symbol        string //证券代码（不带字母)
	symboType     string //证券类型
	swlevelcode   string //行业代码
	swlevelname   string //行业名称
	compcode      string //公司代码
	compname      string //公司名称
	companyDetail string //公司详情
	//sumcompay     string //该行业的所有公司
)

type newData struct {
	commonService.BasicMess
	Data commonService.BasicErrorData `json:"data"`
}

func GetBigData(scode string) (newData, error) {
	var ne newData
	data, _ := commonService.GetCommonData(scode)
	if data.Status == 0 {
		ne.Status = data.Status
		ne.Message = data.Message
	} else {
		ne.Status = data.Status
		ne.Message = data.Message
		ne.Data.Symbol = data.Data.Symbol
		ne.Data.Stname = data.Data.Stname
		ne.Data.NLastPx = data.Data.NLastPx
		ne.Data.NPxChg = data.Data.NPxChg
		ne.Data.NPxChgRatio = data.Data.NPxChgRatio
		ne.Data.Score = data.Data.Score
		ne.Data.Prompt = data.Data.Prompt
	}
	logging.Info("data===%+v", data)
	return ne, nil
}
