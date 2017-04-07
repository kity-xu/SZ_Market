// 现金流量表
package company

import (
	"time"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Cashflow struct {
}

func NewCashflow() *Cashflow {
	return &Cashflow{}
}

func (this *Cashflow) GET(c *gin.Context) {
	scode := c.Query(models.CONTEXT_SCODE)
	stype := c.Query(models.CONTEXT_TYPE)
	spage := c.Query(models.CONTEXT_PAGE)
	sperp := c.Query(models.CONTEXT_PERPAGE)

	req := CheckAndNewRequestParam(scode, stype, sperp, spage)
	if req == nil {
		lib.WriteString(c, 40004, nil)
		return
	}

	data, err := this.getJson(req)
	if err != nil {
		logging.Debug("%v", err)
		lib.WriteString(c, 40002, nil)
		return
	}
	data.SCode = scode

	lib.WriteString(c, 200, data)
}

func (this *Cashflow) getJson(req *RequestParam) (*company.RespFinAnaJson, error) {
	sli := make([]company.CashflowJson, 0, req.PerPage)
	ls, err := company.NewCashflow().GetList(req.SCode, req.Type, req.PerPage, req.Page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		one := company.CashflowJson{
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

	// A股
	jsn := &company.RespFinAnaJson{
		MU:     "人民币元",
		AS:     "新会计准则",
		Length: len(sli),
		List:   sli,
	}
	return jsn, nil
}
