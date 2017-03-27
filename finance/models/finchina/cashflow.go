// 现金流量数据
package finchina

import (
	"time"

	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
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

// TQ_FIN_PROCFSTATEMENTNEW	  一般企业现金流量表(新准则产品表)
// __none__ 前缀的字段是参考其他证券软件的F10功能定义的Json返回字段信息但在数据表中没有找到与之对应的字段,为了不打乱顺序,做个标注
type CashflowGeneral struct {
	Model `db:"-"`

	ENDDATE    dbr.NullString //Date 	放置本次财报的截止日期
	CUR        dbr.NullString //币种 	放置本次财报的币种
	ACCSTACODE dbr.NullString //会计准则 代码 11001 中国会计准则(07年前版) 11002 中国会计准则(07年起版) 本表全为07年后的新会计准则

	__none__CAIOA    dbr.NullFloat64 //CAIOA	经营活动现金流量净额调整项目
	INICASHBALA      dbr.NullFloat64 //CEABg	期初现金及现金等价物 余额
	FINALCASHBALA    dbr.NullFloat64 //CEAEd	期末现金及现金等价物 余额
	PAYINTECASH      dbr.NullFloat64 //CFCms	支付手续费及佣金的现金
	CHARINTECASH     dbr.NullFloat64 //CFIFC	收取利息、手续费及佣金的现金
	RECEOTHERBIZCASH dbr.NullFloat64 //CFOTR	收到其他与经营活动有关的现金
	BIZCASHINFL      dbr.NullFloat64 //CGBOA	经营活动现金流入小计
	INVCASHINFL      dbr.NullFloat64 //CGIIA	投资活动现金流入小计
	FINCASHINFL      dbr.NullFloat64 //CPBFA	筹资活动现金流入小计
	ISSBDRECECASH    dbr.NullFloat64 //CPBID	发行债券收到的现金
	RECEFROMLOAN     dbr.NullFloat64 //CPBLn	取得借款收到的现金
	LABOPAYC         dbr.NullFloat64 //CPFGS	购买商品、接受劳务支付的现金
	PAYWORKCASH      dbr.NullFloat64 //CPFSW	支付给职工以及为职工支付的现金
	PAYTAX           dbr.NullFloat64 //CPFTx	支付的各项税费
	INVRECECASH      dbr.NullFloat64 //CRBIv	吸收投资收到的现金
	WITHINVGETCASH   dbr.NullFloat64 //CRFDI	收回投资收到的现金
	LABORGETCASH     dbr.NullFloat64 //CRFGS	销售商品、提供劳务收到的现金
	INVERETUGETCASH  dbr.NullFloat64 //CRFII	取得投资收益收到的现金
	SUBSRECECASH     dbr.NullFloat64 //CRMSS	其中:子公司吸收少数股东投资收到的现金
	INVCASHOUTF      dbr.NullFloat64 //CUIIA	投资活动现金流出小计
	INVPAYC          dbr.NullFloat64 //CUIIv	投资所支付的现金
	BIZCASHOUTF      dbr.NullFloat64 //CUIOA	经营活动现金流出小计
	FIXEDASSETNETC   dbr.NullFloat64 //GDPES	处置固定资产、无形资产和其他长期资产收回的现金净额
	CHGEXCHGCHGS     dbr.NullFloat64 //IERCE	汇率变动对现金及现金等价物的影响
	FININSTNETR      dbr.NullFloat64 //NBFBI	向其他金融机构拆入资金净增加额
	BANKLOANNETINCR  dbr.NullFloat64 //NBFCB	向中央银行借款净增加额
	__none__NCEAI    dbr.NullFloat64 //NCEAI	现金及现金等价物净增加额的调整项目
	__none__NCEIS    dbr.NullFloat64 //NCEIS	现金及现金等价物净增加额的特殊项目
	MANANETR         dbr.NullFloat64 //NCFOA	经营活动产生的现金流量净额
	INVNETCASHFLOW   dbr.NullFloat64 //NCIIA	投资活动产生的现金流量净额
	FINNETCFLOW      dbr.NullFloat64 //NCPFA	筹资活动产生的现金流量净额
	SUBSNETC         dbr.NullFloat64 //NCRDU	处置子公司及其他营业单位收到的现金净额
	SUBSPAYNETCASH   dbr.NullFloat64 //NCRFU	取得子公司及其他营业单位支付的现金净额
	TRADEPAYMNETR    dbr.NullFloat64 //NDCBI	存放中央银行和同业款项净增加额
	LOANNETR         dbr.NullFloat64 //NIcLn	质押贷款净增加额
	CASHNETR         dbr.NullFloat64 //NIICE	现金及现金等价物净增加额
	__none__NLend    dbr.NullFloat64 //NLend	拆出资金净增加额
	LOANSNETR        dbr.NullFloat64 //NLnAv	客户贷款及垫款净增加额
	ACQUASSETCASH    dbr.NullFloat64 //PcsPE	购建固定资产、无形资产和其他长期资产支付的现金
	DEBTPAYCASH      dbr.NullFloat64 //PmFPy	偿还债务支付的现金
	DIVIPROFPAYCASH  dbr.NullFloat64 //PmISA	分配股利、利润或偿付利息支付的现金
	FINCASHOUTF      dbr.NullFloat64 //PmoFA	筹资活动现金流出小计
}

func NewCashflowGeneral() *CashflowGeneral {
	return &CashflowGeneral{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *CashflowGeneral) getCashflowJsonList(compcode string, req *RequestParam) ([]*CashflowJson, error) {
	logging.Info("getCashflowJsonList %T, compcode %s", *this, compcode)
	var (
		sli_db []CashflowGeneral
	)
	sli := make([]*CashflowJson, 0)

	builder := this.Db.Select("*").From(TABLE_TQ_FIN_PROCFSTATEMENTNEW)
	if req.Type != 0 {
		builder.Where("REPORTDATETYPE=?", req.Type)
	}
	err := builder.Where("COMPCODE = ?", compcode).
		Where("REPORTTYPE = ?", 1).
		OrderBy("ENDDATE DESC").
		Paginate(uint64(req.Page), uint64(req.PerPage)).
		LoadStruct(&sli_db)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}

	for _, v := range sli_db {
		one := NewCashflowJson()

		one.CEABg = v.INICASHBALA.Float64      //期初现金及现金等价物 余额
		one.CEAEd = v.FINALCASHBALA.Float64    //期末现金及现金等价物 余额
		one.CFCms = v.PAYINTECASH.Float64      //支付手续费及佣金的现金
		one.CFIFC = v.CHARINTECASH.Float64     //收取利息、手续费及佣金的现金
		one.CFOTR = v.RECEOTHERBIZCASH.Float64 //收到其他与经营活动有关的现金
		one.CGBOA = v.BIZCASHINFL.Float64      //经营活动现金流入小计
		one.CGIIA = v.INVCASHINFL.Float64      //投资活动现金流入小计
		one.CPBFA = v.FINCASHINFL.Float64      //筹资活动现金流入小计
		one.CPBID = v.ISSBDRECECASH.Float64    //发行债券收到的现金
		one.CPBLn = v.RECEFROMLOAN.Float64     //取得借款收到的现金
		one.CPFGS = v.LABOPAYC.Float64         //购买商品、接受劳务支付的现金
		one.CPFSW = v.PAYWORKCASH.Float64      //支付给职工以及为职工支付的现金
		one.CPFTx = v.PAYTAX.Float64           //支付的各项税费
		one.CRBIv = v.INVRECECASH.Float64      //吸收投资收到的现金
		one.CRFDI = v.WITHINVGETCASH.Float64   //收回投资收到的现金
		one.CRFGS = v.LABORGETCASH.Float64     //销售商品、提供劳务收到的现金
		one.CRFII = v.INVERETUGETCASH.Float64  //取得投资收益收到的现金
		one.CRMSS = v.SUBSRECECASH.Float64     //其中=子公司吸收少数股东投资收到的现金
		one.CUIIA = v.INVCASHOUTF.Float64      //投资活动现金流出小计
		one.CUIIv = v.INVPAYC.Float64          //投资所支付的现金
		one.CUIOA = v.BIZCASHOUTF.Float64      //经营活动现金流出小计
		one.GDPES = v.FIXEDASSETNETC.Float64   //处置固定资产、无形资产和其他长期资产收回的现金净额
		one.IERCE = v.CHGEXCHGCHGS.Float64     //汇率变动对现金及现金等价物的影响
		one.NBFBI = v.FININSTNETR.Float64      //向其他金融机构拆入资金净增加额
		one.NBFCB = v.BANKLOANNETINCR.Float64  //向中央银行借款净增加额
		one.NCFOA = v.MANANETR.Float64         //经营活动产生的现金流量净额
		one.NCIIA = v.INVNETCASHFLOW.Float64   //投资活动产生的现金流量净额
		one.NCPFA = v.FINNETCFLOW.Float64      //筹资活动产生的现金流量净额
		one.NCRDU = v.SUBSNETC.Float64         //处置子公司及其他营业单位收到的现金净额
		one.NCRFU = v.SUBSPAYNETCASH.Float64   //取得子公司及其他营业单位支付的现金净额
		one.NDCBI = v.TRADEPAYMNETR.Float64    //存放中央银行和同业款项净增加额
		one.NIcLn = v.LOANNETR.Float64         //质押贷款净增加额
		one.NIICE = v.CASHNETR.Float64         //现金及现金等价物净增加额
		one.NLnAv = v.LOANSNETR.Float64        //客户贷款及垫款净增加额
		one.PcsPE = v.ACQUASSETCASH.Float64    //购建固定资产、无形资产和其他长期资产支付的现金
		one.PmFPy = v.DEBTPAYCASH.Float64      //偿还债务支付的现金
		one.PmISA = v.DIVIPROFPAYCASH.Float64  //分配股利、利润或偿付利息支付的现金
		one.PmoFA = v.FINCASHOUTF.Float64      //筹资活动现金流出小计

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

//------------------------------------------------------------------------------

type CashflowInfo struct {
}

func NewCashflowInfo() *CashflowInfo {
	return &CashflowInfo{}
}

func (this *CashflowInfo) GetJson(req *RequestParam) (*ResponseFinAnaJson, error) {
	logging.Info("GetJson %T, RequestParam %+v", *this, *req)

	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(req.SCode); err != nil {
		return nil, err
	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, req.SCode)
		return nil, ErrNullComp
	}

	sli := NewCashflowGeneral()
	list, err := sli.getCashflowJsonList(sc.COMPCODE.String, req)
	if err != nil {
		return nil, err
	}

	res := &ResponseFinAnaJson{
		SCode: req.SCodeOrigin,
		MU:    "人民币元",
		AS:    "新会计准则",
	}

	res.List = list
	res.Length = len(list)
	return res, nil
}

//--------------------------------------------------------------------------------

//  TQ_FIN_PROBCFSTATEMENTNEW	  银行现金流量表(新准则产品表)
//  TQ_FIN_PROICFSTATEMENTNEW	  保险现金流量表(新准则产品表)
//  TQ_FIN_PROSCFSTATEMENTNEW	  证券现金流量表(新准则产品表)
