package f10

import (
	"strconv"

	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/logging"
)

type CSDate struct {
	Sid      int32           `json:"sid"`      // 证券ID
	Num      int32           `json:"num"`      // 请求条数(默认10条)
	Capstock []*CapitalStock `json:"capstock"` //
}
type CapitalStock struct {
	CTime      string  `json:"cTime"`      // 日期
	TotalShare float64 `json:"totalShare"` // 总股本（单位:股）
	CircskAmt  float64 `json:"circskAmt"`  // 流通股本（单位:股）
	Cause      string  `json:"cause"`      // 变动原因
}

type HN_F10_CapitalStock struct {
}

func NewHN_F10_CapitalStock() *HN_F10_CapitalStock {
	return &HN_F10_CapitalStock{}
}

// 获取股本变动信息
func GetF10CapitalStock(scode string, limit int) (*CSDate, error) {

	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(scode); err != nil {
		return nil, err
	}

	// 查询股本变动列表
	date, err := finchina.NewEquity().GetShareStruchg(sc.COMPCODE.String, limit)
	if err != nil {
		return nil, err
	}

	var csk []*CapitalStock
	for _, v := range date {
		var cs CapitalStock
		cs.CTime = v.BEGINDATE.String
		cs.TotalShare = v.TOTALSHARE.Float64
		cs.CircskAmt = v.CIRCSKAMT.Float64
		cs.Cause = v.SHCHGRSN.String
		csk = append(csk, &cs)
	}
	sd, err := strconv.Atoi(scode)
	if err != nil {
		logging.Error("%v", err)
	}
	var cs CSDate
	cs.Sid = int32(sd)
	cs.Num = int32(len(csk))
	cs.Capstock = csk
	return &cs, err
}
