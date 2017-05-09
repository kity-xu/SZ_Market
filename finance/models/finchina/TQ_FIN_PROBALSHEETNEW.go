// 资产负债表
// TQ_FIN_PROBALSHEETNEW	  一般企业资产负债表(新准则产品表)
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

//  TQ_FIN_PROBALSHEETNEW	  一般企业资产负债表(新准则产品表)
// __none__ 前缀的字段是参考其他证券软件的F10功能定义的Json返回字段信息,但在数据表中没有找到与之对应的字段,为不打乱与Wiki文档对应顺序而保留
type TQ_FIN_PROBALSHEETNEW struct {
	Model `db:"-"`

	ENDDATE    dbr.NullString //Date 	放置本次财报的截止日期
	CUR        dbr.NullString //币种 	放置本次财报的币种
	ACCSTACODE dbr.NullString //会计准则 代码 11001 中国会计准则(07年前版) 11002 中国会计准则(07年起版) 本表全为07年后的新会计准则

	//资产
	ACCORECE      dbr.NullFloat64 //AcRe  应收账款
	__none__CuDe  dbr.NullFloat64 //CuDe  客户资金存款 TQ_FIN_PROBBALBSHEETNEW 银行资产负债表 CLIEDEPO 客户存款<吸收存款>
	__none__DeMg  dbr.NullFloat64 //DeMg  存出保证金 TQ_FIN_PROBBALBSHEETNEW 银行资产负债表 KEPTMARG
	DERIFINAASSET dbr.NullFloat64 //DFAs  衍生金融资产
	__none__DfSv  dbr.NullFloat64 //DfSv  定期存款
	DIVIDRECE     dbr.NullFloat64 //DiRe  应收股利
	DEFETAXASSET  dbr.NullFloat64 //DTAs  递延所得税资产
	AVAISELLASSE  dbr.NullFloat64 //FAFS  可供出售金融资产
	TRADFINASSET  dbr.NullFloat64 //FAHT  交易性金融资产
	FIXEDASSEIMMO dbr.NullFloat64 //FiAs  固定资产 表中只有 固定资产原值??
	GOODWILL      dbr.NullFloat64 //GWil  商誉
	HOLDINVEDUE   dbr.NullFloat64 //HTMI  持有至到期投资
	__none__IbDe  dbr.NullFloat64 //IbDe  存放同业款项 TQ_FIN_PROBBALBSHEETNEW 银行资产负债表 DEPOOTHERBANK
	INTAASSET     dbr.NullFloat64 //InAs  无形资产
	INTERECE      dbr.NullFloat64 //InRe  应收利息
	PLAC          dbr.NullFloat64 //LdLt  拆出资金
	LONGRECE      dbr.NullFloat64 //LTAR  长期应收款
	EQUIINVE      dbr.NullFloat64 //LTEI  长期股权投资
	LOGPREPEXPE   dbr.NullFloat64 //LTPE  长期待摊费用
	__none__Metl  dbr.NullFloat64 //Metl  贵金属 TQ_FIN_PROBBALBSHEETNEW 银行资产负债表 EXPEMETA
	CURFDS        dbr.NullFloat64 //MnFd  货币资金
	NOTESRECE     dbr.NullFloat64 //NoRe  应收票据
	OTHERRECE     dbr.NullFloat64 //OtRe  其他应收款
	PREPEXPE      dbr.NullFloat64 //PrEx  待摊费用
	PREP          dbr.NullFloat64 //Prpy  预付款项
	INVEPROP      dbr.NullFloat64 //REFI  投资性房地产
	TOTASSET      dbr.NullFloat64 //ToAs  资产总计

	//负债
	ACCREXPE         dbr.NullFloat64 //AcEx   预提费用
	ACCOPAYA         dbr.NullFloat64 //AcPy   应付账款
	ADVAPAYM         dbr.NullFloat64 //AdRE   预收款项
	BDSPAYA          dbr.NullFloat64 //BdPy   应付债券
	COPEPOUN         dbr.NullFloat64 //CmPy   应付手续费及佣金
	DEFEINCOTAXLIAB  dbr.NullFloat64 //DETLb  递延所得税负债
	DEFEREVE         dbr.NullFloat64 //DfIn   递延收益 此表中有(一年内的递延收益 和 长期递延收益)两项,此处用 一年内的递延收益???; TQ_FIN_PROBBALBSHEETNEW 银行资产负债表中 DEFEINCO 递延收益
	DERILIAB         dbr.NullFloat64 //DFLb   衍生金融负债
	__none__DpCl     dbr.NullFloat64 //DpCl   吸收存款 TQ_FIN_PROBBALBSHEETNEW 银行资产负债表 CLIEDEPO 客户存款<吸收存款>
	DEPOSIT          dbr.NullFloat64 //DpFB   同业及其他金融机构存放款项 吸收存款及同业存放 ???
	DIVIPAYA         dbr.NullFloat64 //DvPy   应付股利
	SELLREPASSE      dbr.NullFloat64 //FASR   卖出回购金融资产款
	INTEPAYA         dbr.NullFloat64 //InPy   应付利息
	FDSBORR          dbr.NullFloat64 //LnFB   拆入资金
	CENBANKBORR      dbr.NullFloat64 //LnFC   向中央银行借款
	LONGBORR         dbr.NullFloat64 //LTBw   长期借款
	LONGPAYA         dbr.NullFloat64 //LTPy   长期应付款
	DUENONCLIAB      dbr.NullFloat64 //NCL1   一年内到期的非流动负债
	NOTESPAYA        dbr.NullFloat64 //NtPy   应付票据
	BDSPAYAPERBOND   dbr.NullFloat64 //PCSc   永续债
	__none__PlLn     dbr.NullFloat64 //PlLn   质押借款
	BDSPAYAPREST     dbr.NullFloat64 //PrSk   优先股
	COPEWORKERSAL    dbr.NullFloat64 //SaPy   应付职工薪酬
	SHORTTERMBDSPAYA dbr.NullFloat64 //SBPy   应付短期债券
	SHORTTERMBORR    dbr.NullFloat64 //STLn   短期借款
	TOTLIAB          dbr.NullFloat64 //TaLb   负债合计
	TRADFINLIAB      dbr.NullFloat64 //TFLb   交易性金融负债
	TAXESPAYA        dbr.NullFloat64 //TxPy   应交税费

	//所有者权益
	__none__BPCOEAI dbr.NullFloat64 //BPCOEAI  归属于母公司所有者权益调整项目
	__none__BPCOESI dbr.NullFloat64 //BPCOESI  归属于母公司所有者权益特殊项目
	PARESHARRIGH    dbr.NullFloat64 //BPCSET   归属于母公司股东权益合计
	CURTRANDIFF     dbr.NullFloat64 //CDFCS    外币报表折算差额
	CAPISURP        dbr.NullFloat64 //CpSp     资本公积
	GENERISKRESE    dbr.NullFloat64 //GRPr     一般风险准备
	__none__LEAI    dbr.NullFloat64 //LEAI     负债和权益调整项目
	__none__LESI    dbr.NullFloat64 //LESI     负债和权益特殊项目
	MINYSHARRIGH    dbr.NullFloat64 //MiIt     少数股东权益
	__none__OEAI    dbr.NullFloat64 //OEAI     所有者权益调整项目
	OTHEQUIN        dbr.NullFloat64 //OEIn     其他权益工具
	RIGHAGGR        dbr.NullFloat64 //OESET    所有者权益（或股东权益）合计
	OCL             dbr.NullFloat64 //OtCI     其他综合收益
	PERBOND         dbr.NullFloat64 //PCSe     永续债
	PAIDINCAPI      dbr.NullFloat64 //PICa     实收资本（或股本）
	PREST           dbr.NullFloat64 //PrSc     优先股
	RESE            dbr.NullFloat64 //SpRs     盈余公积
	TOTLIABSHAREQUI dbr.NullFloat64 //TLSE     负债和所有者权益（或股东权益）总计
	TREASTK         dbr.NullFloat64 //TrSc     库存股 表中名称(减：库存股)
	UNDIPROF        dbr.NullFloat64 //UdPr     未分配利润
}

func NewTQ_FIN_PROBALSHEETNEW() *TQ_FIN_PROBALSHEETNEW {
	return &TQ_FIN_PROBALSHEETNEW{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROBALSHEETNEW,
			Db:        MyCat,
		},
	}
}

func (this *TQ_FIN_PROBALSHEETNEW) getListByCompcode(compcode string, report_type int, per_page int, page int) ([]TQ_FIN_PROBALSHEETNEW, error) {
	var sli []TQ_FIN_PROBALSHEETNEW

	builder := this.Db.Select("*").From(this.TableName)
	if report_type != 0 {
		builder.Where("REPORTDATETYPE=?", report_type)
	}
	err := builder.Where("COMPCODE=?", compcode).
		Where("ISVALID=1").
		Where("REPORTTYPE=?", 1).
		OrderBy("ENDDATE DESC").
		Paginate(uint64(page), uint64(per_page)).
		LoadStruct(&sli)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}

	return sli, nil
}

//------------------------------------------------------------------------------

func (this *TQ_FIN_PROBALSHEETNEW) GetList(scode string, market string, report_type int, per_page int, page int) ([]TQ_FIN_PROBALSHEETNEW, error) {

	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		logging.Error("%T GetList error: %s", *this, err)
		return nil, err
	}

	return this.getListByCompcode(sc.COMPCODE.String, report_type, per_page, page)
}

//--------------------------------------------------------------------------------
//  TQ_FIN_PROBBALBSHEETNEW	  银行资产负债表(新准则产品表)
//  TQ_FIN_PROIBALSHEETNEW	  保险资产负债表(新准则产品表)
//  TQ_FIN_PROSBALSHEETNEW	  证券资产负债表(新准则产品表)
