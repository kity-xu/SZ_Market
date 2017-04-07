//配股（Rights Offering）
package company

import (
	"haina.com/market/finance/models/finchina"
)

//Rights Offering
type RO struct {
	RoRaise float64
	RoSc    int64
	ROcount int

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

func (this *RO) GetROListJson(scode string) (*[]*RO, error) {
	list := make([]*RO, 0)
	ros, err := new(finchina.TQ_SK_PROPLACING).GetROList(scode)
	if err != nil {
		return &list, err
	}
	list = this.getROListjson(ros)
	return &list, err
}

func (this *RO) getROListjson(ros []finchina.TQ_SK_PROPLACING) []*RO {
	list := make([]*RO, 0)

	for _, v := range ros {
		var js RO
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

		this.RoSc += v.ISFINSUC.Int64
		this.ROcount++

		list = append(list, &js)
	}
	return list
}
