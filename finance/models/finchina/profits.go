// 利润数据
package finchina

import (
	"time"

	. "haina.com/share/models"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
)

type ProfitsJson struct {
	Date int64 `json:"Date"`

	AAPC float64 `json:"AAPC"` //影响母公司净利润的调整项目
	AILs float64 `json:"AILs"` //资产减值损失
	AREp float64 `json:"AREp"` //分保费用
	BAEp float64 `json:"BAEp"` //业务及管理费
	BPAC float64 `json:"BPAC"` //归属于母公司所有者的净利润
	CoEp float64 `json:"CoEp"` //手续费及佣金支出
	CoRe float64 `json:"CoRe"` //手续费及佣金收入
	CORe float64 `json:"CORe"` //营业成本
	DPES float64 `json:"DPES"` //稀释每股收益
	EPS  float64 `json:"EPS"`  //基本每股收益
	FnEp float64 `json:"FnEp"` //财务费用
	ICEp float64 `json:"ICEp"` //保险手续费及佣金支出
	IDEp float64 `json:"IDEp"` //保单红利支出
	InRe float64 `json:"InRe"` //利息收入
	ItEp float64 `json:"ItEp"` //利息支出
	ITEp float64 `json:"ITEp"` //所得税费用
	MgEp float64 `json:"MgEp"` //管理费用
	MIIn float64 `json:"MIIn"` //少数股东损益
	NCoE float64 `json:"NCoE"` //手续费及佣金净收入
	NInR float64 `json:"NInR"` //利息净收入
	NOEp float64 `json:"NOEp"` //营业外支出
	NORe float64 `json:"NORe"` //营业外收入
	NtIn float64 `json:"NtIn"` //净利润
	OATx float64 `json:"OATx"` //营业税金及附加
	OCOR float64 `json:"OCOR"` //营业总成本
	OOCs float64 `json:"OOCs"` //其他营业成本
	OpEp float64 `json:"OpEp"` //营业支出
	OpPr float64 `json:"OpPr"` //营业利润
	OpRe float64 `json:"OpRe"` //营业收入
	SaEp float64 `json:"SaEp"` //销售费用
	SAPC float64 `json:"SAPC"` //影响母公司净利润的特殊项目
	TOpR float64 `json:"TOpR"` //营业总收入
	ToPr float64 `json:"ToPr"` //利润总额
}

func NewProfitsJson() *ProfitsJson {
	return &ProfitsJson{}
}

// TQ_FIN_PROINCSTATEMENTNEW    中文名称：一般企业利润表(新准则产品表)
// __none__ 前缀的字段是参考其他证券软件的F10功能定义的Json返回字段信息但在数据表中没有找到与之对应的字段,为了不打乱顺序,做个标注
type ProfitsGeneral struct {
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

func NewProfitsGeneral() *ProfitsGeneral {
	return &ProfitsGeneral{
		Model: Model{
			TableName: TABLE_TQ_OA_STCODE,
			Db:        MyCat,
		},
	}
}

func (this *ProfitsGeneral) getProfitsJsonList(compcode string, req *RequestParam) ([]*ProfitsJson, error) {
	logging.Info("getProfitsJsonList %T, compcode %s", *this, compcode)
	var (
		sli_db []ProfitsGeneral
	)
	sli := make([]*ProfitsJson, 0)

	builder := this.Db.Select("*").From(TABLE_TQ_FIN_PROINCSTATEMENTNEW)
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
		one := NewProfitsJson()

		one.AILs = v.ASSEIMPALOSS.Float64
		one.AREp = v.REINEXPE.Float64
		one.BPAC = v.PARENETP.Float64
		one.CoEp = v.POUNEXPE.Float64
		one.CoRe = v.POUNINCO.Float64
		one.CORe = v.BIZCOST.Float64
		one.DPES = v.DILUTEDEPS.Float64
		one.EPS = v.BASICEPS.Float64
		one.FnEp = v.FINEXPE.Float64
		one.IDEp = v.POLIDIVIEXPE.Float64
		one.InRe = v.INTEINCO.Float64
		one.ItEp = v.INTEEXPE.Float64
		one.ITEp = v.INCOTAXEXPE.Float64
		one.MgEp = v.MANAEXPE.Float64
		one.MIIn = v.MINYSHARRIGH.Float64
		one.NOEp = v.NONOEXPE.Float64
		one.NORe = v.NONOREVE.Float64
		one.NtIn = v.NETPROFIT.Float64
		one.OATx = v.BIZTAX.Float64
		one.OCOR = v.BIZTOTCOST.Float64
		one.OpPr = v.PERPROFIT.Float64
		one.OpRe = v.BIZINCO.Float64
		one.SaEp = v.SALESEXPE.Float64
		one.TOpR = v.BIZTOTINCO.Float64
		one.ToPr = v.TOTPROFIT.Float64

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

type ProfitsInfo struct {
}

func NewProfitsInfo() *ProfitsInfo {
	return &ProfitsInfo{}
}

func (this *ProfitsInfo) GetJson(req *RequestParam) (*ResponseFinAnaJson, error) {
	logging.Info("GetJson %T, RequestParam: %+v", *this, *req)

	sc := NewSymbolToCompcode()
	if err := sc.getCompcode(req.SCode); err != nil {
		return nil, err
	}

	if sc.COMPCODE.Valid == false {
		logging.Error("finchina db: select COMPCODE from %s where SYMBOL='%s'", TABLE_TQ_OA_STCODE, req.SCode)
		return nil, ErrNullComp
	}

	sli := NewProfitsGeneral()
	list, err := sli.getProfitsJsonList(sc.COMPCODE.String, req)
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

// TQ_FIN_PROBINCSTATEMENTNEW    中文名称：银行利润表(新准则产品表)
// TQ_FIN_PROIINCSTATEMENTNEW    中文名称：保险利润表(新准则产品表)
// TQ_FIN_PROSINCSTATEMENTNEW    中文名称：证券利润表(新准则产品表)
