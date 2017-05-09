// 现金流量表
package company

import (
	"haina.com/market/finance/models/finchina"
)

type Cashflow struct {
	FinChinaCashflow
}

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

//------------------------------------------------------------------------------
func NewCashflow() *Cashflow {
	return &Cashflow{}
}

func (this *Cashflow) GetList(scode string, market string, report_type int, per_page int, page int) ([]Cashflow, error) {
	return NewFinChinaCashflow().getCashflowList(scode, market, report_type, per_page, page)
}

//------------------------------------------------------------------------------
type FinChinaCashflow struct {
	finchina.TQ_FIN_PROCFSTATEMENTNEW
}

func NewFinChinaCashflow() *FinChinaCashflow {
	return &FinChinaCashflow{}
}

func (this *FinChinaCashflow) getCashflowList(scode string, market string, report_type int, per_page int, page int) ([]Cashflow, error) {
	var (
		slidb []finchina.TQ_FIN_PROCFSTATEMENTNEW
		len1  int
		err   error
	)
	sli := make([]Cashflow, 0, per_page)

	slidb, err = finchina.NewTQ_FIN_PROCFSTATEMENTNEW().GetList(scode, market, report_type, per_page, page)
	if err != nil {
		return nil, err
	}
	if len1 = len(slidb); 0 == len1 {
		return sli, nil
	}

	for _, v := range slidb {
		one := Cashflow{
			FinChinaCashflow: FinChinaCashflow{
				TQ_FIN_PROCFSTATEMENTNEW: v,
			},
		}
		sli = append(sli, one)
	}

	return sli, nil
}
