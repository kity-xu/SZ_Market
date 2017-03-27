// 资产负债数据
package finchina

import (
	"time"

	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

type LiabilitiesJson struct {
	Date int64 `json:"Date"`

	//资产
	AcRe float64 `json:"AcRe"` //应收账款
	CuDe float64 `json:"CuDe"` //客户资金存款
	DeMg float64 `json:"DeMg"` //存出保证金
	DFAs float64 `json:"DFAs"` //衍生金融资产
	DfSv float64 `json:"DfSv"` //定期存款
	DiRe float64 `json:"DiRe"` //应收股利
	DTAs float64 `json:"DTAs"` //递延所得税资产
	FAFS float64 `json:"FAFS"` //可供出售金融资产
	FAHT float64 `json:"FAHT"` //交易性金融资产
	FiAs float64 `json:"FiAs"` //固定资产
	GWil float64 `json:"GWil"` //商誉
	HTMI float64 `json:"HTMI"` //持有至到期投资
	IbDe float64 `json:"IbDe"` //存放同业款项
	InAs float64 `json:"InAs"` //无形资产
	InRe float64 `json:"InRe"` //应收利息
	LdLt float64 `json:"LdLt"` //拆出资金
	LTAR float64 `json:"LTAR"` //长期应收款
	LTEI float64 `json:"LTEI"` //长期股权投资
	LTPE float64 `json:"LTPE"` //长期待摊费用
	Metl float64 `json:"Metl"` //贵金属
	MnFd float64 `json:"MnFd"` //货币资金
	NoRe float64 `json:"NoRe"` //应收票据
	OtRe float64 `json:"OtRe"` //其他应收款
	PrEx float64 `json:"PrEx"` //待摊费用
	Prpy float64 `json:"Prpy"` //预付款项
	REFI float64 `json:"REFI"` //投资性房地产
	ToAs float64 `json:"ToAs"` //资产总计

	//负债
	AcEx  float64 `json:"AcEx"`  //预提费用
	AcPy  float64 `json:"AcPy"`  //应付账款
	AdRE  float64 `json:"AdRE"`  //预收款项
	BdPy  float64 `json:"BdPy"`  //应付债券
	CmPy  float64 `json:"CmPy"`  //应付手续费及佣金
	DETLb float64 `json:"DETLb"` //递延所得税负债
	DfIn  float64 `json:"DfIn"`  //递延收益
	DFLb  float64 `json:"DFLb"`  //衍生金融负债
	DpCl  float64 `json:"DpCl"`  //吸收存款
	DpFB  float64 `json:"DpFB"`  //同业及其他金融机构存放款项
	DvPy  float64 `json:"DvPy"`  //应付股利
	FASR  float64 `json:"FASR"`  //卖出回购金融资产款
	InPy  float64 `json:"InPy"`  //应付利息
	LnFB  float64 `json:"LnFB"`  //拆入资金
	LnFC  float64 `json:"LnFC"`  //向中央银行借款
	LTBw  float64 `json:"LTBw"`  //长期借款
	LTPy  float64 `json:"LTPy"`  //长期应付款
	NCL1  float64 `json:"NCL1"`  //一年内到期的非流动负债
	NtPy  float64 `json:"NtPy"`  //应付票据
	PCSc  float64 `json:"PCSc"`  //永续债
	PlLn  float64 `json:"PlLn"`  //质押借款
	PrSk  float64 `json:"PrSk"`  //优先股
	SaPy  float64 `json:"SaPy"`  //应付职工薪酬
	SBPy  float64 `json:"SBPy"`  //应付短期债券
	STLn  float64 `json:"STLn"`  //短期借款
	TaLb  float64 `json:"TaLb"`  //负债合计
	TFLb  float64 `json:"TFLb"`  //交易性金融负债
	TxPy  float64 `json:"TxPy"`  //应交税费

	//所有者权益
	BPCOEAI float64 `json:"BPCOEAI"` //归属于母公司所有者权益调整项目
	BPCOESI float64 `json:"BPCOESI"` //归属于母公司所有者权益特殊项目
	BPCSET  float64 `json:"BPCSET"`  //归属于母公司股东权益合计
	CDFCS   float64 `json:"CDFCS"`   //外币报表折算差额
	CpSp    float64 `json:"CpSp"`    //资本公积
	GRPr    float64 `json:"GRPr"`    //一般风险准备
	LEAI    float64 `json:"LEAI"`    //负债和权益调整项目
	LESI    float64 `json:"LESI"`    //负债和权益特殊项目
	MiIt    float64 `json:"MiIt"`    //少数股东权益
	OEAI    float64 `json:"OEAI"`    //所有者权益调整项目
	OEIn    float64 `json:"OEIn"`    //其他权益工具
	OESET   float64 `json:"OESET"`   //所有者权益（或股东权益）合计
	OtCI    float64 `json:"OtCI"`    //其他综合收益
	PCSe    float64 `json:"PCSe"`    //永续债
	PICa    float64 `json:"PICa"`    //实收资本（或股本）
	PrSc    float64 `json:"PrSc"`    //优先股
	SpRs    float64 `json:"SpRs"`    //盈余公积
	TLSE    float64 `json:"TLSE"`    //负债和所有者权益（或股东权益）总计
	TrSc    float64 `json:"TrSc"`    //库存股
	UdPr    float64 `json:"UdPr"`    //未分配利润
}

func NewLiabilitiesJson() *LiabilitiesJson {
	return &LiabilitiesJson{}
}

//  TQ_FIN_PROBALSHEETNEW	  一般企业资产负债表(新准则产品表)
// __none__ 前缀的字段是参考其他证券软件的F10功能定义的Json返回字段信息但在数据表中没有找到与之对应的字段,为了不打乱顺序,做个标注
type LiabilitiesGeneral struct {
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

func NewLiabilitiesGeneral() *LiabilitiesGeneral {
	return &LiabilitiesGeneral{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *LiabilitiesGeneral) getLiabilitiesJsonList(compcode string, req *RequestParam) ([]*LiabilitiesJson, error) {
	logging.Info("getLiabilitiesJsonList %T, compcode %s", *this, compcode)
	var (
		sli_db []LiabilitiesGeneral
	)
	sli := make([]*LiabilitiesJson, 0)

	builder := this.Db.Select("*").From(TABLE_TQ_FIN_PROBALSHEETNEW)
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
		one := NewLiabilitiesJson()

		//资产
		one.AcRe = v.ACCORECE.Float64      // 应收账款
		one.DFAs = v.DERIFINAASSET.Float64 // 衍生金融资产
		one.DiRe = v.DIVIDRECE.Float64     // 应收股利
		one.DTAs = v.DEFETAXASSET.Float64  // 递延所得税资产
		one.FAFS = v.AVAISELLASSE.Float64  // 可供出售金融资产
		one.FAHT = v.TRADFINASSET.Float64  // 交易性金融资产
		one.FiAs = v.FIXEDASSEIMMO.Float64 // 固定资产 表中只有 固定资产原值??
		one.GWil = v.GOODWILL.Float64      // 商誉
		one.HTMI = v.HOLDINVEDUE.Float64   // 持有至到期投资
		one.InAs = v.INTAASSET.Float64     // 无形资产
		one.InRe = v.INTERECE.Float64      // 应收利息
		one.LdLt = v.PLAC.Float64          // 拆出资金
		one.LTAR = v.LONGRECE.Float64      // 长期应收款
		one.LTEI = v.EQUIINVE.Float64      // 长期股权投资
		one.LTPE = v.LOGPREPEXPE.Float64   // 长期待摊费用
		one.MnFd = v.CURFDS.Float64        // 货币资金
		one.NoRe = v.NOTESRECE.Float64     // 应收票据
		one.OtRe = v.OTHERRECE.Float64     // 其他应收款
		one.PrEx = v.PREPEXPE.Float64      // 待摊费用
		one.Prpy = v.PREP.Float64          // 预付款项
		one.REFI = v.INVEPROP.Float64      // 投资性房地产
		one.ToAs = v.TOTASSET.Float64      // 资产总计

		//负债
		one.AcEx = v.ACCREXPE.Float64         // 预提费用
		one.AcPy = v.ACCOPAYA.Float64         // 应付账款
		one.AdRE = v.ADVAPAYM.Float64         // 预收款项
		one.BdPy = v.BDSPAYA.Float64          // 应付债券
		one.CmPy = v.COPEPOUN.Float64         // 应付手续费及佣金
		one.DETLb = v.DEFEINCOTAXLIAB.Float64 // 递延所得税负债
		one.DfIn = v.DEFEREVE.Float64         // 递延收益
		one.DFLb = v.DERILIAB.Float64         // 衍生金融负债
		one.DpFB = v.DEPOSIT.Float64          // 同业及其他金融机构存放款项 吸收存款及同业存放 ???
		one.DvPy = v.DIVIPAYA.Float64         // 应付股利
		one.FASR = v.SELLREPASSE.Float64      // 卖出回购金融资产款
		one.InPy = v.INTEPAYA.Float64         // 应付利息
		one.LnFB = v.FDSBORR.Float64          // 拆入资金
		one.LnFC = v.CENBANKBORR.Float64      // 向中央银行借款
		one.LTBw = v.LONGBORR.Float64         // 长期借款
		one.LTPy = v.LONGPAYA.Float64         // 长期应付款
		one.NCL1 = v.DUENONCLIAB.Float64      // 一年内到期的非流动负债
		one.NtPy = v.NOTESPAYA.Float64        // 应付票据
		one.PCSc = v.BDSPAYAPERBOND.Float64   // 永续债
		one.PrSk = v.BDSPAYAPREST.Float64     // 优先股
		one.SaPy = v.COPEWORKERSAL.Float64    // 应付职工薪酬
		one.SBPy = v.SHORTTERMBDSPAYA.Float64 // 应付短期债券
		one.STLn = v.SHORTTERMBORR.Float64    // 短期借款
		one.TaLb = v.TOTLIAB.Float64          // 负债合计
		one.TFLb = v.TRADFINLIAB.Float64      // 交易性金融负债
		one.TxPy = v.TAXESPAYA.Float64        // 应交税费

		//所有者权益
		one.BPCSET = v.PARESHARRIGH.Float64  //  归属于母公司股东权益合计
		one.CDFCS = v.CURTRANDIFF.Float64    //  外币报表折算差额
		one.CpSp = v.CAPISURP.Float64        //  资本公积
		one.GRPr = v.GENERISKRESE.Float64    //  一般风险准备
		one.MiIt = v.MINYSHARRIGH.Float64    //  少数股东权益
		one.OEIn = v.OTHEQUIN.Float64        //  其他权益工具
		one.OESET = v.RIGHAGGR.Float64       //  所有者权益（或股东权益）合计
		one.OtCI = v.OCL.Float64             //  其他综合收益
		one.PCSe = v.PERBOND.Float64         //  永续债
		one.PICa = v.PAIDINCAPI.Float64      //  实收资本（或股本）
		one.PrSc = v.PREST.Float64           //  优先股
		one.SpRs = v.RESE.Float64            //  盈余公积
		one.TLSE = v.TOTLIABSHAREQUI.Float64 //  负债和所有者权益（或股东权益）总计
		one.TrSc = v.TREASTK.Float64         //  库存股 表中名称(减：库存股)
		one.UdPr = v.UNDIPROF.Float64        //  未分配利润

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

type LiabilitiesInfo struct {
}

func NewLiabilitiesInfo() *LiabilitiesInfo {
	return &LiabilitiesInfo{}
}

func (this *LiabilitiesInfo) GetJson(req *RequestParam) (*ResponseFinAnaJson, error) {
	logging.Info("GetJson %T, RequestParam %+v", *this, *req)

	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(req.SCode); err != nil {
		return nil, err
	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, req.SCode)
		return nil, ErrNullComp
	}

	sli := NewLiabilitiesGeneral()
	list, err := sli.getLiabilitiesJsonList(sc.COMPCODE.String, req)
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
//  TQ_FIN_PROBBALBSHEETNEW	  银行资产负债表(新准则产品表)
//  TQ_FIN_PROIBALSHEETNEW	  保险资产负债表(新准则产品表)
//  TQ_FIN_PROSBALSHEETNEW	  证券资产负债表(新准则产品表)
