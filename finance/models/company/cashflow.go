// 现金流量表
package company

import (
	"time"

	"haina.com/market/finance/models/finchina"
)

type CashflowJson struct {
	Date  int64   `json:"Date"`  //截止日期 unix时间戳
	CAIOA float64 `json:"CAIOA"` //经营活动现金流量净额调整项目
	CEABg float64 `json:"CEABg"` //期初现金及现金等价物 余额
	CEAEd float64 `json:"CEAEd"` //期末现金及现金等价物 余额
	CFCms float64 `json:"CFCms"` //支付手续费及佣金的现金
	CFIFC float64 `json:"CFIFC"` //收取利息、手续费及佣金的现金
	CFOTR float64 `json:"CFOTR"` //收到其他与经营活动有关的现金
	CGBOA float64 `json:"CGBOA"` //经营活动现金流入小计
	CGIIA float64 `json:"CGIIA"` //投资活动现金流入小计
	CPBFA float64 `json:"CPBFA"` //筹资活动现金流入小计
	CPBID float64 `json:"CPBID"` //发行债券收到的现金
	CPBLn float64 `json:"CPBLn"` //取得借款收到的现金
	CPFGS float64 `json:"CPFGS"` //购买商品、接受劳务支付的现金
	CPFSW float64 `json:"CPFSW"` //支付给职工以及为职工支付的现金
	CPFTx float64 `json:"CPFTx"` //支付的各项税费
	CRBIv float64 `json:"CRBIv"` //吸收投资收到的现金
	CRFDI float64 `json:"CRFDI"` //收回投资收到的现金
	CRFGS float64 `json:"CRFGS"` //销售商品、提供劳务收到的现金
	CRFII float64 `json:"CRFII"` //取得投资收益收到的现金
	CRMSS float64 `json:"CRMSS"` //其中:子公司吸收少数股东投资收到的现金
	CUIIA float64 `json:"CUIIA"` //投资活动现金流出小计
	CUIIv float64 `json:"CUIIv"` //投资支付的现金
	CUIOA float64 `json:"CUIOA"` //经营活动现金流出小计
	GDPES float64 `json:"GDPES"` //处置固定资产、无形资产和其他长期资产收回的现金净额
	IERCE float64 `json:"IERCE"` //汇率变动对现金及现金等价物的影响
	NBFBI float64 `json:"NBFBI"` //向其他金融机构拆入资金净增加额
	NBFCB float64 `json:"NBFCB"` //向中央银行借款净增加额
	NCEAI float64 `json:"NCEAI"` //现金及现金等价物净增加额的调整项目
	NCEIS float64 `json:"NCEIS"` //现金及现金等价物净增加额的特殊项目
	NCFOA float64 `json:"NCFOA"` //经营活动产生的现金流量净额
	NCIIA float64 `json:"NCIIA"` //投资活动产生的现金流量净额
	NCPFA float64 `json:"NCPFA"` //筹资活动产生的现金流量净额
	NCRDU float64 `json:"NCRDU"` //处置子公司及其他营业单位收到的现金净额
	NCRFU float64 `json:"NCRFU"` //取得子公司及其他营业单位支付的现金净额
	NDCBI float64 `json:"NDCBI"` //存放中央银行和同业款项净增加额
	NIcLn float64 `json:"NIcLn"` //质押贷款净增加额
	NIICE float64 `json:"NIICE"` //现金及现金等价物 净增加额
	NLend float64 `json:"NLend"` //拆出资金净增加额
	NLnAv float64 `json:"NLnAv"` //客户贷款及垫款净增加额
	PcsPE float64 `json:"PcsPE"` //购建固定资产、无形资产和其他长期资产支付的现金
	PmFPy float64 `json:"PmFPy"` //偿还债务支付的现金
	PmISA float64 `json:"PmISA"` //分配股利、利润或偿付利息支付的现金
	PmoFA float64 `json:"PmoFA"` //筹资活动现金流出小计
}

func NewCashflowJson() *CashflowJson {
	return &CashflowJson{}
}

//------------------------------------------------------------------------------
type Cashflow struct {
}

func NewCashflow() *Cashflow {
	return &Cashflow{}
}

func (this *Cashflow) getList(scode string, report_type int, per_page int, page int) ([]Cashflow, error) {
	return nil, nil
}
func (this *Cashflow) GetList(scode string, report_type int, per_page int, page int) ([]Cashflow, error) {
	return nil, nil
}
func (this *Cashflow) getJson(scode string, report_type int, per_page int, page int) ([]CashflowJson, error) {
	return NewFinChinaCashflow().getJson(scode, report_type, per_page, page)
}
func (this *Cashflow) GetJson(scode string, report_type int, per_page int, page int) (*RespFinAnaJson, error) {
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
type FinChinaCashflow struct {
}

func NewFinChinaCashflow() *FinChinaCashflow {
	return &FinChinaCashflow{}
}

func (this *FinChinaCashflow) getJson(scode string, report_type int, per_page int, page int) ([]CashflowJson, error) {
	sli := make([]CashflowJson, 0, per_page)
	ls, err := finchina.NewCashflow().GetList(scode, report_type, per_page, page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		one := CashflowJson{
			CEABg: v.INICASHBALA.Float64,      //期初现金及现金等价物 余额
			CEAEd: v.FINALCASHBALA.Float64,    //期末现金及现金等价物 余额
			CFCms: v.PAYINTECASH.Float64,      //支付手续费及佣金的现金
			CFIFC: v.CHARINTECASH.Float64,     //收取利息、手续费及佣金的现金
			CFOTR: v.RECEOTHERBIZCASH.Float64, //收到其他与经营活动有关的现金
			CGBOA: v.BIZCASHINFL.Float64,      //经营活动现金流入小计
			CGIIA: v.INVCASHINFL.Float64,      //投资活动现金流入小计
			CPBFA: v.FINCASHINFL.Float64,      //筹资活动现金流入小计
			CPBID: v.ISSBDRECECASH.Float64,    //发行债券收到的现金
			CPBLn: v.RECEFROMLOAN.Float64,     //取得借款收到的现金
			CPFGS: v.LABOPAYC.Float64,         //购买商品、接受劳务支付的现金
			CPFSW: v.PAYWORKCASH.Float64,      //支付给职工以及为职工支付的现金
			CPFTx: v.PAYTAX.Float64,           //支付的各项税费
			CRBIv: v.INVRECECASH.Float64,      //吸收投资收到的现金
			CRFDI: v.WITHINVGETCASH.Float64,   //收回投资收到的现金
			CRFGS: v.LABORGETCASH.Float64,     //销售商品、提供劳务收到的现金
			CRFII: v.INVERETUGETCASH.Float64,  //取得投资收益收到的现金
			CRMSS: v.SUBSRECECASH.Float64,     //其中:子公司吸收少数股东投资收到的现金
			CUIIA: v.INVCASHOUTF.Float64,      //投资活动现金流出小计
			CUIIv: v.INVPAYC.Float64,          //投资所支付的现金
			CUIOA: v.BIZCASHOUTF.Float64,      //经营活动现金流出小计
			GDPES: v.FIXEDASSETNETC.Float64,   //处置固定资产、无形资产和其他长期资产收回的现金净额
			IERCE: v.CHGEXCHGCHGS.Float64,     //汇率变动对现金及现金等价物的影响
			NBFBI: v.FININSTNETR.Float64,      //向其他金融机构拆入资金净增加额
			NBFCB: v.BANKLOANNETINCR.Float64,  //向中央银行借款净增加额
			NCFOA: v.MANANETR.Float64,         //经营活动产生的现金流量净额
			NCIIA: v.INVNETCASHFLOW.Float64,   //投资活动产生的现金流量净额
			NCPFA: v.FINNETCFLOW.Float64,      //筹资活动产生的现金流量净额
			NCRDU: v.SUBSNETC.Float64,         //处置子公司及其他营业单位收到的现金净额
			NCRFU: v.SUBSPAYNETCASH.Float64,   //取得子公司及其他营业单位支付的现金净额
			NDCBI: v.TRADEPAYMNETR.Float64,    //存放中央银行和同业款项净增加额
			NIcLn: v.LOANNETR.Float64,         //质押贷款净增加额
			NIICE: v.CASHNETR.Float64,         //现金及现金等价物净增加额
			NLnAv: v.LOANSNETR.Float64,        //客户贷款及垫款净增加额
			PcsPE: v.ACQUASSETCASH.Float64,    //购建固定资产、无形资产和其他长期资产支付的现金
			PmFPy: v.DEBTPAYCASH.Float64,      //偿还债务支付的现金
			PmISA: v.DIVIPROFPAYCASH.Float64,  //分配股利、利润或偿付利息支付的现金
			PmoFA: v.FINCASHOUTF.Float64,      //筹资活动现金流出小计
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
