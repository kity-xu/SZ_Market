// 资产负债表
package company

type LiabilitiesJson struct {
	Date int64 `json:"Date"`

	//资产
	AcRe float64 `json:"AcRe"` //应收账款
	CuDe float64 `json:"CuDe"` //客户资金存款
	DeMg float64 `json:"DeMg"` //存出保证金
	DFAs float64 `json:"DFAs"` //衍生金融资产
	DfSv float64 `json:"DfSv"` //定期存款
	DiRe float64 `json:"DiRe"` //应收股利
	DTAs float64 `json:"DTAs"` //递延所得税资产
	FAFS float64 `json:"FAFS"` //可供出售金融资产
	FAHT float64 `json:"FAHT"` //交易性金融资产
	FiAs float64 `json:"FiAs"` //固定资产
	GWil float64 `json:"GWil"` //商誉
	HTMI float64 `json:"HTMI"` //持有至到期投资
	IbDe float64 `json:"IbDe"` //存放同业款项
	InAs float64 `json:"InAs"` //无形资产
	InRe float64 `json:"InRe"` //应收利息
	LdLt float64 `json:"LdLt"` //拆出资金
	LTAR float64 `json:"LTAR"` //长期应收款
	LTEI float64 `json:"LTEI"` //长期股权投资
	LTPE float64 `json:"LTPE"` //长期待摊费用
	Metl float64 `json:"Metl"` //贵金属
	MnFd float64 `json:"MnFd"` //货币资金
	NoRe float64 `json:"NoRe"` //应收票据
	OtRe float64 `json:"OtRe"` //其他应收款
	PrEx float64 `json:"PrEx"` //待摊费用
	Prpy float64 `json:"Prpy"` //预付款项
	REFI float64 `json:"REFI"` //投资性房地产
	ToAs float64 `json:"ToAs"` //资产总计

	//负债
	AcEx  float64 `json:"AcEx"`  //预提费用
	AcPy  float64 `json:"AcPy"`  //应付账款
	AdRE  float64 `json:"AdRE"`  //预收款项
	BdPy  float64 `json:"BdPy"`  //应付债券
	CmPy  float64 `json:"CmPy"`  //应付手续费及佣金
	DETLb float64 `json:"DETLb"` //递延所得税负债
	DfIn  float64 `json:"DfIn"`  //递延收益
	DFLb  float64 `json:"DFLb"`  //衍生金融负债
	DpCl  float64 `json:"DpCl"`  //吸收存款
	DpFB  float64 `json:"DpFB"`  //同业及其他金融机构存放款项
	DvPy  float64 `json:"DvPy"`  //应付股利
	FASR  float64 `json:"FASR"`  //卖出回购金融资产款
	InPy  float64 `json:"InPy"`  //应付利息
	LnFB  float64 `json:"LnFB"`  //拆入资金
	LnFC  float64 `json:"LnFC"`  //向中央银行借款
	LTBw  float64 `json:"LTBw"`  //长期借款
	LTPy  float64 `json:"LTPy"`  //长期应付款
	NCL1  float64 `json:"NCL1"`  //一年内到期的非流动负债
	NtPy  float64 `json:"NtPy"`  //应付票据
	PCSc  float64 `json:"PCSc"`  //永续债
	PlLn  float64 `json:"PlLn"`  //质押借款
	PrSk  float64 `json:"PrSk"`  //优先股
	SaPy  float64 `json:"SaPy"`  //应付职工薪酬
	SBPy  float64 `json:"SBPy"`  //应付短期债券
	STLn  float64 `json:"STLn"`  //短期借款
	TaLb  float64 `json:"TaLb"`  //负债合计
	TFLb  float64 `json:"TFLb"`  //交易性金融负债
	TxPy  float64 `json:"TxPy"`  //应交税费

	//所有者权益
	BPCOEAI float64 `json:"BPCOEAI"` //归属于母公司所有者权益调整项目
	BPCOESI float64 `json:"BPCOESI"` //归属于母公司所有者权益特殊项目
	BPCSET  float64 `json:"BPCSET"`  //归属于母公司股东权益合计
	CDFCS   float64 `json:"CDFCS"`   //外币报表折算差额
	CpSp    float64 `json:"CpSp"`    //资本公积
	GRPr    float64 `json:"GRPr"`    //一般风险准备
	LEAI    float64 `json:"LEAI"`    //负债和权益调整项目
	LESI    float64 `json:"LESI"`    //负债和权益特殊项目
	MiIt    float64 `json:"MiIt"`    //少数股东权益
	OEAI    float64 `json:"OEAI"`    //所有者权益调整项目
	OEIn    float64 `json:"OEIn"`    //其他权益工具
	OESET   float64 `json:"OESET"`   //所有者权益（或股东权益）合计
	OtCI    float64 `json:"OtCI"`    //其他综合收益
	PCSe    float64 `json:"PCSe"`    //永续债
	PICa    float64 `json:"PICa"`    //实收资本（或股本）
	PrSc    float64 `json:"PrSc"`    //优先股
	SpRs    float64 `json:"SpRs"`    //盈余公积
	TLSE    float64 `json:"TLSE"`    //负债和所有者权益（或股东权益）总计
	TrSc    float64 `json:"TrSc"`    //库存股
	UdPr    float64 `json:"UdPr"`    //未分配利润
}

func NewLiabilitiesJson() *LiabilitiesJson {
	return &LiabilitiesJson{}
}
