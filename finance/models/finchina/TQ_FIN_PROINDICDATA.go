// 关键指标
// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
//------------------------------------------------------------------------------
// 与以下3表组合中的数据构成关键指标数据
// TQ_FIN_PROCFSTTMSUBJECT	TTM现金科目产品表
// TQ_FIN_PROFINMAININDEX   主要财务指标（产品表）
// TQ_FIN_PROTTMINDIC       财务数据_TTM指标（产品表）
package finchina

import (
	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

type TQ_FIN_PROINDICDATA struct {
	Model   `db:"-"`
	ENDDATE dbr.NullString //Date 	放置本次财报的截止日期

	//// 每股指标
	CRPS     dbr.NullFloat64 //PSCR      每股资本公积金
	FCFFPS   dbr.NullFloat64 //PSECF     每股企业自由现金流量
	OPREVPS  dbr.NullFloat64 //PSR       每股营业收入
	REPS     dbr.NullFloat64 //PSRE      每股留存收益
	FCFEPS   dbr.NullFloat64 //PSSCF     每股股东自由现金流量
	SRPS     dbr.NullFloat64 //PSSR      每股盈余公积金
	TOPREVPS dbr.NullFloat64 //PSTR      每股营业总收入
	UPPS     dbr.NullFloat64 //PSUP      每股未分配利润
	NCFPS    dbr.NullFloat64 //PSNCF     每股现金流量净额
	OPNCFPS  dbr.NullFloat64 //PSNBCF    每股经营活动产生的现金流量净额

	//// 盈利能力
	PROTOTCRT      dbr.NullFloat64 //CPR       成本费用利润率
	EBITTOTOPI     dbr.NullFloat64 //EBITDR    息税前利润／营业总收入
	ROAAANNUAL     dbr.NullFloat64 //JROA      总资产净利率
	NPCUT          dbr.NullFloat64 //NPADNRGL  扣除非经常性损益后的净利润
	OPPRORT        dbr.NullFloat64 //OPR       营业利润率
	ROEAVG         dbr.NullFloat64 //ROEA      净资产收益率_平均
	ROEDILUTED     dbr.NullFloat64 //ROED      净资产收益率_摊薄
	ROEDILUTEDCUT  dbr.NullFloat64 //ROEDD     净资产收益率_扣除,摊薄
	ROEWEIGHTEDCUT dbr.NullFloat64 //ROEDW     净资产收益率_扣除,加权
	ROEWEIGHTED    dbr.NullFloat64 //ROEW      净资产收益率_加权
	SCOSTRT        dbr.NullFloat64 //SCR       销售成本率
	SGPMARGIN      dbr.NullFloat64 //SGPM      销售毛利率

	//// 偿债能力
	OPNCFTOCURLIAB  dbr.NullFloat64 //CFLR       现金流动负债比
	CURRENTRT       dbr.NullFloat64 //CR         流动比率
	EQURT           dbr.NullFloat64 //ER         产权比率
	LTMLIABTOOPCAP  dbr.NullFloat64 //LTLRWC     长期负债与营运资金比率 ???
	QUICKRT         dbr.NullFloat64 //QR         速动比率
	NTANGASSTONDEBT dbr.NullFloat64 //TNWDND     有形净值／净债务
	TDEBTTOFART     dbr.NullFloat64 //TNWDR      有形净值债务率

	//// 成长能力
	//// 营运能力
	ACCPAYTDAYS     dbr.NullFloat64 //APTD 应付帐款周转天数
	ACCPAYRT        dbr.NullFloat64 //APTR 应付帐款周转率
	ACCRECGTURNDAYS dbr.NullFloat64 //ARTD 应收帐款周转天数
	ACCRECGTURNRT   dbr.NullFloat64 //ARTR 应收帐款周转率
	CURASSTURNRT    dbr.NullFloat64 //CATR 流动资产周转率
	FATURNRT        dbr.NullFloat64 //FATR 固定资产周转率
	INVTURNDAYS     dbr.NullFloat64 //ITD  存货周转天数
	INVTURNRT       dbr.NullFloat64 //ITR  存货周转率
	OPCYCLE         dbr.NullFloat64 //OC   营业周期
	EQUTURNRT       dbr.NullFloat64 //SETR 股东权益周转率
	TATURNRT        dbr.NullFloat64 //TATR 总资产周转率

	//// 现金状况
	FCFF            dbr.NullFloat64 //FCFl         自由现金流量
	OPANCFTOOPNI    dbr.NullFloat64 //NBAGCFDNE    经营活动产生的现金流量净额／经营活动净收益
	SCASHREVTOOPIRT dbr.NullFloat64 //SGPCRSDR     销售商品提供劳务收到的现金／营业收入

	//// 分红能力
	CDCOVER  dbr.NullFloat64 //CDPM 现金股利保障倍数
	DIVCOVER dbr.NullFloat64 //DPM  股利保障倍数
	DPR      dbr.NullFloat64 //DPR  股利支付率
	CDPS     dbr.NullFloat64 //DPS  每股股利

	//// 资本结构
	ASSLIABRT    dbr.NullFloat64 //ALR     资产负债率
	TDEBTTOTA    dbr.NullFloat64 //BPDTA   应付债券／总资产
	EM           dbr.NullFloat64 //EM      权益乘数
	LTMLIABTOTFA dbr.NullFloat64 //FAR     固定资产比率 长期负债与固定资产比率???
	TCAPTOTART   dbr.NullFloat64 //IAR     无形资产比率 资本与资产比率???
	LTMASSRT     dbr.NullFloat64 //LTASR   长期资产适合率
	LTMLIABTOTA  dbr.NullFloat64 //LTBDTA  长期借款／总资产 长期负债/总资产???
	WORKCAP      dbr.NullFloat64 //WC      营运资金

	//// 收益质量
	INCOTAXTOTP  dbr.NullFloat64 //IDP         所得税／利润总额
	NVALCHGITOTP dbr.NullFloat64 //ONIDTP      价值变动净收益／利润总额
	NNONOPITOTP  dbr.NullFloat64 //NIVCDTP     营业外收支净额／利润总额
	NPCUTTONP    dbr.NullFloat64 //NNOIDTP     扣除非经常损益后的净利润／净利润 扣除非经常性损益后的净利润/归属母公司的净利润???
	OPANITOTP    dbr.NullFloat64 //NNOIDTPTTM  经营活动净收益／利润总额

	//// 杜邦分析
	NPTONOCONMS dbr.NullFloat64 //BPCNPSNP 归属母公司股东的净利润／净利润 归属母公司股东的净利润/含少数股东损益的净利润???
	EMCONMS     dbr.NullFloat64 //EMDA     权益乘数_杜邦分析 资本结构里已使用了EM(权益乘数) 此处为权益乘数(含少数股权的净资产)!!!???
	NPTOTP      dbr.NullFloat64 //NIDTP    净利润／利润总额 归属母公司的净利润/利润总额???
	OPNCFTOOPTI dbr.NullFloat64 //NPDBR    净利润／营业总收入 经营性现金净流量/营业总收入???

}

func NewTQ_FIN_PROINDICDATA() *TQ_FIN_PROINDICDATA {
	return &TQ_FIN_PROINDICDATA{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROINDICDATA,
			Db:        MyCat,
		},
	}
}

// 从 TQ_FIN_PROINDICDATA     衍生财务指标（产品表）    取数据
func (this *TQ_FIN_PROINDICDATA) getListByCompcode(compcode string, report_type int, per_page int, page int, date DateList) ([]TQ_FIN_PROINDICDATA, error) {
	var (
		sli_db []TQ_FIN_PROINDICDATA
		err    error
	)

	builder := this.Db.Select("*").From(this.TableName)
	if report_type != 0 {
		builder.Where("REPORTDATETYPE=?", report_type)
	}
	if len(date) > 0 {
		sets := ""
		for _, v := range date {
			sets += v + ","
		}
		sets = sets[:len(sets)-1]
		builder.Where("ENDDATE in (" + sets + ")")
	}
	err = builder.Where("COMPCODE=?", compcode).
		Where("ISVALID=1").
		Where("REPORTTYPE=?", 3).
		OrderBy("ENDDATE DESC").
		Limit(uint64(per_page)).
		LoadStruct(&sli_db)
	if err != nil && err != dbr.ErrNotFound {
		logging.Error("%T getListByCompcode: %v", *this, err)
		return nil, err
	}
	return sli_db, nil
}
func (this *TQ_FIN_PROINDICDATA) GetListByEnddates(scode string, market string, report_type int, per_page int, page int, date DateList) ([]TQ_FIN_PROINDICDATA, error) {
	sc := NewTQ_OA_STCODE()
	if err := sc.getCompcode(scode, market); err != nil {
		logging.Error("%T GetListByEnddates error: %s", *this, err)
		return nil, err
	}
	return this.getListByCompcode(sc.COMPCODE.String, report_type, per_page, page, date)
}

/***************************以下是移动端f10页面*****************************************/
// 该处是财务数据部分字段

type F10_MB_PROINDICDATA struct {
	Model     `db:"-" `
	ROEAVG    dbr.NullFloat64 //净资产收益率_平均(%)
	ASSLIABRT dbr.NullFloat64 //资产负债率(%)
	CRPS      dbr.NullFloat64 //每股资本公积金(元)
	UPPS      dbr.NullFloat64 //每股未分配利润(元)
	ENDDATE   dbr.NullString  //截止时间
}

func NewF10_MB_PROINDICDATA() *F10_MB_PROINDICDATA {
	return &F10_MB_PROINDICDATA{
		Model: Model{
			TableName: TABLE_TQ_FIN_PROINDICDATA,
			Db:        MyCat,
		},
	}
}

func (this *F10_MB_PROINDICDATA) GetF10_MB_PROINDICDATA(compCode string) (*F10_MB_PROINDICDATA, error) {
	exps := map[string]interface{}{
		"COMPCODE=?":   compCode,
		"REPORTTYPE=?": 3,
		"ISVALID=?":    1,
	}
	builder := this.Db.Select("*").From(this.TableName).OrderBy("ENDDATE desc") //变动起始日
	err := this.SelectWhere(builder, exps).Limit(1).LoadStruct(this)
	if err != nil {
		return nil, err
	}
	return this, nil
}
