package company

import (
	"errors"
	"strconv"

	"haina.com/market/finance/models/finchina"
)

type Top10Json struct {
	//SCode string `json:"SCode"` // 股票代码
	Name  string `json:"Name"`  // 股东名称
	Ovwet string `json:"Ovwet"` // 增持股份
	Posi  string `json:"Posi"`  // 持股数
	Prop  string `json:"Prop"`  // 持股数量占总股本比例
	ISHIS int    `json:"ISHIS"` // 是否上一报告期存在股东
}

type TopList interface{}
type RetTopInfoJson struct {
	SCode   string      `json:"scode"`
	Sum     string      `json:"Sum"`
	Rate    string      `json:"Rate"`
	CR      float64     `json:"CR"`
	EndDate string      `json:"EndDate"`
	TopList interface{} `json:"TLSG"`
}

/**
  获取十大流通股东信息
*/
func GetTop10List(enddate string, sCode string, limit int) (RetTopInfoJson, error) {

	dataEnd, err := finchina.NewTop10().GetEndDate(sCode)
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

	data, err := finchina.NewTop10().GetTop10List(zendtata, sCode, limit)

	var rij RetTopInfoJson
	jsns10 := []*Top10Json{}

	for _, item := range data {

		jsn, err := GetTop10Json(item)
		if err != nil {
			//return jsns, err
		}

		jsns10 = append(jsns10, jsn)
	}

	tp := finchina.NewCalculate().GetSingleByExps(zendtata, sCode)
	//计算上次累计持股
	var sSum float64
	for index, item := range dataEnd {
		if zendtata == item.ENDDATE {
			if index < len(dataEnd)-1 {
				tp1 := finchina.NewCalculate().GetSingleByExps(dataEnd[index+1].ENDDATE, sCode)
				sSum, err = strconv.ParseFloat(tp1.Sumh, 64)
				break
			}
		}

	}

	rij.TopList = jsns10
	rij.SCode = sCode
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

// 获取JSON
func GetTop10Json(top10 *finchina.Top10) (*Top10Json, error) {
	var jsn Top10Json
	if len(top10.SHHOLDERNAME) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &Top10Json{
		Name:  top10.SHHOLDERNAME,
		Ovwet: top10.HOLDERSUMCHG.String,
		Posi:  top10.HOLDERAMT,
		Prop:  top10.HOLDERRTO,
		ISHIS: top10.ISHIS,
	}, nil
}
