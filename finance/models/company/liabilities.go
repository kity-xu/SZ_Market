// 资产负债表
package company

import (
	"time"

	"haina.com/market/finance/models/finchina"
)

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

//------------------------------------------------------------------------------
type Liabilities struct {
}

func NewLiabilities() *Liabilities {
	return &Liabilities{}
}

func (this *Liabilities) getList(scode string, report_type int, per_page int, page int) ([]Liabilities, error) {
	return nil, nil
}
func (this *Liabilities) GetList(scode string, report_type int, per_page int, page int) ([]Liabilities, error) {
	return nil, nil
}
func (this *Liabilities) getJson(scode string, report_type int, per_page int, page int) ([]LiabilitiesJson, error) {
	return NewFinChinaLiabilities().getJson(scode, report_type, per_page, page)
}
func (this *Liabilities) GetJson(scode string, report_type int, per_page int, page int) (*RespFinAnaJson, error) {
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
type FinChinaLiabilities struct {
}

func NewFinChinaLiabilities() *FinChinaLiabilities {
	return &FinChinaLiabilities{}
}

func (this *FinChinaLiabilities) getJson(scode string, report_type int, per_page int, page int) ([]LiabilitiesJson, error) {
	sli := make([]LiabilitiesJson, 0, per_page)
	ls, err := finchina.NewLiabilities().GetList(scode, report_type, per_page, page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		one := LiabilitiesJson{
			//资产
			AcRe: v.ACCORECE.Float64,      // 应收账款
			DFAs: v.DERIFINAASSET.Float64, // 衍生金融资产
			DiRe: v.DIVIDRECE.Float64,     // 应收股利
			DTAs: v.DEFETAXASSET.Float64,  // 递延所得税资产
			FAFS: v.AVAISELLASSE.Float64,  // 可供出售金融资产
			FAHT: v.TRADFINASSET.Float64,  // 交易性金融资产
			FiAs: v.FIXEDASSEIMMO.Float64, // 固定资产 表中只有 固定资产原值??
			GWil: v.GOODWILL.Float64,      // 商誉
			HTMI: v.HOLDINVEDUE.Float64,   // 持有至到期投资
			InAs: v.INTAASSET.Float64,     // 无形资产
			InRe: v.INTERECE.Float64,      // 应收利息
			LdLt: v.PLAC.Float64,          // 拆出资金
			LTAR: v.LONGRECE.Float64,      // 长期应收款
			LTEI: v.EQUIINVE.Float64,      // 长期股权投资
			LTPE: v.LOGPREPEXPE.Float64,   // 长期待摊费用
			MnFd: v.CURFDS.Float64,        // 货币资金
			NoRe: v.NOTESRECE.Float64,     // 应收票据
			OtRe: v.OTHERRECE.Float64,     // 其他应收款
			PrEx: v.PREPEXPE.Float64,      // 待摊费用
			Prpy: v.PREP.Float64,          // 预付款项
			REFI: v.INVEPROP.Float64,      // 投资性房地产
			ToAs: v.TOTASSET.Float64,      // 资产总计

			//负债
			AcEx:  v.ACCREXPE.Float64,         // 预提费用
			AcPy:  v.ACCOPAYA.Float64,         // 应付账款
			AdRE:  v.ADVAPAYM.Float64,         // 预收款项
			BdPy:  v.BDSPAYA.Float64,          // 应付债券
			CmPy:  v.COPEPOUN.Float64,         // 应付手续费及佣金
			DETLb: v.DEFEINCOTAXLIAB.Float64,  // 递延所得税负债
			DfIn:  v.DEFEREVE.Float64,         // 递延收益
			DFLb:  v.DERILIAB.Float64,         // 衍生金融负债
			DpFB:  v.DEPOSIT.Float64,          // 同业及其他金融机构存放款项 吸收存款及同业存放 ???
			DvPy:  v.DIVIPAYA.Float64,         // 应付股利
			FASR:  v.SELLREPASSE.Float64,      // 卖出回购金融资产款
			InPy:  v.INTEPAYA.Float64,         // 应付利息
			LnFB:  v.FDSBORR.Float64,          // 拆入资金
			LnFC:  v.CENBANKBORR.Float64,      // 向中央银行借款
			LTBw:  v.LONGBORR.Float64,         // 长期借款
			LTPy:  v.LONGPAYA.Float64,         // 长期应付款
			NCL1:  v.DUENONCLIAB.Float64,      // 一年内到期的非流动负债
			NtPy:  v.NOTESPAYA.Float64,        // 应付票据
			PCSc:  v.BDSPAYAPERBOND.Float64,   // 永续债
			PrSk:  v.BDSPAYAPREST.Float64,     // 优先股
			SaPy:  v.COPEWORKERSAL.Float64,    // 应付职工薪酬
			SBPy:  v.SHORTTERMBDSPAYA.Float64, // 应付短期债券
			STLn:  v.SHORTTERMBORR.Float64,    // 短期借款
			TaLb:  v.TOTLIAB.Float64,          // 负债合计
			TFLb:  v.TRADFINLIAB.Float64,      // 交易性金融负债
			TxPy:  v.TAXESPAYA.Float64,        // 应交税费

			//所有者权益
			BPCSET: v.PARESHARRIGH.Float64,    //  归属于母公司股东权益合计
			CDFCS:  v.CURTRANDIFF.Float64,     //  外币报表折算差额
			CpSp:   v.CAPISURP.Float64,        //  资本公积
			GRPr:   v.GENERISKRESE.Float64,    //  一般风险准备
			MiIt:   v.MINYSHARRIGH.Float64,    //  少数股东权益
			OEIn:   v.OTHEQUIN.Float64,        //  其他权益工具
			OESET:  v.RIGHAGGR.Float64,        //  所有者权益（或股东权益）合计
			OtCI:   v.OCL.Float64,             //  其他综合收益
			PCSe:   v.PERBOND.Float64,         //  永续债
			PICa:   v.PAIDINCAPI.Float64,      //  实收资本（或股本）
			PrSc:   v.PREST.Float64,           //  优先股
			SpRs:   v.RESE.Float64,            //  盈余公积
			TLSE:   v.TOTLIABSHAREQUI.Float64, //  负债和所有者权益（或股东权益）总计
			TrSc:   v.TREASTK.Float64,         //  库存股 表中名称(减：库存股)
			UdPr:   v.UNDIPROF.Float64,        //  未分配利润
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
