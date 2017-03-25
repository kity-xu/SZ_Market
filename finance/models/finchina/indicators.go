// 关键指标数据
package finchina

import (
	"time"

	. "haina.com/share/models"

	. "haina.com/market/finance/models/company"
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

// __none__ 前缀的字段是参考其他证券软件的F10功能定义的Json返回字段信息但在数据表中没有找到与之对应的字段,为了不打乱顺序,做个标注
type IndicatorsGeneral struct {
	Model `db:"-"`

	//// 每股指标
	__none__EPSE   dbr.NullFloat64 //EPSE      每股收益_期末股本摊薄
	__none__EPSTTM dbr.NullFloat64 //EPSTTM    每股收益_TTM
	__none__PSAF   dbr.NullFloat64 //PSAF      每股公积金
	__none__PSOE   dbr.NullFloat64 //PSOE      每股营业利润
	__none__PSRTTM dbr.NullFloat64 //PSRTTM    每股营业收入_TTM

	//// 盈利能力
	__none__DSR       dbr.NullFloat64 //DSR       销售期间费用率
	__none__DSRTTM    dbr.NullFloat64 //DSRTTM    销售期间费用率_TTM
	__none__FEDBR     dbr.NullFloat64 //FEDBR     财务费用／营业总收入
	__none__FEDBRTTM  dbr.NullFloat64 //FEDBRTTM  财务费用／营业总收入_TTM
	__none__JROATTM   dbr.NullFloat64 //JROATTM   总资产净利率_TTM
	__none__LAIDBR    dbr.NullFloat64 //LAIDBR    资产减值损失／营业总收入
	__none__LAIDBRTTM dbr.NullFloat64 //LAIDBRTTM 资产减值损失／营业总收入_TTM
	__none__MEDBR     dbr.NullFloat64 //MEDBR     管理费用／营业总收入
	__none__MEDBRTTM  dbr.NullFloat64 //MEDBRTTM  管理费用／营业总收入_TTM
	__none__NPAPC     dbr.NullFloat64 //NPAPC     归属母公司净利润
	__none__NPOR      dbr.NullFloat64 //NPOR      净利润／营业总收入
	__none__NPORTTM   dbr.NullFloat64 //NPORTTM   净利润／营业总收入_TTM
	__none__NSR       dbr.NullFloat64 //NSR       销售净利率
	__none__NSRTTM    dbr.NullFloat64 //NSRTTM    销售净利率_TTM
	__none__OCDBR     dbr.NullFloat64 //OCDBR     营业总成本／营业总收入
	__none__OCDBRTTM  dbr.NullFloat64 //OCDBRTTM  营业总成本／营业总收入_TTM
	__none__ROETTM    dbr.NullFloat64 //ROETTM    净资产收益率_TTM
	__none__SEDBR     dbr.NullFloat64 //SEDBR     销售费用／营业总收入
	__none__SEDBRTTM  dbr.NullFloat64 //SEDBRTTM  销售费用／营业总收入_TTM

	//// 偿债能力
	__none__BPCSDIBD   dbr.NullFloat64 //BPCSDIBD   归属母公司股东的权益／带息债务
	__none__BPCSDTL    dbr.NullFloat64 //BPCSDTL    归属母公司股东的权益／负债合计
	__none__NBAGCFDCL  dbr.NullFloat64 //NBAGCFDCL  经营活动产生现金流量净额／流动负债
	__none__NBAGCFDIBD dbr.NullFloat64 //NBAGCFDIBD 经营活动产生现金流量净额／带息债务
	__none__NBAGCFDND  dbr.NullFloat64 //NBAGCFDND  经营活动产生现金流量净额／净债务
	__none__NBAGCFDTL  dbr.NullFloat64 //NBAGCFDTL  经营活动产生现金流量净额／负债合计
	__none__SQR        dbr.NullFloat64 //SQR        超速动比率
	__none__TIE        dbr.NullFloat64 //TIE        利息保障倍数
	__none__TNWDIBD    dbr.NullFloat64 //TNWDIBD    有形净值／带息债务

	//// 成长能力
	__none__APCSNPYG      dbr.NullFloat64 //APCSNPYG      归属母公司股东的净利润同比增长
	__none__APCSNPYGD     dbr.NullFloat64 //APCSNPYGD     归属母公司股东的净利润(扣除)同比增长
	__none__BAGCFNYOYG    dbr.NullFloat64 //BAGCFNYOYG    经营活动产生的现金流量净额同比增长
	__none__BAGCFPSNYOYG  dbr.NullFloat64 //BAGCFPSNYOYG  每股经营活动产生的现金流量净额同比增长
	__none__BEPSYG        dbr.NullFloat64 //BEPSYG        基本每股收益同比增长
	__none__BIYG          dbr.NullFloat64 //BIYG          营业收入同比增长
	__none__BPYG          dbr.NullFloat64 //BPYG          营业利润同比增长
	__none__BSPCERBYGR    dbr.NullFloat64 //BSPCERBYGR    归属母公司股东的权益相对年初增长率
	__none__DEPSYG        dbr.NullFloat64 //DEPSYG        稀释每股收益同比增长
	__none__NAPSRBYGR     dbr.NullFloat64 //NAPSRBYGR     每股净资产相对年初增长率
	__none__NAYG          dbr.NullFloat64 //NAYG          净资产同比增长
	__none__NPYG          dbr.NullFloat64 //NPYG          净利润同比增长
	__none__OP5YSPBPCNPGA dbr.NullFloat64 //OP5YSPBPCNPGA 过去五年同期归属母公司净利润平均增幅
	__none__REDYG         dbr.NullFloat64 //REDYG         净资产收益率(摊薄)同比增长
	__none__SGR           dbr.NullFloat64 //SGR           可持续增长率
	__none__TARBYGR       dbr.NullFloat64 //TARBYGR       资产总计相对年初增长率
	__none__TAYG          dbr.NullFloat64 //TAYG          总资产同比增长
	__none__TPYG          dbr.NullFloat64 //TPYG          利润总额同比增长

	//// 营运能力

	//// 现金状况
	__none__CSDAA       dbr.NullFloat64 //CSDAA        资本支出／折旧和摊销
	__none__NBAGCF      dbr.NullFloat64 //NBAGCF       经营活动产生的现金流量净额
	__none__NBAGCFDR    dbr.NullFloat64 //NBAGCFDR     经营活动产生的现金流量净额／营业收入
	__none__NBAGCFDRTTM dbr.NullFloat64 //NBAGCFDRTTM  经营活动产生的现金流量净额／营业收入_TTM
	__none__NPCL        dbr.NullFloat64 //NPCL         净利润现金含量
	__none__SGPCRS      dbr.NullFloat64 //SGPCRS       销售商品提供劳务收到的现金
	__none__TACRR       dbr.NullFloat64 //TACRR        总资产现金回收率

	//// 分红能力
	__none__CCEB dbr.NullFloat64 //CCEB 每股现金及现金等价物 余额
	__none__RER  dbr.NullFloat64 //RER  留存盈余比率

	//// 资本结构
	__none__BPCSDIC   dbr.NullFloat64 //BPCSDIC 归属母公司股东的权益／全部投入资本
	__none__CADTA     dbr.NullFloat64 //CADTA   流动资产／总资产
	__none__CLDTL     dbr.NullFloat64 //CLDTL   流动负债／负债合计
	__none__IBDDIC    dbr.NullFloat64 //IBDDIC  带息债务／全部投入资本
	__none__NCADTA    dbr.NullFloat64 //NCADTA  非流动资产／总资产
	__none__NCLDTL    dbr.NullFloat64 //NCLDTL  非流动负债／负债合计
	__none__EQUTURNRT dbr.NullFloat64 //SHER    股东权益比率 股东权益周转率???
	__none__LTLDSET   dbr.NullFloat64 //LTLDSET 长期负债／股东权益合计

	//// 收益质量
	__none__IIJVCDTP    dbr.NullFloat64 //IIJVCDTP    对联营合营公司投资收益／利润总额
	__none__IIJVCDTPTTM dbr.NullFloat64 //IIJVCDTPTTM 对联营合营公司投资收益／利润总额_TTM

	//// 杜邦分析
}

// TQ_FIN_PROFINMAININDEX     主要财务指标（产品表）
type TQ_FIN_PROFINMAININDEX struct {
	ENDDATE    dbr.NullString //Date 	放置本次财报的截止日期
	CUR        dbr.NullString //币种 	放置本次财报的币种
	ACCSTACODE dbr.NullString //会计准则 代码 11001 中国会计准则(07年前版) 11002 中国会计准则(07年起版) 本表全为07年后的新会计准则

	//// 每股指标
	EPSFULLDILUTED dbr.NullFloat64 //DEPS      稀释每股收益
	EPSBASIC       dbr.NullFloat64 //EPS       基本每股收益
	NAPS           dbr.NullFloat64 //PSNA      每股净资产(值)

	//// 盈利能力

	//// 偿债能力

	//// 成长能力

	//// 营运能力

	//// 现金状况

	//// 分红能力

	//// 资本结构

	//// 收益质量
	NVALCHGITOTP dbr.NullFloat64 //ONIDTPTTM   价值变动净收益／利润总额_TTM
	NNONOPITOTP  dbr.NullFloat64 //NIVCDTPTTM  营业外收支净额／利润总额_TTM
	OPANITOTP    dbr.NullFloat64 //NPADNRGALNP 经营活动净收益／利润总额_TTM

	//// 杜邦分析

}

// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
type TQ_FIN_PROINDICDATA struct {
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

// TQ_FIN_PROTTMINDIC     财务数据_TTM指标（产品表）
type TQ_FIN_PROTTMINDIC struct {
	ENDDATE dbr.NullString //Date 	放置本次财报的截止日期

	//// 每股指标
	OPNCFPS dbr.NullFloat64 //PSNBCFTTM 每股经营活动产生的现金流量净额_TTM
	NCFPS   dbr.NullFloat64 //PSNCFTTM  每股现金流量净额_TTM

	//// 盈利能力
	EBITTOTOPI dbr.NullFloat64 //EBITDRTTM 息税前利润／营业总收入_TTM
	SGPMARGIN  dbr.NullFloat64 //SGPMTTM   销售毛利率_TTM

	//// 偿债能力
	//// 成长能力
	//// 营运能力

	//// 现金状况
	OPANCFTOOPNI    dbr.NullFloat64 //NBAGCFDNETTM 经营活动产生的现金流量净额／经营活动净收益_TTM
	SCASHREVTOOPIRT dbr.NullFloat64 //SGPCRSDRTTM  销售商品提供劳务收到的现金／营业收入_TTM

	//// 分红能力
	//// 资本结构
	//// 收益质量
	//// 杜邦分析
}

// TQ_FIN_PROCFSTTMSUBJECT	  TTM现金科目产品表
type TQ_FIN_PROCFSTTMSUBJECT struct {
	ENDDATE dbr.NullString //Date 	放置本次财报的截止日期

	CASHNETR dbr.NullFloat64 //CACENI       现金及现金等价物 净增加额
}

func NewIndicatorsGeneral() *IndicatorsGeneral {
	return &IndicatorsGeneral{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

// TQ_FIN_PROFINMAININDEX     主要财务指标（产品表）
func (this *IndicatorsGeneral) getDataListForTQ_FIN_PROFINMAININDEX(compcode string, req *RequestParam) ([]*TQ_FIN_PROFINMAININDEX, error) {
	var (
		sli_db []*TQ_FIN_PROFINMAININDEX
		err    error
	)
	/*
		err = this.Db.Select("*").From(TABLE_TQ_FIN_PROFINMAININDEX).
			Where("COMPCODE = ?", compcode).
			//Where("REPORTTYPE = 1 or REPORTTYPE = 3").
			Where("REPORTTYPE = 3").
			OrderBy("REPORTTYPE ASE, ENDDATE DESC").
			Paginate(uint64(req.Page), uint64(req.PerPage)).
			LoadStruct(&sli_db)
			//*/

	builder := this.Db.Select("*").From(TABLE_TQ_FIN_PROFINMAININDEX)
	if req.Type != 0 {
		builder.Where("REPORTDATETYPE=?", req.Type)
	}
	//参考富途F10 财务分析 关键指标 的数据表 应由 finchina数据库的 主要财务指标表 衍生财务指标表 财务数据_TTM指标表 TTM现金科目产品表 四张表而来
	//表中 REPORTTYPE 给的释义: 放置本次财报的类型（1、3为合并报表；2、4为母公司报表；1和2是第一次披露的期末值，3和4是最新一次披露的数值，结合是否实际披露字段，可得治是否发生过再次披露）
	//按道理讲应该都用1(合并期末)
	//	finchina 数据库 主要财务指标表 REPORTTYPE 有记录类型 1,2,3,4
	//           而 衍生财务指标表 财务数据_TTM指标表 TTM现金科目产品表 REPORTTYPE 记录类型只有 3,4, 没有1,2
	//所以此处都统一用3(合并期末调整)
	err = builder.Where("COMPCODE = ?", compcode).
		//Where("REPORTTYPE = 1 or REPORTTYPE = 3").
		Where("REPORTTYPE = 3").
		OrderBy("REPORTTYPE ASC, ENDDATE DESC").
		Paginate(uint64(req.Page), uint64(req.PerPage)).
		LoadStruct(&sli_db)

	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}
	return sli_db, nil
}

// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
func (this *IndicatorsGeneral) getDataListForTQ_FIN_PROINDICDATA(compcode string, req *RequestParam) ([]*TQ_FIN_PROINDICDATA, error) {
	var (
		sli_db []*TQ_FIN_PROINDICDATA
		err    error
	)
	builder := this.Db.Select("*").From(TABLE_TQ_FIN_PROINDICDATA)
	if req.Type != 0 {
		builder.Where("REPORTDATETYPE=?", req.Type)
	}
	err = builder.Where("COMPCODE = ?", compcode).
		Where("REPORTTYPE = 3").
		OrderBy("REPORTTYPE ASC, ENDDATE DESC").
		Paginate(uint64(req.Page), uint64(req.PerPage)).
		LoadStruct(&sli_db)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}
	return sli_db, nil
}

// TQ_FIN_PROTTMINDIC     财务数据_TTM指标（产品表）
func (this *IndicatorsGeneral) getDataListForTQ_FIN_PROTTMINDIC(compcode string, req *RequestParam) ([]*TQ_FIN_PROTTMINDIC, error) {
	var (
		sli_db []*TQ_FIN_PROTTMINDIC
		err    error
	)

	builder := this.Db.Select("*").From(TABLE_TQ_FIN_PROTTMINDIC)
	if req.Type != 0 {
		builder.Where("REPORTDATETYPE=?", req.Type)
	}
	err = builder.Where("COMPCODE = ?", compcode).
		Where("REPORTTYPE = 3").
		OrderBy("REPORTTYPE ASC, ENDDATE DESC").
		Paginate(uint64(req.Page), uint64(req.PerPage)).
		LoadStruct(&sli_db)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}
	return sli_db, nil
}

// TQ_FIN_PROCFSTTMSUBJECT	  TTM现金科目产品表
func (this *IndicatorsGeneral) getDataListForTQ_FIN_PROCFSTTMSUBJECT(compcode string, req *RequestParam) ([]*TQ_FIN_PROCFSTTMSUBJECT, error) {
	var (
		sli_db []*TQ_FIN_PROCFSTTMSUBJECT
		err    error
	)

	builder := this.Db.Select("*").From(TABLE_TQ_FIN_PROCFSTTMSUBJECT)
	if req.Type != 0 {
		builder.Where("REPORTDATETYPE=?", req.Type)
	}
	err = builder.Where("COMPCODE = ?", compcode).
		Where("REPORTTYPE = 3").
		OrderBy("REPORTTYPE ASC, ENDDATE DESC").
		Paginate(uint64(req.Page), uint64(req.PerPage)).
		LoadStruct(&sli_db)
	if err != nil && err != dbr.ErrNotFound {
		return nil, err
	}
	return sli_db, nil
}

func (this *IndicatorsGeneral) getIndicatorsJsonList(compcode string, req *RequestParam) ([]*IndicatorsJson, error) {
	logging.Info("getIndicatorsJsonList %T, compcode %s", *this, compcode)
	var (
		sli_db_TQ_FIN_PROFINMAININDEX  []*TQ_FIN_PROFINMAININDEX
		sli_db_TQ_FIN_PROINDICDATA     []*TQ_FIN_PROINDICDATA
		sli_db_TQ_FIN_PROTTMINDIC      []*TQ_FIN_PROTTMINDIC
		sli_db_TQ_FIN_PROCFSTTMSUBJECT []*TQ_FIN_PROCFSTTMSUBJECT
		err                            error
		len1, len2, len3, len4         int
	)
	sli := make([]*IndicatorsJson, 0)

	if sli_db_TQ_FIN_PROFINMAININDEX, err = this.getDataListForTQ_FIN_PROFINMAININDEX(compcode, req); err != nil {
		return nil, err
	}
	if len1 = len(sli_db_TQ_FIN_PROFINMAININDEX); 0 == len1 {
		return sli, nil
	}

	if sli_db_TQ_FIN_PROINDICDATA, err = this.getDataListForTQ_FIN_PROINDICDATA(compcode, req); err != nil {
		return nil, err
	}
	if len2 = len(sli_db_TQ_FIN_PROINDICDATA); len2 != len1 {
		logging.Error("finchina db: TQ_FIN_PROINDICDATA %d != TQ_FIN_PROFINMAININDEX %d", len2, len1)
		return nil, ErrIncData
	}

	if sli_db_TQ_FIN_PROTTMINDIC, err = this.getDataListForTQ_FIN_PROTTMINDIC(compcode, req); err != nil {
		return nil, err
	}
	if len3 = len(sli_db_TQ_FIN_PROTTMINDIC); len3 != len1 {
		logging.Error("finchina db: TQ_FIN_PROTTMINDIC %d != TQ_FIN_PROFINMAININDEX %d", len3, len1)
		return nil, ErrIncData

	}

	if sli_db_TQ_FIN_PROCFSTTMSUBJECT, err = this.getDataListForTQ_FIN_PROCFSTTMSUBJECT(compcode, req); err != nil {
		return nil, err
	}
	if len4 = len(sli_db_TQ_FIN_PROCFSTTMSUBJECT); len4 != len1 {
		logging.Error("finchina db: TQ_FIN_PROCFSTTMSUBJECT %d != TQ_FIN_PROFINMAININDEX %d", len4, len1)
		return nil, ErrIncData

	}

	for n := 0; n < len1; n++ {
		v1 := sli_db_TQ_FIN_PROFINMAININDEX[n]
		v2 := sli_db_TQ_FIN_PROINDICDATA[n]
		v3 := sli_db_TQ_FIN_PROTTMINDIC[n]
		v4 := sli_db_TQ_FIN_PROCFSTTMSUBJECT[n]

		one := &IndicatorsJson{

			// TQ_FIN_PROFINMAININDEX     主要财务指标（产品表）
			DEPS:        v1.EPSFULLDILUTED.Float64, //稀释每股收益
			EPS:         v1.EPSBASIC.Float64,       //基本每股收益
			PSNA:        v1.NAPS.Float64,           //每股净资产(值)
			ONIDTPTTM:   v1.NVALCHGITOTP.Float64,   //价值变动净收益／利润总额_TTM
			NIVCDTPTTM:  v1.NNONOPITOTP.Float64,    //营业外收支净额／利润总额_TTM
			NPADNRGALNP: v1.OPANITOTP.Float64,      //经营活动净收益／利润总额_TTM

			// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
			//// 每股指标
			PSCR:   v2.CRPS.Float64,     //每股资本公积金
			PSECF:  v2.FCFFPS.Float64,   //每股企业自由现金流量
			PSR:    v2.OPREVPS.Float64,  //每股营业收入
			PSRE:   v2.REPS.Float64,     //每股留存收益
			PSSCF:  v2.FCFEPS.Float64,   //每股股东自由现金流量
			PSSR:   v2.SRPS.Float64,     //每股盈余公积金
			PSTR:   v2.TOPREVPS.Float64, //每股营业总收入
			PSUP:   v2.UPPS.Float64,     //每股未分配利润
			PSNCF:  v2.NCFPS.Float64,    //每股现金流量净额
			PSNBCF: v2.OPNCFPS.Float64,  //每股经营活动产生的现金流量净额

			//// 盈利能力
			CPR:      v2.PROTOTCRT.Float64,      // 成本费用利润率
			EBITDR:   v2.EBITTOTOPI.Float64,     // 息税前利润／营业总收入
			JROA:     v2.ROAAANNUAL.Float64,     // 总资产净利率
			NPADNRGL: v2.NPCUT.Float64,          // 扣除非经常性损益后的净利润
			OPR:      v2.OPPRORT.Float64,        // 营业利润率
			ROEA:     v2.ROEAVG.Float64,         // 净资产收益率_平均
			ROED:     v2.ROEDILUTED.Float64,     // 净资产收益率_摊薄
			ROEDD:    v2.ROEDILUTEDCUT.Float64,  // 净资产收益率_扣除摊薄
			ROEDW:    v2.ROEWEIGHTEDCUT.Float64, // 净资产收益率_扣除加权
			ROEW:     v2.ROEWEIGHTED.Float64,    // 净资产收益率_加权
			SCR:      v2.SCOSTRT.Float64,        // 销售成本率
			SGPM:     v2.SGPMARGIN.Float64,      // 销售毛利率

			//// 偿债能力
			CFLR:   v2.OPNCFTOCURLIAB.Float64,  //    现金流动负债比
			CR:     v2.CURRENTRT.Float64,       //    流动比率
			ER:     v2.EQURT.Float64,           //    产权比率
			LTLRWC: v2.LTMLIABTOOPCAP.Float64,  //    长期负债与营运资金比率 ???
			QR:     v2.QUICKRT.Float64,         //    速动比率
			TNWDND: v2.NTANGASSTONDEBT.Float64, //    有形净值／净债务
			TNWDR:  v2.TDEBTTOFART.Float64,     //    有形净值债务率

			//// 成长能力

			//// 营运能力
			APTD: v2.ACCPAYTDAYS.Float64,     // 应付帐款周转天数
			APTR: v2.ACCPAYRT.Float64,        // 应付帐款周转率
			ARTD: v2.ACCRECGTURNDAYS.Float64, // 应收帐款周转天数
			ARTR: v2.ACCRECGTURNRT.Float64,   // 应收帐款周转率
			CATR: v2.CURASSTURNRT.Float64,    // 流动资产周转率
			FATR: v2.FATURNRT.Float64,        // 固定资产周转率
			ITD:  v2.INVTURNDAYS.Float64,     // 存货周转天数
			ITR:  v2.INVTURNRT.Float64,       // 存货周转率
			OC:   v2.OPCYCLE.Float64,         // 营业周期
			SETR: v2.EQUTURNRT.Float64,       // 股东权益周转率
			TATR: v2.TATURNRT.Float64,        // 总资产周转率

			//// 现金状况
			FCFl:      v2.FCFF.Float64,            // 自由现金流量
			NBAGCFDNE: v2.OPANCFTOOPNI.Float64,    // 经营活动产生的现金流量净额／经营活动净收益
			SGPCRSDR:  v2.SCASHREVTOOPIRT.Float64, // 销售商品提供劳务收到的现金／营业收入

			//// 分红能力
			CDPM: v2.CDCOVER.Float64,  // 现金股利保障倍数
			DPM:  v2.DIVCOVER.Float64, // 股利保障倍数
			DPR:  v2.DPR.Float64,      // 股利支付率
			DPS:  v2.CDPS.Float64,     // 每股股利

			//// 资本结构
			ALR:    v2.ASSLIABRT.Float64,    // 资产负债率
			BPDTA:  v2.TDEBTTOTA.Float64,    // 应付债券／总资产
			EM:     v2.EM.Float64,           // 权益乘数
			FAR:    v2.LTMLIABTOTFA.Float64, // 固定资产比率 长期负债与固定资产比率???
			IAR:    v2.TCAPTOTART.Float64,   // 无形资产比率 资本与资产比率???
			LTASR:  v2.LTMASSRT.Float64,     // 长期资产适合率
			LTBDTA: v2.LTMLIABTOTA.Float64,  // 长期借款／总资产 长期负债/总资产???
			WC:     v2.WORKCAP.Float64,      // 营运资金

			//// 收益质量
			IDP:        v2.INCOTAXTOTP.Float64,  // 所得税／利润总额
			ONIDTP:     v2.NVALCHGITOTP.Float64, // 价值变动净收益／利润总额
			NIVCDTP:    v2.NNONOPITOTP.Float64,  // 营业外收支净额／利润总额
			NNOIDTP:    v2.NPCUTTONP.Float64,    // 扣除非经常损益后的净利润／净利润 扣除非经常性损益后的净利润/归属母公司的净利润???
			NNOIDTPTTM: v2.OPANITOTP.Float64,    // 经营活动净收益／利润总额

			//// 杜邦分析
			BPCNPSNP: v2.NPTONOCONMS.Float64, // 归属母公司股东的净利润／净利润 归属母公司股东的净利润/含少数股东损益的净利润???
			EMDA:     v2.EMCONMS.Float64,     // 权益乘数_杜邦分析 资本结构里已使用了EM(权益乘数) 此处为权益乘数(含少数股权的净资产)!!!???
			NIDTP:    v2.NPTOTP.Float64,      // 净利润／利润总额 归属母公司的净利润/利润总额???
			NPDBR:    v2.OPNCFTOOPTI.Float64, // 净利润／营业总收入 经营性现金净流量/营业总收入???

			// TQ_FIN_PROTTMINDIC     财务数据_TTM指标（产品表）
			//// 每股指标
			PSNBCFTTM: v3.OPNCFPS.Float64, // 每股经营活动产生的现金流量净额_TTM
			PSNCFTTM:  v3.NCFPS.Float64,   // 每股现金流量净额_TTM

			//// 盈利能力
			EBITDRTTM: v3.EBITTOTOPI.Float64, // 息税前利润／营业总收入_TTM
			SGPMTTM:   v3.SGPMARGIN.Float64,  // 销售毛利率_TTM

			//// 偿债能力
			//// 成长能力
			//// 营运能力
			//// 现金状况
			NBAGCFDNETTM: v3.OPANCFTOOPNI.Float64,    // 经营活动产生的现金流量净额／经营活动净收益_TTM
			SGPCRSDRTTM:  v3.SCASHREVTOOPIRT.Float64, // 销售商品提供劳务收到的现金／营业收入_TTM

			//// 分红能力
			//// 资本结构
			//// 收益质量
			//// 杜邦分析

			// TQ_FIN_PROCFSTTMSUBJECT	  TTM现金科目产品表
			CACENI: v4.CASHNETR.Float64, //     现金及现金等价物 净增加额
		}

		if v1.ENDDATE.Valid {
			tm, err := time.Parse("20060102", v1.ENDDATE.String)
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

type IndicatorsInfo struct {
}

func NewIndicatorsInfo() *IndicatorsInfo {
	return &IndicatorsInfo{}
}

func (this *IndicatorsInfo) GetJson(req *RequestParam) (*ResponseFinAnaJson, error) {
	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(req.SCode); err != nil {
		return nil, err
	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, req.SCode)
		return nil, ErrNullComp
	}

	sli := NewIndicatorsGeneral()
	list, err := sli.getIndicatorsJsonList(sc.COMPCODE.String, req)
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
