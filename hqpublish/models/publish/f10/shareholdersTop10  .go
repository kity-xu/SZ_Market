package f10

import "haina.com/market/hqpublish/models/finchina"

type Date struct {
	Sid      int32       `json:"sid"`     // 证券ID
	Num      int32       `json:"num"`     // 条数
	Htype    int32       `json:"htype"`   // 1：股东； 2：流通股东
	Ndate    []*NEndDate `json:"ndate"`   // 日期数组
	HoldersL []*Holders  `json:"holders"` // 股东信息
}

// 统计日期
type NEndDate struct {
	Date string `json:"date"`
}

// 十大股东信息
type Holders struct {
	Name      string  `json:"name"`      // 股东名称
	Holdings  float64 `json:"holdings"`  // 持股数量
	CircskAmt float64 `json:"circskAmt"` // 流通股本（单位:股）
	Cause     string  `json:"cause"`     // 变动原因
}

type HN_F10_ShareholdersTop10 struct {
}

func NewHN_F10_ShareholdersTop10() *HN_F10_ShareholdersTop10 {
	return &HN_F10_ShareholdersTop10{}
}

// 获取十大股东信息
func GetHN_F10_ShareholdersTop10(scode string, limit int, htype int) (*Date, error) {

	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(scode); err != nil {
		return nil, err
	}

	return nil, nil
	//
	//	// 查询股本变动列表
	//	date, err := finchina.NewEquity().GetShareStruchg(sc.COMPCODE.String, limit)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	var csk []*CapitalStock
	//	for _, v := range date {
	//		var cs CapitalStock
	//		cs.CTime = v.BEGINDATE.String
	//		cs.TotalShare = v.TOTALSHARE.Float64
	//		cs.CircskAmt = v.CIRCSKAMT.Float64
	//		cs.Cause = v.SHCHGRSN.String
	//		csk = append(csk, &cs)
	//	}
	//	sd, err := strconv.Atoi(scode)
	//	if err != nil {
	//		logging.Error("%v", err)
	//	}
	//	var cs CSDate
	//	cs.Sid = int32(sd)
	//	cs.Num = int32(len(csk))
	//	cs.Capstock = csk
	//	return &cs, err
}
