// 利润表
package company

type ProfitsJson struct {
	Date int64 `json:"Date"`

	AAPC float64 `json:"AAPC"` //影响母公司净利润的调整项目
	AILs float64 `json:"AILs"` //资产减值损失
	AREp float64 `json:"AREp"` //分保费用
	BAEp float64 `json:"BAEp"` //业务及管理费
	BPAC float64 `json:"BPAC"` //归属于母公司所有者的净利润
	CoEp float64 `json:"CoEp"` //手续费及佣金支出
	CoRe float64 `json:"CoRe"` //手续费及佣金收入
	CORe float64 `json:"CORe"` //营业成本
	DPES float64 `json:"DPES"` //稀释每股收益
	EPS  float64 `json:"EPS"`  //基本每股收益
	FnEp float64 `json:"FnEp"` //财务费用
	ICEp float64 `json:"ICEp"` //保险手续费及佣金支出
	IDEp float64 `json:"IDEp"` //保单红利支出
	InRe float64 `json:"InRe"` //利息收入
	ItEp float64 `json:"ItEp"` //利息支出
	ITEp float64 `json:"ITEp"` //所得税费用
	MgEp float64 `json:"MgEp"` //管理费用
	MIIn float64 `json:"MIIn"` //少数股东损益
	NCoE float64 `json:"NCoE"` //手续费及佣金净收入
	NInR float64 `json:"NInR"` //利息净收入
	NOEp float64 `json:"NOEp"` //营业外支出
	NORe float64 `json:"NORe"` //营业外收入
	NtIn float64 `json:"NtIn"` //净利润
	OATx float64 `json:"OATx"` //营业税金及附加
	OCOR float64 `json:"OCOR"` //营业总成本
	OOCs float64 `json:"OOCs"` //其他营业成本
	OpEp float64 `json:"OpEp"` //营业支出
	OpPr float64 `json:"OpPr"` //营业利润
	OpRe float64 `json:"OpRe"` //营业收入
	SaEp float64 `json:"SaEp"` //销售费用
	SAPC float64 `json:"SAPC"` //影响母公司净利润的特殊项目
	TOpR float64 `json:"TOpR"` //营业总收入
	ToPr float64 `json:"ToPr"` //利润总额
}

func NewProfitsJson() *ProfitsJson {
	return &ProfitsJson{}
}
