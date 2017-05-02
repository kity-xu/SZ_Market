// 关键指标
package company

import (
	"haina.com/market/finance/models/finchina"
	"haina.com/share/logging"
)

type Indicators struct {
	FinChinaIndicators
}

type IndicatorsJson struct {
	Date int64 `json:"Date"`

	// 每股指标
	DEPS      float64 `json:"DEPS"`      //稀释每股收益
	EPS       float64 `json:"EPS"`       //基本每股收益
	EPSE      float64 `json:"EPSE"`      //每股收益_期末股本摊薄
	EPSTTM    float64 `json:"EPSTTM"`    //每股收益_TTM
	PSAF      float64 `json:"PSAF"`      //每股公积金
	PSCR      float64 `json:"PSCR"`      //每股资本公积金
	PSECF     float64 `json:"PSECF"`     //每股企业自由现金流量
	PSNA      float64 `json:"PSNA"`      //每股净资产(值)
	PSNBCF    float64 `json:"PSNBCF"`    //每股经营活动产生的现金流量净额
	PSNBCFTTM float64 `json:"PSNBCFTTM"` //每股经营活动产生的现金流量净额_TTM
	PSNCF     float64 `json:"PSNCF"`     //每股现金流量净额
	PSNCFTTM  float64 `json:"PSNCFTTM"`  //每股现金流量净额_TTM
	PSOE      float64 `json:"PSOE"`      //每股营业利润
	PSR       float64 `json:"PSR"`       //每股营业收入
	PSRE      float64 `json:"PSRE"`      //每股留存收益
	PSRTTM    float64 `json:"PSRTTM"`    //每股营业收入_TTM
	PSSCF     float64 `json:"PSSCF"`     //每股股东自由现金流量
	PSSR      float64 `json:"PSSR"`      //每股盈余公积金
	PSTR      float64 `json:"PSTR"`      //每股营业总收入
	PSUP      float64 `json:"PSUP"`      //每股未分配利润

	// 盈利能力
	CPR       float64 `json:"CPR"`       //成本费用利润率
	DSR       float64 `json:"DSR"`       //销售期间费用率
	DSRTTM    float64 `json:"DSRTTM"`    //销售期间费用率_TTM
	EBITDR    float64 `json:"EBITDR"`    //息税前利润／营业总收入
	EBITDRTTM float64 `json:"EBITDRTTM"` //息税前利润／营业总收入_TTM
	FEDBR     float64 `json:"FEDBR"`     //财务费用／营业总收入
	FEDBRTTM  float64 `json:"FEDBRTTM"`  //财务费用／营业总收入_TTM
	JROA      float64 `json:"JROA"`      //总资产净利率
	JROATTM   float64 `json:"JROATTM"`   //总资产净利率_TTM
	LAIDBR    float64 `json:"LAIDBR"`    //资产减值损失／营业总收入
	LAIDBRTTM float64 `json:"LAIDBRTTM"` //资产减值损失／营业总收入_TTM
	MEDBR     float64 `json:"MEDBR"`     //管理费用／营业总收入
	MEDBRTTM  float64 `json:"MEDBRTTM"`  //管理费用／营业总收入_TTM
	NPADNRGL  float64 `json:"NPADNRGL"`  //扣除非经常性损益后的净利润
	NPAPC     float64 `json:"NPAPC"`     //归属母公司净利润
	NPOR      float64 `json:"NPOR"`      //净利润／营业总收入
	NPORTTM   float64 `json:"NPORTTM"`   //净利润／营业总收入_TTM
	NSR       float64 `json:"NSR"`       //销售净利率
	NSRTTM    float64 `json:"NSRTTM"`    //销售净利率_TTM
	OCDBR     float64 `json:"OCDBR"`     //营业总成本／营业总收入
	OCDBRTTM  float64 `json:"OCDBRTTM"`  //营业总成本／营业总收入_TTM
	OPR       float64 `json:"OPR"`       //营业利润率
	ROEA      float64 `json:"ROEA"`      //净资产收益率_平均
	ROED      float64 `json:"ROED"`      //净资产收益率_摊薄
	ROEDD     float64 `json:"ROEDD"`     //净资产收益率_扣除,摊薄
	ROEDW     float64 `json:"ROEDW"`     //净资产收益率_扣除,加权
	ROETTM    float64 `json:"ROETTM"`    //净资产收益率_TTM
	ROEW      float64 `json:"ROEW"`      //净资产收益率_加权
	SCR       float64 `json:"SCR"`       //销售成本率
	SEDBR     float64 `json:"SEDBR"`     //销售费用／营业总收入
	SEDBRTTM  float64 `json:"SEDBRTTM"`  //销售费用／营业总收入_TTM
	SGPM      float64 `json:"SGPM"`      //销售毛利率
	SGPMTTM   float64 `json:"SGPMTTM"`   //销售毛利率_TTM

	// 偿债能力
	BPCSDIBD   float64 `json:"BPCSDIBD"`   //归属母公司股东的权益／带息债务
	BPCSDTL    float64 `json:"BPCSDTL"`    //归属母公司股东的权益／负债合计
	CFLR       float64 `json:"CFLR"`       //现金流动负债比
	CR         float64 `json:"CR"`         //流动比率
	ER         float64 `json:"ER"`         //产权比率
	LTLRWC     float64 `json:"LTLRWC"`     //长期负债与营运资金比率
	NBAGCFDCL  float64 `json:"NBAGCFDCL"`  //经营活动产生现金流量净额／流动负债
	NBAGCFDIBD float64 `json:"NBAGCFDIBD"` //经营活动产生现金流量净额／带息债务
	NBAGCFDND  float64 `json:"NBAGCFDND"`  //经营活动产生现金流量净额／净债务
	NBAGCFDTL  float64 `json:"NBAGCFDTL"`  //经营活动产生现金流量净额／负债合计
	QR         float64 `json:"QR"`         //速动比率
	SQR        float64 `json:"SQR"`        //超速动比率
	TIE        float64 `json:"TIE"`        //利息保障倍数
	TNWDIBD    float64 `json:"TNWDIBD"`    //有形净值／带息债务
	TNWDND     float64 `json:"TNWDND"`     //有形净值／净债务
	TNWDR      float64 `json:"TNWDR"`      //有形净值债务率

	// 成长能力
	APCSNPYG      float64 `json:"APCSNPYG"`      //归属母公司股东的净利润同比增长
	APCSNPYGD     float64 `json:"APCSNPYGD"`     //归属母公司股东的净利润(扣除)同比增长
	BAGCFNYOYG    float64 `json:"BAGCFNYOYG"`    //经营活动产生的现金流量净额同比增长
	BAGCFPSNYOYG  float64 `json:"BAGCFPSNYOYG"`  //每股经营活动产生的现金流量净额同比增长
	BEPSYG        float64 `json:"BEPSYG"`        //基本每股收益同比增长
	BIYG          float64 `json:"BIYG"`          //营业收入同比增长
	BPYG          float64 `json:"BPYG"`          //营业利润同比增长
	BSPCERBYGR    float64 `json:"BSPCERBYGR"`    //归属母公司股东的权益相对年初增长率
	DEPSYG        float64 `json:"DEPSYG"`        //稀释每股收益同比增长
	NAPSRBYGR     float64 `json:"NAPSRBYGR"`     //每股净资产相对年初增长率
	NAYG          float64 `json:"NAYG"`          //净资产同比增长
	NPYG          float64 `json:"NPYG"`          //净利润同比增长
	OP5YSPBPCNPGA float64 `json:"OP5YSPBPCNPGA"` //过去五年同期归属母公司净利润平均增幅
	REDYG         float64 `json:"REDYG"`         //净资产收益率(摊薄)同比增长
	SGR           float64 `json:"SGR"`           //可持续增长率
	TARBYGR       float64 `json:"TARBYGR"`       //资产总计相对年初增长率
	TAYG          float64 `json:"TAYG"`          //总资产同比增长
	TPYG          float64 `json:"TPYG"`          //利润总额同比增长

	// 营运能力
	APTD float64 `json:"APTD"` //应付帐款周转天数
	APTR float64 `json:"APTR"` //应付帐款周转率
	ARTD float64 `json:"ARTD"` //应收帐款周转天数
	ARTR float64 `json:"ARTR"` //应收帐款周转率
	CATR float64 `json:"CATR"` //流动资产周转率
	FATR float64 `json:"FATR"` //固定资产周转率
	ITD  float64 `json:"ITD"`  //存货周转天数
	ITR  float64 `json:"ITR"`  //存货周转率
	OC   float64 `json:"OC"`   //营业周期
	SETR float64 `json:"SETR"` //股东权益周转率
	TATR float64 `json:"TATR"` //总资产周转率

	// 现金状况
	CACENI       float64 `json:"CACENI"`       //现金及现金等价物 净增加额
	CSDAA        float64 `json:"CSDAA"`        //资本支出／折旧和摊销
	FCFl         float64 `json:"FCFl"`         //自由现金流量
	NBAGCF       float64 `json:"NBAGCF"`       //经营活动产生的现金流量净额
	NBAGCFDNE    float64 `json:"NBAGCFDNE"`    //经营活动产生的现金流量净额／经营活动净收益
	NBAGCFDNETTM float64 `json:"NBAGCFDNETTM"` //经营活动产生的现金流量净额／经营活动净收益_TTM
	NBAGCFDR     float64 `json:"NBAGCFDR"`     //经营活动产生的现金流量净额／营业收入
	NBAGCFDRTTM  float64 `json:"NBAGCFDRTTM"`  //经营活动产生的现金流量净额／营业收入_TTM
	NPCL         float64 `json:"NPCL"`         //净利润现金含量
	SGPCRS       float64 `json:"SGPCRS"`       //销售商品提供劳务收到的现金
	SGPCRSDR     float64 `json:"SGPCRSDR"`     //销售商品提供劳务收到的现金／营业收入
	SGPCRSDRTTM  float64 `json:"SGPCRSDRTTM"`  //销售商品提供劳务收到的现金／营业收入_TTM
	TACRR        float64 `json:"TACRR"`        //总资产现金回收率

	// 分红能力
	CCEB float64 `json:"CCEB"` //每股现金及现金等价物 余额
	CDPM float64 `json:"CDPM"` //现金股利保障倍数
	DPM  float64 `json:"DPM"`  //股利保障倍数
	DPR  float64 `json:"DPR"`  //股利支付率
	DPS  float64 `json:"DPS"`  //每股股利
	RER  float64 `json:"RER"`  //留存盈余比率

	// 资本结构
	ALR     float64 `json:"ALR"`     //资产负债率
	BPCSDIC float64 `json:"BPCSDIC"` //归属母公司股东的权益／全部投入资本
	BPDTA   float64 `json:"BPDTA"`   //应付债券／总资产
	CADTA   float64 `json:"CADTA"`   //流动资产／总资产
	CLDTL   float64 `json:"CLDTL"`   //流动负债／负债合计
	EM      float64 `json:"EM"`      //权益乘数
	FAR     float64 `json:"FAR"`     //固定资产比率
	IAR     float64 `json:"IAR"`     //无形资产比率
	IBDDIC  float64 `json:"IBDDIC"`  //带息债务／全部投入资本
	LTASR   float64 `json:"LTASR"`   //长期资产适合率
	LTBDTA  float64 `json:"LTBDTA"`  //长期借款／总资产
	LTLDSET float64 `json:"LTLDSET"` //长期负债／股东权益合计
	NCADTA  float64 `json:"NCADTA"`  //非流动资产／总资产
	NCLDTL  float64 `json:"NCLDTL"`  //非流动负债／负债合计
	SHER    float64 `json:"SHER"`    //股东权益比率
	WC      float64 `json:"WC"`      //营运资金

	// 收益质量
	IDP         float64 `json:"IDP"`         //所得税／利润总额
	IIJVCDTP    float64 `json:"IIJVCDTP"`    //对联营合营公司投资收益／利润总额
	IIJVCDTPTTM float64 `json:"IIJVCDTPTTM"` //对联营合营公司投资收益／利润总额_TTM
	ONIDTP      float64 `json:"NIDTP"`       //价值变动净收益／利润总额
	ONIDTPTTM   float64 `json:"NIDTPTTM"`    //价值变动净收益／利润总额_TTM
	NIVCDTP     float64 `json:"NIVCDTP"`     //营业外收支净额／利润总额
	NIVCDTPTTM  float64 `json:"NIVCDTPTTM"`  //营业外收支净额／利润总额_TTM
	NNOIDTP     float64 `json:"NNOIDTP"`     //扣除非经常损益后的净利润／净利润
	NNOIDTPTTM  float64 `json:"NNOIDTPTTM"`  //经营活动净收益／利润总额
	NPADNRGALNP float64 `json:"NPADNRGALNP"` //经营活动净收益／利润总额_TTM

	// 杜邦分析
	BPCNPSNP float64 `json:"BPCNPSNP"` //归属母公司股东的净利润／净利润
	EMDA     float64 `json:"EMDA"`     //权益乘数_杜邦分析
	NIDTP    float64 `json:"NIDTP"`    //净利润／利润总额
	NPDBR    float64 `json:"NPDBR"`    //净利润／营业总收入
}

//------------------------------------------------------------------------------

func NewIndicators() *Indicators {
	return &Indicators{}
}

// 获取数据列表
func (this *Indicators) GetList(scode string, report_type int, per_page int, page int) ([]Indicators, error) {
	return NewFinChinaIndicators().getIndicatorsList(scode, report_type, per_page, page)
}

//------------------------------------------------------------------------------

type FinChinaIndicators struct {
	PROFINMAININDEX  finchina.TQ_FIN_PROFINMAININDEX  //主要财务指标（产品表）
	PROINDICDATA     finchina.TQ_FIN_PROINDICDATA     //衍生财务指标（产品表）
	PROTTMINDIC      finchina.TQ_FIN_PROTTMINDIC      //财务数据_TTM指标（产品表）
	PROCFSTTMSUBJECT finchina.TQ_FIN_PROCFSTTMSUBJECT //TTM现金科目产品表
}

func NewFinChinaIndicators() *FinChinaIndicators {
	return &FinChinaIndicators{}
}

func (this *FinChinaIndicators) getIndicatorsList(scode string, report_type int, per_page int, page int) ([]Indicators, error) {
	var (
		slidb_TQ_FIN_PROFINMAININDEX  []finchina.TQ_FIN_PROFINMAININDEX
		slidb_TQ_FIN_PROINDICDATA     []finchina.TQ_FIN_PROINDICDATA
		slidb_TQ_FIN_PROTTMINDIC      []finchina.TQ_FIN_PROTTMINDIC
		slidb_TQ_FIN_PROCFSTTMSUBJECT []finchina.TQ_FIN_PROCFSTTMSUBJECT
		dates                         finchina.DateList
		len1, len2, len3, len4        int
		err                           error
	)
	sli := make([]Indicators, 0, per_page)

	// 从 TQ_FIN_PROFINMAININDEX  主要财务指标（产品表）    取数据
	slidb_TQ_FIN_PROFINMAININDEX, err = finchina.NewTQ_FIN_PROFINMAININDEX().GetList(scode, report_type, per_page, page)
	if err != nil {
		return nil, err
	}
	if len1 = len(slidb_TQ_FIN_PROFINMAININDEX); 0 == len1 {
		return sli, nil
	}
	// 生成截止日期数组,其余表按日期数组取数据
	for _, v := range slidb_TQ_FIN_PROFINMAININDEX {
		if v.ENDDATE.Valid {
			date := v.ENDDATE.String
			dates = append(dates, date)
		}
	}

	// 从 TQ_FIN_PROINDICDATA     衍生财务指标（产品表）    取数据
	slidb_TQ_FIN_PROINDICDATA, err = finchina.NewTQ_FIN_PROINDICDATA().GetListByEnddates(scode, report_type, per_page, page, dates)
	if err != nil {
		return nil, err
	}
	if len2 = len(slidb_TQ_FIN_PROINDICDATA); len2 != len1 {
		logging.Error("finchina db: TQ_FIN_PROINDICDATA %d != TQ_FIN_PROFINMAININDEX %d", len2, len1)
		return nil, finchina.ErrIncData
	}

	// 从 TQ_FIN_PROTTMINDIC      财务数据_TTM指标（产品表）取数据
	slidb_TQ_FIN_PROTTMINDIC, err = finchina.NewTQ_FIN_PROTTMINDIC().GetListByEnddates(scode, report_type, per_page, page, dates)
	if err != nil {
		return nil, err
	}
	if len3 = len(slidb_TQ_FIN_PROTTMINDIC); len3 != len1 {
		logging.Error("finchina db: TQ_FIN_PROTTMINDIC %d != TQ_FIN_PROFINMAININDEX %d", len3, len1)
		return nil, finchina.ErrIncData

	}

	// 从 TQ_FIN_PROCFSTTMSUBJECT TTM现金科目产品表        取数据
	slidb_TQ_FIN_PROCFSTTMSUBJECT, err = finchina.NewTQ_FIN_PROCFSTTMSUBJECT().GetListByEnddates(scode, report_type, per_page, page, dates)
	if err != nil {
		return nil, err
	}
	if len4 = len(slidb_TQ_FIN_PROCFSTTMSUBJECT); len4 != len1 {
		logging.Error("finchina db: TQ_FIN_PROCFSTTMSUBJECT %d != TQ_FIN_PROFINMAININDEX %d", len4, len1)
		return nil, finchina.ErrIncData

	}

	for n := 0; n < len1; n++ {
		one := Indicators{
			FinChinaIndicators: FinChinaIndicators{
				PROFINMAININDEX:  slidb_TQ_FIN_PROFINMAININDEX[n],  //主要财务指标（产品表）
				PROINDICDATA:     slidb_TQ_FIN_PROINDICDATA[n],     //衍生财务指标（产品表）
				PROTTMINDIC:      slidb_TQ_FIN_PROTTMINDIC[n],      //财务数据_TTM指标（产品表）
				PROCFSTTMSUBJECT: slidb_TQ_FIN_PROCFSTTMSUBJECT[n], //TTM现金科目产品表
			},
		}
		sli = append(sli, one)
	}

	return sli, nil
}
