// 利润表
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

// TQ_FIN_PROINCSTATEMENTNEW    中文名称：一般企业利润表(新准则产品表)
// __none__ 前缀的字段是参考其他证券软件的F10功能定义的Json返回字段信息,但在数据表中没有找到与之对应的字段,为不打乱与Wiki文档对应顺序而保留
type TQ_FIN_PROINCSTATEMENTNEW struct {
	Model `db:"-"`

	ENDDATE    dbr.NullString //Date 	放置本次财报的截止日期
	CUR        dbr.NullString //币种 	放置本次财报的币种
	ACCSTACODE dbr.NullString //会计准则 代码 11001 中国会计准则(07年前版) 11002 中国会计准则(07年起版) 本表全为07年后的新会计准则

	__none__AAPC dbr.NullFloat64 //AAPC		影响母公司净利润的调整项目
	ASSEIMPALOSS dbr.NullFloat64 //AILs		资产减值损失
	REINEXPE     dbr.NullFloat64 //AREp		分保费用
	__none__BAEp dbr.NullFloat64 //BAEp		业务及管理费用 银行,保险利润表 存在该字段
	PARENETP     dbr.NullFloat64 //BPAC		归属于母公司所有者的净利润
	POUNEXPE     dbr.NullFloat64 //CoEp		手续费及佣金支出
	POUNINCO     dbr.NullFloat64 //CoRe		手续费及佣金收入
	BIZCOST      dbr.NullFloat64 //CORe		营业成本
	DILUTEDEPS   dbr.NullFloat64 //DPES		稀释每股收益
	BASICEPS     dbr.NullFloat64 //EPS		基本每股收益
	FINEXPE      dbr.NullFloat64 //FnEp		财务费用
	__none__ICEp dbr.NullFloat64 //ICEp		保险手续费及佣金支出
	POLIDIVIEXPE dbr.NullFloat64 //IDEp		保单红利支出
	INTEINCO     dbr.NullFloat64 //InRe		利息收入
	INTEEXPE     dbr.NullFloat64 //ItEp		利息支出
	INCOTAXEXPE  dbr.NullFloat64 //ITEp		所得税费用
	MANAEXPE     dbr.NullFloat64 //MgEp		管理费用
	MINYSHARRIGH dbr.NullFloat64 //MIIn		少数股东损益
	__none__NCoE dbr.NullFloat64 //NCoE		手续费及佣金净收入 银行,保险利润表 存在该字段
	__none__NInR dbr.NullFloat64 //NInR		利息净收入 银行,保险,证券利润表 存在该字段
	NONOEXPE     dbr.NullFloat64 //NOEp		营业外支出
	NONOREVE     dbr.NullFloat64 //NORe		营业外收入
	NETPROFIT    dbr.NullFloat64 //NtIn		净利润
	BIZTAX       dbr.NullFloat64 //OATx		营业税金及附加
	BIZTOTCOST   dbr.NullFloat64 //OCOR		营业总成本
	__none__OOCs dbr.NullFloat64 //OOCs		其他营业成本
	__none__OpEp dbr.NullFloat64 //OpEp		营业支出 银行,保险,证券利润表存在此字段
	PERPROFIT    dbr.NullFloat64 //OpPr		营业利润
	BIZINCO      dbr.NullFloat64 //OpRe		营业收入
	SALESEXPE    dbr.NullFloat64 //SaEp		销售费用
	__none__SAPC dbr.NullFloat64 //SAPC		影响母公司净利润的特殊项目
	BIZTOTINCO   dbr.NullFloat64 //TOpR		营业总收入
	TOTPROFIT    dbr.NullFloat64 //ToPr		利润总额
}

func NewTQ_FIN_PROINCSTATEMENTNEW() *TQ_FIN_PROINCSTATEMENTNEW {
	return &TQ_FIN_PROINCSTATEMENTNEW{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROINCSTATEMENTNEW,
			Db:        MyCat,
		},
	}
}

func (this *TQ_FIN_PROINCSTATEMENTNEW) getListByCompcode(compcode string, report_type int, per_page int, page int) ([]TQ_FIN_PROINCSTATEMENTNEW, error) {
	var sli []TQ_FIN_PROINCSTATEMENTNEW

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

func (this *TQ_FIN_PROINCSTATEMENTNEW) GetList(scode string, market string, report_type int, per_page int, page int) ([]TQ_FIN_PROINCSTATEMENTNEW, error) {

	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode); err != nil {
		logging.Error("%T GetList error: %s", *this, err)
		return nil, err
	}

	return this.getListByCompcode(sc.COMPCODE.String, report_type, per_page, page)
}

// TQ_FIN_PROBINCSTATEMENTNEW    中文名称：银行利润表(新准则产品表)
// TQ_FIN_PROIINCSTATEMENTNEW    中文名称：保险利润表(新准则产品表)
// TQ_FIN_PROSINCSTATEMENTNEW    中文名称：证券利润表(新准则产品表)

/***************************以下是移动端f10页面*****************************************/
// 该处是财务数据

type F10_MB_PROINCSTATEMENTNEW struct {
	Model    `db:"-" `
	BIZINCO  dbr.NullFloat64 //营业收入（元）
	BASICEPS dbr.NullFloat64 //基本每股收益（元）
	PARENETP dbr.NullFloat64 //归属净利润（元）
}

func NewF10_MB_PROINCSTATEMENTNEW() *F10_MB_PROINCSTATEMENTNEW {
	return &F10_MB_PROINCSTATEMENTNEW{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROINCSTATEMENTNEW,
			Db:        MyCat,
		},
	}
}

func (this *F10_MB_PROINCSTATEMENTNEW) GetF10_MB_PROINCSTATEMENTNEW(compCode string) ([]F10_MB_PROINCSTATEMENTNEW, error) {
	var res []F10_MB_PROINCSTATEMENTNEW

	exps := map[string]interface{}{
		"COMPCODE=?":   compCode,
		"REPORTTYPE=?": 1,
		"ISVALID=?":    1,
	}
	builder := this.Db.Select("BIZINCO,BASICEPS,PARENETP").From(this.TableName).OrderBy("ENDDATE desc") //变动起始日
	err := this.SelectWhere(builder, exps).Limit(5).LoadStruct(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
