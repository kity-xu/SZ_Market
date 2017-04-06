package company

import (
	"errors"
	"strconv"

	"haina.com/market/finance/models/finchina"
)

type Top10 struct {
	ID    int64  `json:"-"`     // ID
	ISHIS int    `json:"ISHIS"` // 是否上一报告期存在股东
	Name  string `json:"Name"`  // 股东名称
	Ovwet string `json:"Ovwet"` // 增持股份
	Posi  string `json:"Posi"`  // 持股数
	Prop  string `json:"Prop"`  // 持股数量占总股本比例
}

type TopList interface{}
type RetTopInfoJson struct {
	CR      float64     `json:"CR"`
	EndDate string      `json:"EndDate"`
	Rate    string      `json:"Rate"`
	SCode   string      `json:"scode"`
	Sum     string      `json:"Sum"`
	TopList interface{} `json:"TLSG"`
}

/**
  获取十大流通股东信息
*/
func GetTop10Group(enddate string, scode string, limit int) (RetTopInfoJson, error) {

	dataEnd, err := finchina.NewTQ_SK_OTSHOLDER().GetEndDate(scode)
	var endStr = ""
	for _, item := range dataEnd {
		endStr += item.ENDDATE + ","
	}
	var zendtata = ""
	if enddate != "" {
		zendtata = enddate
	} else {
		zendtata = dataEnd[0].ENDDATE
	}

	data, err := finchina.NewTQ_SK_OTSHOLDER().GetTop10Group(zendtata, scode, limit)

	var rij RetTopInfoJson
	jsns10 := []*Top10{}

	for _, item := range data {

		jsn, err := GetTop10Info(item)
		if err != nil {
			return rij, err
		}

		jsns10 = append(jsns10, jsn)
	}

	tp := finchina.NewCalculate().GetSingleCalculate(zendtata, scode)
	//计算上次累计持股
	var sSum float64
	for index, item := range dataEnd {
		if zendtata == item.ENDDATE {
			if index < len(dataEnd)-1 {
				tp1 := finchina.NewCalculate().GetSingleCalculate(dataEnd[index+1].ENDDATE, scode)
				sSum, err = strconv.ParseFloat(tp1.Sumh, 64)
				break
			}
		}

	}

	rij.TopList = jsns10
	rij.SCode = scode
	rij.Sum = tp.Sumh
	rij.Rate = tp.Rate
	var Sums float64
	Sums, err = strconv.ParseFloat(tp.Sumh, 64)
	rij.CR = Sums - sSum

	var enddatestr = ""
	if len(endStr) > 1 {
		enddatestr = endStr[0 : len(endStr)-1]
	}
	rij.EndDate = enddatestr
	return rij, err
}

// 获取Top10信息
func GetTop10Info(tso *finchina.TQ_SK_OTSHOLDER) (*Top10, error) {
	var jsn Top10
	if len(tso.SHHOLDERNAME) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &Top10{
		//ID:    top10.ID,
		ISHIS: tso.ISHIS,
		Name:  tso.SHHOLDERNAME,
		Ovwet: tso.HOLDERSUMCHG.String,
		Posi:  tso.HOLDERAMT,
		Prop:  tso.HOLDERRTO,
	}, nil
}
