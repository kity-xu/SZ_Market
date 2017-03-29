// 利润表
package company

import (
	"time"

	"haina.com/market/finance/models/finchina"
)

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

//------------------------------------------------------------------------------
type Profits struct {
}

func NewProfits() *Profits {
	return &Profits{}
}

func (this *Profits) getList(scode string, report_type int, per_page int, page int) ([]Profits, error) {
	return nil, nil
}
func (this *Profits) GetList(scode string, report_type int, per_page int, page int) ([]Profits, error) {
	return nil, nil
}
func (this *Profits) getJson(scode string, report_type int, per_page int, page int) ([]ProfitsJson, error) {
	return NewFinChinaProfits().getJson(scode, report_type, per_page, page)
}
func (this *Profits) GetJson(scode string, report_type int, per_page int, page int) (*RespFinAnaJson, error) {
	ls, err := this.getJson(scode, report_type, per_page, page)
	if err != nil {
		return nil, err
	}
	jsn := &RespFinAnaJson{
		MU:     "人民币元",
		AS:     "新会计准则",
		Length: len(ls),
		List:   ls,
	}
	return jsn, nil
}

//------------------------------------------------------------------------------
type FinChinaProfits struct {
}

func NewFinChinaProfits() *FinChinaProfits {
	return &FinChinaProfits{}
}

func (this *FinChinaProfits) getJson(scode string, report_type int, per_page int, page int) ([]ProfitsJson, error) {
	sli := make([]ProfitsJson, 0, per_page)
	ls, err := finchina.NewProfits().GetList(scode, report_type, per_page, page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		one := ProfitsJson{
			AILs: v.ASSEIMPALOSS.Float64,
			AREp: v.REINEXPE.Float64,
			BPAC: v.PARENETP.Float64,
			CoEp: v.POUNEXPE.Float64,
			CoRe: v.POUNINCO.Float64,
			CORe: v.BIZCOST.Float64,
			DPES: v.DILUTEDEPS.Float64,
			EPS:  v.BASICEPS.Float64,
			FnEp: v.FINEXPE.Float64,
			IDEp: v.POLIDIVIEXPE.Float64,
			InRe: v.INTEINCO.Float64,
			ItEp: v.INTEEXPE.Float64,
			ITEp: v.INCOTAXEXPE.Float64,
			MgEp: v.MANAEXPE.Float64,
			MIIn: v.MINYSHARRIGH.Float64,
			NOEp: v.NONOEXPE.Float64,
			NORe: v.NONOREVE.Float64,
			NtIn: v.NETPROFIT.Float64,
			OATx: v.BIZTAX.Float64,
			OCOR: v.BIZTOTCOST.Float64,
			OpPr: v.PERPROFIT.Float64,
			OpRe: v.BIZINCO.Float64,
			SaEp: v.SALESEXPE.Float64,
			TOpR: v.BIZTOTINCO.Float64,
			ToPr: v.TOTPROFIT.Float64,
		}

		if v.ENDDATE.Valid {
			tm, err := time.Parse("20060102", v.ENDDATE.String)
			if err != nil {
				return nil, err
			}
			one.Date = tm.Unix()
		}

		sli = append(sli, one)
	}
	return sli, nil
}
