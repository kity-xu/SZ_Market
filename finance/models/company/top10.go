package company

import (
	"errors"

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
	CR      string      `json:"CR"`
	TopList interface{} `json:"TLSG"`
}

/**
  获取十大流通股东信息
*/
func GetTop10List(enddate string, sCode string, limit int) (RetTopInfoJson, error) {

	data, err, scpcod := finchina.NewTop10().GetTop10List(enddate, sCode, limit)
	var rij RetTopInfoJson
	jsns10 := []*Top10Json{}

	for _, item := range data {

		jsn, err := GetTop10Json(item)
		if err != nil {
			//return jsns, err
		}

		jsns10 = append(jsns10, jsn)
	}

	tp := finchina.NewTop10().GetSingleByExps(enddate, scpcod)

	rij.TopList = jsns10
	rij.SCode = sCode
	rij.Sum = tp.Sumh
	rij.Rate = tp.Rate
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
