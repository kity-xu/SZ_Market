//增发（Seasoned Equity Offerings）
package company

import (
	"haina.com/market/finance/models/finchina"
)

//Seasoned Equity Offerings
type SEO struct {
	SeoRaise float64
	SeoSc    int64
	SEOcount int

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

func (this *SEO) GetSEOList(scode string) (*[]*SEO, error) {
	list := make([]*SEO, 0)
	seos, err := new(finchina.TQ_SK_PROADDISS).GetSEOListFromFC(scode)
	if err != nil {
		return &list, err
	}
	list = this.newSEOListjson(seos)
	return &list, err
}

func (this *SEO) newSEOListjson(seos []finchina.TQ_SK_PROADDISS) []*SEO {
	list := make([]*SEO, 0)

	for _, v := range seos {
		var js SEO
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

		this.SeoSc += v.ISFINSUC.Int64
		this.SEOcount++

		list = append(list, &js)
	}
	return list
}
