package company

import (
	"errors"

	"haina.com/market/finance/models/finchina"
)

type ShareHolderJson struct {
	CRPS string `json:"CRPS"` // 户均持股较上期变化（%）
	TNS  string `json:"TNS"`  // 股东总户数（户）
	Date string `json:"Date"` // 指标\日期
	ANS  string `json:"ANS"`  // 户均持股数（股/户）
	APS  string `json:"APS"`  // 户均持股比例（%）
}

type SharList interface{}
type RetInfoJson struct {
	SCode    string      `json:"scode"`
	SharList interface{} `json:"Shareholders"`
}

/**
  获取股东人数信息
*/
func GetShareholderList(enddate string, sCode string, limit int) (RetInfoJson, error) {

	data, err := finchina.NewShareHolder().GetListByExps(enddate, sCode, limit)
	var js RetInfoJson
	jsns := []*ShareHolderJson{}

	for _, item := range data {
		jsn, err := GetJson(item)
		if err != nil {
			//	return jsns, err
		}

		jsns = append(jsns, jsn)

		js.SCode = sCode
		js.SharList = jsns

	}
	return js, err
}

// 获取JSON
func GetJson(shareHolder *finchina.ShareHolder) (*ShareHolderJson, error) {
	var jsn ShareHolderJson
	if len(shareHolder.ENDDATE) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &ShareHolderJson{
		Date: shareHolder.ENDDATE,
		TNS:  shareHolder.TOTALSHAMT,
		ANS:  shareHolder.KAVGSH,
		APS:  shareHolder.HOLDPROPORTIONPACC,
		CRPS: shareHolder.PROPORTIONCHG,
	}, nil
}
