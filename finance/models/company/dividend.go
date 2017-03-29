package company

import (
	"strings"

	"haina.com/market/finance/models/finchina"
)

type FinDivid struct {
	ToCash   float64 //分红金额合计
	SeoRaise float64 //增发实际募资净额合计
	RoRaise  float64 //配股实际募资金额合计
}

//Dividend
type DivJson struct {
	Bene     string  `json:"Bene"`     //分红对象
	Bonus    float64 `json:"Bonus"`    //送股（股）
	Date     string  `json:"Data"`     //年度
	Dividend float64 `json:"Dividend"` //分红（元，税前）
	DivDate  string  `json:"DivDate"`  //红利发放日
	DivRate  string  `json:"DivRate"`  //股利支付率（%）
	DNDate   string  `json:"DNDate"`   //决案公布日
	Evolve   string  `json:"Evolve"`   //事情进展
	ExDate   string  `json:"ExDate"`   //除权除息日
	INDate   string  `json:"INDate"`   //实施公告日
	LisDate  string  `json:"LisDate"`  //转股上市日
	PNDate   string  `json:"PNDate"`   //预案公布日
	RegDate  string  `json:"RegDate"`  //股权登记日
	Tran     float64 `json:"Tran"`     //转股（股）
}

//Seasoned Equity Offerings
type SEOJson struct {
	AGMD    string  `json:"AGMD"`    //股东大会决议公告日
	IECD    string  `json:"IECD"`    //发审委公告日
	LisDate string  `json:"LisDate"` //新股上市日
	PNDate  string  `json:"PNDate"`  //预案公布日
	PPrice  float64 `json:"PPrice"`  //预案发行价格
	Price   float64 `json:"Price"`   //实际发行价格
	PVal    float64 `json:"PVal"`    //预案募资金额
	PVol    float64 `json:"PVol"`    //预案发行数量
	Range   string  `json:"Range"`   //发行对象类型
	SEOD    string  `json:"SEOD"`    //发行新股日
	SRCD    string  `json:"SRCD"`    //证监会核准公告日
	Step    string  `json:"Step"`    //事情进展
	Type    string  `json:"Type"`    //发行方式
	Val     float64 `json:"Val"`     //实际募资金额
	Vol     float64 `json:"Vol"`     //实际发行数量
}

//Rights Offering
type ROJson struct {
	AGMD    string  `json:"AGMD"`    //股东大会决议公告日
	Code    string  `json:"Code"`    //配股代码
	DNDate  string  `json:"DNDate"`  //决案公布日
	ERDate  string  `json:"ERDate"`  //配股除权日
	LisDate string  `json:"LisDate"` //配股上市日
	PNDate  string  `json:"PNDate"`  //预案公布日
	PProp   float64 `json:"PProp"`   //计划配股比例
	Price   float64 `json:"Price"`   //实际配股价格
	Prop    float64 `json:"Prop"`    //实际配股比例
	PVol    float64 `json:"PVol"`    //计划配股数量
	RegDate string  `json:"RegDate"` //股权登记日
	ROPD    string  `json:"ROPD"`    //配股缴款起止日
	Short   string  `json:"Short"`   //配股简称
	Vol     float64 `json:"Vol"`     //实际配股数量
}

// Repo(repurchase agreement)
type RepoJson struct {
}

func (this *FinDivid) GetDivListJson(sets uint64, secode string) (*[]*DivJson, error) {
	list := make([]*DivJson, 0)
	divs, err := new(finchina.Dividend).GetDivList(sets, secode)
	if err != nil {
		return &list, err
	}
	list = getDivListjson(this, divs)
	return &list, err
}
func (this *FinDivid) GetSEOListJson(secode string) (*[]*SEOJson, error) {
	list := make([]*SEOJson, 0)
	seos, err := new(finchina.Dividend).GetSEOList(secode)
	if err != nil {
		return &list, err
	}
	list = getSEOListjson(this, seos)
	return &list, err
}
func (this *FinDivid) GetROListJson(secode string) (*[]*ROJson, error) {
	list := make([]*ROJson, 0)
	ros, err := new(finchina.Dividend).GetROList(secode)
	if err != nil {
		return &list, err
	}
	list = getROListjson(this, ros)
	return &list, err
}

func getDivListjson(this *FinDivid, divs []finchina.Div) []*DivJson {
	list := make([]*DivJson, 0)
	for _, v := range divs {
		var js DivJson
		js.Bene = v.GRAOBJ.String
		js.Bonus = v.PROBONUSRT.Float64
		js.Date = v.DIVIYEAR.String
		js.Dividend = v.PRETAXCASHMAXDVCNY.Float64
		js.DivDate = v.CASHDVARRBEGDATE.String
		js.INDate = v.PUBLISHDATE.String
		js.DNDate = v.SHHDMEETRESPUBDATE.String
		js.ExDate = v.XDRDATE.String
		js.RegDate = v.EQURECORDDATE.String
		js.Tran = v.TRANADDRT.Float64
		this.ToCash += v.TOTCASHDV.Float64
		if strings.EqualFold(v.LISTDATE.String, "19000101") {
			js.LisDate = "--"
		} else {
			js.LisDate = v.LISTDATE.String
		}

		list = append(list, &js)
	}
	return list
}

func getSEOListjson(this *FinDivid, seos []finchina.SEO) []*SEOJson {
	list := make([]*SEOJson, 0)

	for _, v := range seos {
		var js SEOJson
		js.IECD = v.CSRCAPPRAGREEDATE.String
		js.LisDate = v.UPDATEDATE.String
		//js.PNDate=  v.ADDISSPUBDATE.String
		js.PPrice = v.ENQUMAXPRICE.Float64
		js.Price = v.ISSPRICE.Float64
		js.PVal = v.PLANTOTRAISEAMT.Float64
		js.PVol = v.PLANISSMAXQTY.Float64
		js.Range = v.ISSUEOBJECT.String
		js.SEOD = v.SHARECORDDATE.String
		js.SRCD = v.CSRCAPPDPUBDATE.String
		js.Val = v.ACTNETRAISEAMT.Float64
		js.Vol = v.ACTISSQTY.Float64
		js.Type = v.ISSUEMODEMEMO.String
		this.SeoRaise += v.ACTNETRAISEAMT.Float64

		list = append(list, &js)
	}
	return list
}

func getROListjson(this *FinDivid, ros []finchina.RO) []*ROJson {
	list := make([]*ROJson, 0)

	for _, v := range ros {
		var js ROJson
		js.AGMD = v.UPDATEDATE.String
		js.Code = v.ALLOTCODE.String
		js.DNDate = v.LISTPUBDATE.String
		js.ERDate = v.EXRIGHTDATE.String
		js.LisDate = v.LISTDATE.String
		js.PNDate = v.PUBLISHDATE.String
		js.PProp = v.ACTTOTALLOTRT.Float64
		js.Price = v.ALLOTPRICE.Float64
		js.Prop = v.ALLOTRT.Float64
		js.PVol = v.PLANISSMAXQTY.Float64
		js.RegDate = v.EQURECORDDATE.String
		js.ROPD = v.PAYBEGDATE.String + "~" + v.PAYENDDATE.String
		js.Short = v.ALLOTSNAME.String
		js.Vol = v.ACTISSQTY.Float64
		this.RoRaise += v.ACTNETRAISEAMT.Float64

		list = append(list, &js)
	}
	return list
}
