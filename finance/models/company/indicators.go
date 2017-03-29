// 关键指标
package company

import (
	"time"

	"haina.com/market/finance/models/finchina"
)

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
	NPOR      float64 `json:"NPDBR"`     //净利润／营业总收入
	NPORTTM   float64 `json:"NPDBRTTM"`  //净利润／营业总收入_TTM
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

func NewIndicatorsJson() *IndicatorsJson {
	return &IndicatorsJson{}
}

//------------------------------------------------------------------------------
type Indicators struct {
}

func NewIndicators() *Indicators {
	return &Indicators{}
}

func (this *Indicators) getList(scode string, report_type int, per_page int, page int) ([]Indicators, error) {
	return nil, nil
}
func (this *Indicators) GetList(scode string, report_type int, per_page int, page int) ([]Indicators, error) {
	return nil, nil
}
func (this *Indicators) getJson(scode string, report_type int, per_page int, page int) ([]IndicatorsJson, error) {
	return NewFinChinaIndicators().getJson(scode, report_type, per_page, page)
}
func (this *Indicators) GetJson(scode string, report_type int, per_page int, page int) (*RespFinAnaJson, error) {
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
type FinChinaIndicators struct {
}

func NewFinChinaIndicators() *FinChinaIndicators {
	return &FinChinaIndicators{}
}

func (this *FinChinaIndicators) getJson(scode string, report_type int, per_page int, page int) ([]IndicatorsJson, error) {
	sli := make([]IndicatorsJson, 0, per_page)
	ls, err := finchina.NewIndicators().GetList(scode, report_type, per_page, page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		node := IndicatorsJson{
			// TQ_FIN_PROFINMAININDEX     主要财务指标（产品表）
			DEPS:        v.PROFINMAININDEX.EPSFULLDILUTED.Float64, //稀释每股收益
			EPS:         v.PROFINMAININDEX.EPSBASIC.Float64,       //基本每股收益
			PSNA:        v.PROFINMAININDEX.NAPS.Float64,           //每股净资产(值)
			ONIDTPTTM:   v.PROFINMAININDEX.NVALCHGITOTP.Float64,   //价值变动净收益／利润总额_TTM
			NIVCDTPTTM:  v.PROFINMAININDEX.NNONOPITOTP.Float64,    //营业外收支净额／利润总额_TTM
			NPADNRGALNP: v.PROFINMAININDEX.OPANITOTP.Float64,      //经营活动净收益／利润总额_TTM

			// TQ_FIN_PROINDICDATA      衍生财务指标（产品表）
			//// 每股指标
			PSCR:   v.PROINDICDATA.CRPS.Float64,     //每股资本公积金
			PSECF:  v.PROINDICDATA.FCFFPS.Float64,   //每股企业自由现金流量
			PSR:    v.PROINDICDATA.OPREVPS.Float64,  //每股营业收入
			PSRE:   v.PROINDICDATA.REPS.Float64,     //每股留存收益
			PSSCF:  v.PROINDICDATA.FCFEPS.Float64,   //每股股东自由现金流量
			PSSR:   v.PROINDICDATA.SRPS.Float64,     //每股盈余公积金
			PSTR:   v.PROINDICDATA.TOPREVPS.Float64, //每股营业总收入
			PSUP:   v.PROINDICDATA.UPPS.Float64,     //每股未分配利润
			PSNCF:  v.PROINDICDATA.NCFPS.Float64,    //每股现金流量净额
			PSNBCF: v.PROINDICDATA.OPNCFPS.Float64,  //每股经营活动产生的现金流量净额

			//// 盈利能力
			CPR:      v.PROINDICDATA.PROTOTCRT.Float64,      // 成本费用利润率
			EBITDR:   v.PROINDICDATA.EBITTOTOPI.Float64,     // 息税前利润／营业总收入
			JROA:     v.PROINDICDATA.ROAAANNUAL.Float64,     // 总资产净利率
			NPADNRGL: v.PROINDICDATA.NPCUT.Float64,          // 扣除非经常性损益后的净利润
			OPR:      v.PROINDICDATA.OPPRORT.Float64,        // 营业利润率
			ROEA:     v.PROINDICDATA.ROEAVG.Float64,         // 净资产收益率_平均
			ROED:     v.PROINDICDATA.ROEDILUTED.Float64,     // 净资产收益率_摊薄
			ROEDD:    v.PROINDICDATA.ROEDILUTEDCUT.Float64,  // 净资产收益率_扣除摊薄
			ROEDW:    v.PROINDICDATA.ROEWEIGHTEDCUT.Float64, // 净资产收益率_扣除加权
			ROEW:     v.PROINDICDATA.ROEWEIGHTED.Float64,    // 净资产收益率_加权
			SCR:      v.PROINDICDATA.SCOSTRT.Float64,        // 销售成本率
			SGPM:     v.PROINDICDATA.SGPMARGIN.Float64,      // 销售毛利率

			//// 偿债能力
			CFLR:   v.PROINDICDATA.OPNCFTOCURLIAB.Float64,  //    现金流动负债比
			CR:     v.PROINDICDATA.CURRENTRT.Float64,       //    流动比率
			ER:     v.PROINDICDATA.EQURT.Float64,           //    产权比率
			LTLRWC: v.PROINDICDATA.LTMLIABTOOPCAP.Float64,  //    长期负债与营运资金比率 ???
			QR:     v.PROINDICDATA.QUICKRT.Float64,         //    速动比率
			TNWDND: v.PROINDICDATA.NTANGASSTONDEBT.Float64, //    有形净值／净债务
			TNWDR:  v.PROINDICDATA.TDEBTTOFART.Float64,     //    有形净值债务率

			//// 成长能力
			//// 营运能力
			APTD: v.PROINDICDATA.ACCPAYTDAYS.Float64,     // 应付帐款周转天数
			APTR: v.PROINDICDATA.ACCPAYRT.Float64,        // 应付帐款周转率
			ARTD: v.PROINDICDATA.ACCRECGTURNDAYS.Float64, // 应收帐款周转天数
			ARTR: v.PROINDICDATA.ACCRECGTURNRT.Float64,   // 应收帐款周转率
			CATR: v.PROINDICDATA.CURASSTURNRT.Float64,    // 流动资产周转率
			FATR: v.PROINDICDATA.FATURNRT.Float64,        // 固定资产周转率
			ITD:  v.PROINDICDATA.INVTURNDAYS.Float64,     // 存货周转天数
			ITR:  v.PROINDICDATA.INVTURNRT.Float64,       // 存货周转率
			OC:   v.PROINDICDATA.OPCYCLE.Float64,         // 营业周期
			SETR: v.PROINDICDATA.EQUTURNRT.Float64,       // 股东权益周转率
			TATR: v.PROINDICDATA.TATURNRT.Float64,        // 总资产周转率

			//// 现金状况
			FCFl:      v.PROINDICDATA.FCFF.Float64,            // 自由现金流量
			NBAGCFDNE: v.PROINDICDATA.OPANCFTOOPNI.Float64,    // 经营活动产生的现金流量净额／经营活动净收益
			SGPCRSDR:  v.PROINDICDATA.SCASHREVTOOPIRT.Float64, // 销售商品提供劳务收到的现金／营业收入

			//// 分红能力
			CDPM: v.PROINDICDATA.CDCOVER.Float64,  // 现金股利保障倍数
			DPM:  v.PROINDICDATA.DIVCOVER.Float64, // 股利保障倍数
			DPR:  v.PROINDICDATA.DPR.Float64,      // 股利支付率
			DPS:  v.PROINDICDATA.CDPS.Float64,     // 每股股利

			//// 资本结构
			ALR:    v.PROINDICDATA.ASSLIABRT.Float64,    // 资产负债率
			BPDTA:  v.PROINDICDATA.TDEBTTOTA.Float64,    // 应付债券／总资产
			EM:     v.PROINDICDATA.EM.Float64,           // 权益乘数
			FAR:    v.PROINDICDATA.LTMLIABTOTFA.Float64, // 固定资产比率 长期负债与固定资产比率???
			IAR:    v.PROINDICDATA.TCAPTOTART.Float64,   // 无形资产比率 资本与资产比率???
			LTASR:  v.PROINDICDATA.LTMASSRT.Float64,     // 长期资产适合率
			LTBDTA: v.PROINDICDATA.LTMLIABTOTA.Float64,  // 长期借款／总资产 长期负债/总资产???
			WC:     v.PROINDICDATA.WORKCAP.Float64,      // 营运资金

			//// 收益质量
			IDP:        v.PROINDICDATA.INCOTAXTOTP.Float64,  // 所得税／利润总额
			ONIDTP:     v.PROINDICDATA.NVALCHGITOTP.Float64, // 价值变动净收益／利润总额
			NIVCDTP:    v.PROINDICDATA.NNONOPITOTP.Float64,  // 营业外收支净额／利润总额
			NNOIDTP:    v.PROINDICDATA.NPCUTTONP.Float64,    // 扣除非经常损益后的净利润／净利润 扣除非经常性损益后的净利润/归属母公司的净利润???
			NNOIDTPTTM: v.PROINDICDATA.OPANITOTP.Float64,    // 经营活动净收益／利润总额

			//// 杜邦分析
			BPCNPSNP: v.PROINDICDATA.NPTONOCONMS.Float64, // 归属母公司股东的净利润／净利润 归属母公司股东的净利润/含少数股东损益的净利润???
			EMDA:     v.PROINDICDATA.EMCONMS.Float64,     // 权益乘数_杜邦分析 资本结构里已使用了EM(权益乘数) 此处为权益乘数(含少数股权的净资产)!!!???
			NIDTP:    v.PROINDICDATA.NPTOTP.Float64,      // 净利润／利润总额 归属母公司的净利润/利润总额???
			NPDBR:    v.PROINDICDATA.OPNCFTOOPTI.Float64, // 净利润／营业总收入 经营性现金净流量/营业总收入???

			// TQ_FIN_PROTTMINDIC     财务数据_TTM指标（产品表）
			//// 每股指标
			PSNBCFTTM: v.PROTTMINDIC.OPNCFPS.Float64, // 每股经营活动产生的现金流量净额_TTM
			PSNCFTTM:  v.PROTTMINDIC.NCFPS.Float64,   // 每股现金流量净额_TTM

			//// 盈利能力
			EBITDRTTM: v.PROTTMINDIC.EBITTOTOPI.Float64, // 息税前利润／营业总收入_TTM
			SGPMTTM:   v.PROTTMINDIC.SGPMARGIN.Float64,  // 销售毛利率_TTM

			//// 偿债能力
			//// 成长能力
			//// 营运能力
			//// 现金状况
			NBAGCFDNETTM: v.PROTTMINDIC.OPANCFTOOPNI.Float64,    // 经营活动产生的现金流量净额／经营活动净收益_TTM
			SGPCRSDRTTM:  v.PROTTMINDIC.SCASHREVTOOPIRT.Float64, // 销售商品提供劳务收到的现金／营业收入_TTM

			//// 分红能力
			//// 资本结构
			//// 收益质量
			//// 杜邦分析

			// TQ_FIN_PROCFSTTMSUBJECT	  TTM现金科目产品表
			CACENI: v.PROCFSTTMSUBJECT.CASHNETR.Float64, //     现金及现金等价物 净增加额
		}
		if v.PROFINMAININDEX.ENDDATE.Valid {
			tm, err := time.Parse("20060102", v.PROFINMAININDEX.ENDDATE.String)
			if err != nil {
				return nil, err
			}
			node.Date = tm.Unix()
		}

		sli = append(sli, node)
	}

	return sli, nil
}
