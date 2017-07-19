package company

import (
	"haina.com/market/finance/models/finchina"
	"haina.com/share/logging"
)

type F10MobileTerminal struct {
	CompInfo F10_Compinfo           `json:"compinfo"`
	Equity   F10_Equity_Shareholder `json:"equity"`
	Dividend F10_Dividend_Ro        `json:"dividend"`
	Finance  F10_Finance            `json:"finance"`
}

//1.公司资料
type F10_Compinfo struct {
	Indus string `json:"Indus"` //公司所属证监会行业（聚源）
	Prov  string `json:"Prov"`  //省份
	Main  string `json:"Main"`  //经营范围-主营
}

//2.股本股东
type F10_Equity_Shareholder struct {
	Totalshare       float64 `json:"TotalShare"`       //总股本(万股)          ///TQ_SK_SHARESTRUCHG
	Circskamt        float64 `json:"Circskamt"`        //流通股本
	Totalshamt       float64 `json:"Totalshamt"`       //股东总户数            ///TQ_SK_SHAREHOLDERNUM
	Totalshrto       float64 `json:"Totalshrto"`       //股东总户数较上期增减
	Top1sha          string  `json:"No1share"`         //第一大股东           ///TQ_SK_SHAREHOLDER
	Top10Rate        float64 `json:"Top10rate"`        //前十大股东占比
	LegalPersonsRate float64 `json:"LegalPersonsRate"` //法人所占比例         ///TQ_SK_SHAREHOLDERNUM
}

//3.分红配股
type F10_Dividend_Ro struct {
	List []DividendRo `json:"list"`
}

type DividendRo struct {
	Date     string  `json:"Date"`     //年度
	Dividend float64 `json:"Dividend"` //分红（元，税前）
	RegDate  string  `json:"RegDate"`  //股权登记日
}

//4.财务数据
type F10_Finance struct {
	MainIncome  float64 `json:"mainIncome"`  //主营业务收入        ///TQ_FIN_PROINCSTATEMENTNEW
	MainBizRate float64 `json:"mainBizRate"` //主营收入同比增长率
	EPS         float64 `json:"EPS"`         //每股收益
	NetProfit   float64 `json:"netProfit"`   //净利润

	CapReserve float64 `json:"capReserve"` //每股公积金
	NetYield   float64 `json:"netYield"`   //净资产收益率_平均    ///TQ_FIN_PROINDICDATA
	Ratio      float64 `json:"ratio"`      //资产负债率
	UPPS       float64 `json:"UPPS"`       //每股未分配利润
}

func F10Mobile(scode string, market string) (*F10MobileTerminal, *string, error) {
	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(scode, market); err != nil {
		return nil, nil, err
	}
	/*-------------------------------------------------------------------*/
	/*----------------------------公司信息--------------------------------*/
	comp := finchina.NewCompInfo()
	cinfo, err := comp.GetCompInfo(sc.COMPCODE.String)
	industry, err := comp.GetCompTrade(sc.COMPCODE.String)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, nil, err
	}

	t1 := F10_Compinfo{
		Indus: industry,
		Prov:  getProvince(cinfo.REGION.String),
		Main:  cinfo.MAJORBIZ.String,
	}

	/*-------------------------------------------------------------------*/
	/*-------------------------股本股东-----------------------------------*/
	equity, err := finchina.NewEquity().GetEquity(sc.COMPCODE.String)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, nil, err
	}

	shnum, err := finchina.NewShareHolders().GetShareHolders(sc.COMPCODE.String)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, nil, err
	}

	top10, err := finchina.NewShareHoldersTop10().GetShareHoldersTop10(sc.COMPCODE.String)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, nil, err
	}

	var top10rate float64 = 0.0
	var nametop1 string

	for _, v := range *top10 {
		top10rate += v.HOLDERRTO.Float64
		if v.RANK.Int64 == int64(1) {
			nametop1 = v.SHHOLDERNAME.String
		}
	}

	t2 := F10_Equity_Shareholder{
		Totalshare: equity.TOTALSHARE.Float64,
		Circskamt:  equity.CIRCSKAMT.Float64,
		Totalshamt: shnum.TOTALSHAMT.Float64,
		Totalshrto: shnum.TOTALSHRTO.Float64,
		Top1sha:    nametop1,
		Top10Rate:  top10rate,
		//LegalPersonsRate: shnum.CORPSHAMT / float32(shnum.TOTALSHAMT*1.0),
	}

	/*-------------------------------------------------------------------*/
	/*--------------------------分红配股----------------------------------*/
	divs, err := finchina.NewDividendRO().GetDividendRO(sc.COMPCODE.String)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, nil, err
	}

	t3 := F10_Dividend_Ro{}
	for _, v := range *divs {
		div := DividendRo{
			Date:     v.DIVIYEAR.String,
			Dividend: v.PRETAXCASHMAXDVCNY.Float64,
			RegDate:  v.EQURECORDDATE.String,
		}
		t3.List = append(t3.List, div)
	}

	/*-------------------------------------------------------------------*/
	/*--------------------------财务数据----------------------------------*/
	f1, err := finchina.NewF10_MB_PROINCSTATEMENTNEW().GetF10_MB_PROINCSTATEMENTNEW(sc.COMPCODE.String)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, nil, err
	}

	f2, err := finchina.NewF10_MB_PROINDICDATA().GetF10_MB_PROINDICDATA(sc.COMPCODE.String)
	if err != nil {
		logging.Debug("%v", err.Error())
		return nil, nil, err
	}

	t4 := F10_Finance{
		MainIncome:  f1.MAINBIZINCO.Float64,
		MainBizRate: 0,
		EPS:         f1.BASICEPS.Float64,
		NetProfit:   f1.NETPROFIT.Float64,
		CapReserve:  f2.CRPS.Float64,
		NetYield:    f2.ROEAVG.Float64,
		Ratio:       f2.ASSLIABRT.Float64,
		UPPS:        f2.UPPS.Float64,
	}

	f10 := &F10MobileTerminal{
		CompInfo: t1,
		Equity:   t2,
		Dividend: t3,
		Finance:  t4,
	}
	return f10, &cinfo.COMPSNAME.String, nil
}
