// 资产负债表
package company

import (
	"time"

	"github.com/gin-gonic/gin"
	"haina.com/market/finance/models"
	"haina.com/market/finance/models/company"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Liabilities struct {
}

func NewLiabilities() *Liabilities {
	return &Liabilities{}
}

// 获取资产负债列表
func (this *Liabilities) GET(c *gin.Context) {
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

func (this *Liabilities) getJson(req *RequestParam) (*company.RespFinAnaJson, error) {
	sli := make([]company.LiabilitiesJson, 0, req.PerPage)
	ls, err := company.NewLiabilities().GetList(req.SCode, req.Type, req.PerPage, req.Page)
	if err != nil {
		return nil, err
	}
	for _, v := range ls {
		one := company.LiabilitiesJson{
			//资产
			AcRe: v.ACCORECE.Float64,      // 应收账款
			DFAs: v.DERIFINAASSET.Float64, // 衍生金融资产
			DiRe: v.DIVIDRECE.Float64,     // 应收股利
			DTAs: v.DEFETAXASSET.Float64,  // 递延所得税资产
			FAFS: v.AVAISELLASSE.Float64,  // 可供出售金融资产
			FAHT: v.TRADFINASSET.Float64,  // 交易性金融资产
			FiAs: v.FIXEDASSEIMMO.Float64, // 固定资产 表中只有 固定资产原值??
			GWil: v.GOODWILL.Float64,      // 商誉
			HTMI: v.HOLDINVEDUE.Float64,   // 持有至到期投资
			InAs: v.INTAASSET.Float64,     // 无形资产
			InRe: v.INTERECE.Float64,      // 应收利息
			LdLt: v.PLAC.Float64,          // 拆出资金
			LTAR: v.LONGRECE.Float64,      // 长期应收款
			LTEI: v.EQUIINVE.Float64,      // 长期股权投资
			LTPE: v.LOGPREPEXPE.Float64,   // 长期待摊费用
			MnFd: v.CURFDS.Float64,        // 货币资金
			NoRe: v.NOTESRECE.Float64,     // 应收票据
			OtRe: v.OTHERRECE.Float64,     // 其他应收款
			PrEx: v.PREPEXPE.Float64,      // 待摊费用
			Prpy: v.PREP.Float64,          // 预付款项
			REFI: v.INVEPROP.Float64,      // 投资性房地产
			ToAs: v.TOTASSET.Float64,      // 资产总计

			//负债
			AcEx:  v.ACCREXPE.Float64,         // 预提费用
			AcPy:  v.ACCOPAYA.Float64,         // 应付账款
			AdRE:  v.ADVAPAYM.Float64,         // 预收款项
			BdPy:  v.BDSPAYA.Float64,          // 应付债券
			CmPy:  v.COPEPOUN.Float64,         // 应付手续费及佣金
			DETLb: v.DEFEINCOTAXLIAB.Float64,  // 递延所得税负债
			DfIn:  v.DEFEREVE.Float64,         // 递延收益
			DFLb:  v.DERILIAB.Float64,         // 衍生金融负债
			DpFB:  v.DEPOSIT.Float64,          // 同业及其他金融机构存放款项 吸收存款及同业存放 ???
			DvPy:  v.DIVIPAYA.Float64,         // 应付股利
			FASR:  v.SELLREPASSE.Float64,      // 卖出回购金融资产款
			InPy:  v.INTEPAYA.Float64,         // 应付利息
			LnFB:  v.FDSBORR.Float64,          // 拆入资金
			LnFC:  v.CENBANKBORR.Float64,      // 向中央银行借款
			LTBw:  v.LONGBORR.Float64,         // 长期借款
			LTPy:  v.LONGPAYA.Float64,         // 长期应付款
			NCL1:  v.DUENONCLIAB.Float64,      // 一年内到期的非流动负债
			NtPy:  v.NOTESPAYA.Float64,        // 应付票据
			PCSc:  v.BDSPAYAPERBOND.Float64,   // 永续债
			PrSk:  v.BDSPAYAPREST.Float64,     // 优先股
			SaPy:  v.COPEWORKERSAL.Float64,    // 应付职工薪酬
			SBPy:  v.SHORTTERMBDSPAYA.Float64, // 应付短期债券
			STLn:  v.SHORTTERMBORR.Float64,    // 短期借款
			TaLb:  v.TOTLIAB.Float64,          // 负债合计
			TFLb:  v.TRADFINLIAB.Float64,      // 交易性金融负债
			TxPy:  v.TAXESPAYA.Float64,        // 应交税费

			//所有者权益
			BPCSET: v.PARESHARRIGH.Float64,    //  归属于母公司股东权益合计
			CDFCS:  v.CURTRANDIFF.Float64,     //  外币报表折算差额
			CpSp:   v.CAPISURP.Float64,        //  资本公积
			GRPr:   v.GENERISKRESE.Float64,    //  一般风险准备
			MiIt:   v.MINYSHARRIGH.Float64,    //  少数股东权益
			OEIn:   v.OTHEQUIN.Float64,        //  其他权益工具
			OESET:  v.RIGHAGGR.Float64,        //  所有者权益（或股东权益）合计
			OtCI:   v.OCL.Float64,             //  其他综合收益
			PCSe:   v.PERBOND.Float64,         //  永续债
			PICa:   v.PAIDINCAPI.Float64,      //  实收资本（或股本）
			PrSc:   v.PREST.Float64,           //  优先股
			SpRs:   v.RESE.Float64,            //  盈余公积
			TLSE:   v.TOTLIABSHAREQUI.Float64, //  负债和所有者权益（或股东权益）总计
			TrSc:   v.TREASTK.Float64,         //  库存股 表中名称(减：库存股)
			UdPr:   v.UNDIPROF.Float64,        //  未分配利润
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
