package company

import (
	"haina.com/market/finance/models/finchina"
)

//Dividend
type Dividend struct {
	ToCash   float64
	Divcount int

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

func (this *Dividend) GetDividendList(sets uint64, scode string) (*[]*Dividend, error) {
	list := make([]*Dividend, 0)
	divs, err := new(finchina.TQ_SK_DIVIDENTS).GetDivListFromFC(sets, scode)
	if err != nil {
		return &list, err
	}
	list = this.newDivListjson(divs)
	return &list, err
}

func (this *Dividend) newDivListjson(divs []finchina.TQ_SK_DIVIDENTS) []*Dividend {
	list := make([]*Dividend, 0)
	for _, v := range divs {
		var js Dividend
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
		js.LisDate = v.LISTDATE.String

		this.Divcount++
		list = append(list, &js)
	}
	return list
}
