// TQ_FIN_PROCFSTATEMENTNEW	  一般企业现金流量表(新准则产品表)
package io_finchina

import (
	. "haina.com/share/models"

	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

//// TQ_FIN_PROCFSTATEMENTNEW	  一般企业现金流量表(新准则产品表)
//// __none__ 前缀的字段是参考其他证券软件的F10功能定义的Json返回字段信息,但在数据表中没有找到与之对应的字段,为不打乱与Wiki文档对应顺序而保留
//type TQ_FIN_PROCFSTATEMENTNEW struct {
//	Model `db:"-"`
//
//	ENDDATE    dbr.NullString //Date 	放置本次财报的截止日期
//	CUR        dbr.NullString //币种 	放置本次财报的币种
//	ACCSTACODE dbr.NullString //会计准则 代码 11001 中国会计准则(07年前版) 11002 中国会计准则(07年起版) 本表全为07年后的新会计准则
//
//	__none__CAIOA    dbr.NullFloat64 //CAIOA	经营活动现金流量净额调整项目
//	INICASHBALA      dbr.NullFloat64 //CEABg	期初现金及现金等价物 余额
//	FINALCASHBALA    dbr.NullFloat64 //CEAEd	期末现金及现金等价物 余额
//	PAYINTECASH      dbr.NullFloat64 //CFCms	支付手续费及佣金的现金
//	CHARINTECASH     dbr.NullFloat64 //CFIFC	收取利息、手续费及佣金的现金
//	RECEOTHERBIZCASH dbr.NullFloat64 //CFOTR	收到其他与经营活动有关的现金
//	BIZCASHINFL      dbr.NullFloat64 //CGBOA	经营活动现金流入小计
//	INVCASHINFL      dbr.NullFloat64 //CGIIA	投资活动现金流入小计
//	FINCASHINFL      dbr.NullFloat64 //CPBFA	筹资活动现金流入小计
//	ISSBDRECECASH    dbr.NullFloat64 //CPBID	发行债券收到的现金
//	RECEFROMLOAN     dbr.NullFloat64 //CPBLn	取得借款收到的现金
//	LABOPAYC         dbr.NullFloat64 //CPFGS	购买商品、接受劳务支付的现金
//	PAYWORKCASH      dbr.NullFloat64 //CPFSW	支付给职工以及为职工支付的现金
//	PAYTAX           dbr.NullFloat64 //CPFTx	支付的各项税费
//	INVRECECASH      dbr.NullFloat64 //CRBIv	吸收投资收到的现金
//	WITHINVGETCASH   dbr.NullFloat64 //CRFDI	收回投资收到的现金
//	LABORGETCASH     dbr.NullFloat64 //CRFGS	销售商品、提供劳务收到的现金
//	INVERETUGETCASH  dbr.NullFloat64 //CRFII	取得投资收益收到的现金
//	SUBSRECECASH     dbr.NullFloat64 //CRMSS	其中:子公司吸收少数股东投资收到的现金
//	INVCASHOUTF      dbr.NullFloat64 //CUIIA	投资活动现金流出小计
//	INVPAYC          dbr.NullFloat64 //CUIIv	投资所支付的现金
//	BIZCASHOUTF      dbr.NullFloat64 //CUIOA	经营活动现金流出小计
//	FIXEDASSETNETC   dbr.NullFloat64 //GDPES	处置固定资产、无形资产和其他长期资产收回的现金净额
//	CHGEXCHGCHGS     dbr.NullFloat64 //IERCE	汇率变动对现金及现金等价物的影响
//	FININSTNETR      dbr.NullFloat64 //NBFBI	向其他金融机构拆入资金净增加额
//	BANKLOANNETINCR  dbr.NullFloat64 //NBFCB	向中央银行借款净增加额
//	__none__NCEAI    dbr.NullFloat64 //NCEAI	现金及现金等价物净增加额的调整项目
//	__none__NCEIS    dbr.NullFloat64 //NCEIS	现金及现金等价物净增加额的特殊项目
//	MANANETR         dbr.NullFloat64 //NCFOA	经营活动产生的现金流量净额
//	INVNETCASHFLOW   dbr.NullFloat64 //NCIIA	投资活动产生的现金流量净额
//	FINNETCFLOW      dbr.NullFloat64 //NCPFA	筹资活动产生的现金流量净额
//	SUBSNETC         dbr.NullFloat64 //NCRDU	处置子公司及其他营业单位收到的现金净额
//	SUBSPAYNETCASH   dbr.NullFloat64 //NCRFU	取得子公司及其他营业单位支付的现金净额
//	TRADEPAYMNETR    dbr.NullFloat64 //NDCBI	存放中央银行和同业款项净增加额
//	LOANNETR         dbr.NullFloat64 //NIcLn	质押贷款净增加额
//	CASHNETR         dbr.NullFloat64 //NIICE	现金及现金等价物净增加额
//	__none__NLend    dbr.NullFloat64 //NLend	拆出资金净增加额
//	LOANSNETR        dbr.NullFloat64 //NLnAv	客户贷款及垫款净增加额
//	ACQUASSETCASH    dbr.NullFloat64 //PcsPE	购建固定资产、无形资产和其他长期资产支付的现金
//	DEBTPAYCASH      dbr.NullFloat64 //PmFPy	偿还债务支付的现金
//	DIVIPROFPAYCASH  dbr.NullFloat64 //PmISA	分配股利、利润或偿付利息支付的现金
//	FINCASHOUTF      dbr.NullFloat64 //PmoFA	筹资活动现金流出小计
//}

type TQ_FIN_PROCFSTATEMENTNEW struct {
	Model `db:"-"`

	ENDDATE        dbr.NullString  //Date 	放置本次财报的截止日期
	MANANETR       dbr.NullFloat64 //NCFOA	经营活动产生的现金流量净额
	INVNETCASHFLOW dbr.NullFloat64 //NCIIA	投资活动产生的现金流量净额
	FINNETCFLOW    dbr.NullFloat64 //NCPFA	筹资活动产生的现金流量净额
}

func (this *TQ_FIN_PROCFSTATEMENTNEW) Write() error {
	panic("implement me")
}

func NewTQ_FIN_PROCFSTATEMENTNEW() *TQ_FIN_PROCFSTATEMENTNEW {
	return &TQ_FIN_PROCFSTATEMENTNEW{
		Model: Model{
			TableName: finchina.TABLE_TQ_FIN_PROCFSTATEMENTNEW,
			Db:        MyCat,
		},
	}
}

func (this *TQ_FIN_PROCFSTATEMENTNEW) getListByCompcode(compcode string, report_data_type int, per_page int, page int) ([]TQ_FIN_PROCFSTATEMENTNEW, error) {
	var sli []TQ_FIN_PROCFSTATEMENTNEW

	builder := this.Db.Select("ENDDATE", "MANANETR", "INVNETCASHFLOW", "FINNETCFLOW").From(this.TableName)
	if report_data_type != 0 {
		builder.Where("REPORTDATETYPE=?", report_data_type)
	}
	err := builder.Where("COMPCODE=?", compcode).
		Where("ISVALID=1").
		Where("REPORTTYPE=?", 3).
		OrderBy("ENDDATE DESC").
		Paginate(uint64(page), uint64(per_page)).
		LoadStruct(&sli)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}

	return sli, nil
}

//------------------------------------------------------------------------------

func (this *TQ_FIN_PROCFSTATEMENTNEW) GetList(scode string, market string, report_data_type int, per_page int, page int) ([]TQ_FIN_PROCFSTATEMENTNEW, error) {

	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		logging.Error("%T GetList error: %s", *this, err)
		return nil, err
	}

	return this.getListByCompcode(sc.COMPCODE.String, report_data_type, per_page, page)
}
