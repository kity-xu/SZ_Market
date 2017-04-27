package company

import (
	"errors"

	"haina.com/market/finance/models/finchina"
)

type ShareHolder struct {
	ID   int64  `json:"-"`    // ID
	ANS  string `json:"ANS"`  // 户均持股数（股/户）
	APS  string `json:"APS"`  // 户均持股比例（%）
	CRPS string `json:"CRPS"` // 户均持股较上期变化（%）
	Date string `json:"Date"` // 指标\日期
	TNS  string `json:"TNS"`  // 股东总户数（户）
}

type SharList interface{}
type RetInfoJson struct {
	SCode    string      `json:"scode"`
	SharList interface{} `json:"Shareholders"`
}

/**
  获取股东人数信息
*/
func GetShareholderGroup(scode string, limit int, strdate string) (RetInfoJson, error) {

	data, err := finchina.NewTQ_SK_SHAREHOLDERNUM().GetListByExps(scode, limit, strdate)
	var js RetInfoJson
	jsns := []*ShareHolder{}

	for _, item := range data {
		jsn, err := GetSharNum(item)
		if err != nil {
			return js, err
		}

		jsns = append(jsns, jsn)

		js.SCode = scode
		js.SharList = jsns

	}
	return js, err
}

// 获取JSON
func GetSharNum(shareHolder *finchina.TQ_SK_SHAREHOLDERNUM) (*ShareHolder, error) {
	var jsn ShareHolder
	if len(shareHolder.ENDDATE) < 1 {
		return &jsn, errors.New("obj is nil")
	}

	return &ShareHolder{
		ANS:  shareHolder.KAVGSH,
		APS:  shareHolder.HOLDPROPORTIONPACC,
		CRPS: shareHolder.PROPORTIONCHG,
		Date: shareHolder.ENDDATE,
		TNS:  shareHolder.TOTALSHAMT,
	}, nil
}
