//增发（Seasoned Equity Offerings）
package company

import (
	"haina.com/market/finance/models/finchina"
)

type SEO struct {
	Scode   string  `json:"scode"`
	Cash    float64 `json:"tocash"`
	Count   int     `json:"count"`
	Success int     `json:"success"`
	Faild   int     `json:"faild"`
	Ing     int     `json:"ing"`
	List    interface{}
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

func (this *SEO) GetSEOList(scode string, market string) (*SEO, error) {
	var list SEO
	seos, err := new(finchina.TQ_SK_PROADDISS).GetSEOListFromFC(scode, market)
	if err != nil {
		return &list, err
	}
	list = this.newSEOListjson(seos)
	return &list, err
}

func (this *SEO) newSEOListjson(seos []finchina.TQ_SK_PROADDISS) SEO {
	var seo SEO
	list := make([]SEOJson, 0)

	var cash float64
	//	var sc int
	var count int

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
		cash += v.ACTNETRAISEAMT.Float64
		//sc += int(v.ISFINSUC.Int64)
		count++
		if "1" == v.ISSUESTATUS.String {
			seo.Ing += 1
		} else if "7" == v.ISSUESTATUS.String || "8" == v.ISSUESTATUS.String {
			seo.Success += 1
		}

		list = append(list, js)
	}
	seo.Count = count
	seo.Cash = cash
	seo.Faild = seo.Count - seo.Success - seo.Ing
	seo.List = list
	return seo
}
