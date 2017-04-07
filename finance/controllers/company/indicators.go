// 关键指标
package company

import (
	"time"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Indicators struct {
}

func NewIndicators() *Indicators {
	return &Indicators{}
}

func (this *Indicators) GET(c *gin.Context) {
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

func (this *Indicators) getJson(req *RequestParam) (*company.RespFinAnaJson, error) {
	sli := make([]company.IndicatorsJson, 0, req.PerPage)
	ls, err := company.NewIndicators().GetList(req.SCode, req.Type, req.PerPage, req.Page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		node := company.IndicatorsJson{
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

	// A股
	jsn := &company.RespFinAnaJson{
		MU:     "人民币元",
		AS:     "新会计准则",
		Length: len(sli),
		List:   sli,
	}

	return jsn, nil
}
