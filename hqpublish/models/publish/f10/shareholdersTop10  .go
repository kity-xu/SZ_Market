package f10

import (
	"strconv"

	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/logging"
)

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
	Name     string  `json:"name"`     // 股东名称
	Holdings float64 `json:"holdings"` // 持股数量
	Rate     float64 `json:"rate"`     // 占比
	Change   float64 `json:"change"`   // 变动
}

// 获取十大股东信息
func GetHN_F10_ShareholdersTop10(scode string, limit int32, htype int32, enddate string) (*Date, error) {
	var date Date
	scd, _ := strconv.Atoi(scode)
	date.Sid = int32(scd)
	date.Htype = htype

	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(scode); err != nil {
		return nil, err
	}

	if limit < 10 {
		limit = 10
	}

	// 判断是流通股东还是正常股东
	if htype == 1 { // 1、股东
		// 十大股东发布日期
		rdate, err := finchina.NewTQ_SK_SHAREHOLDER().GetSharEndDate(sc.COMPCODE.String)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		var nd []*NEndDate
		for _, v := range rdate {
			var d NEndDate
			d.Date = v.ENDDATE
			nd = append(nd, &d)
		}
		// 十大股东信息
		ldate, err := finchina.NewTQ_SK_SHAREHOLDER().GetSharBaseL(sc.COMPCODE.String, limit, enddate)
		if err != nil {
			logging.Error("%v", err)
			return nil, err
		}
		var hd []*Holders
		for _, v := range ldate {
			var h Holders
			h.Name = v.SHHOLDERNAME
			h.Holdings = v.HOLDERAMT
			h.Rate = v.HOLDERRTO
			h.Change = v.CURCHG.Float64
			hd = append(hd, &h)
		}
		date.Num = int32(len(hd))
		date.Ndate = nd
		date.HoldersL = hd
		return &date, err
	} //else if htype == 2 { // 2、流通股东
	// 查询日期列表
	rdate, err := finchina.NewTQ_SK_OTSHOLDER().GetOtshEndDate(sc.COMPCODE.String)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	var nd []*NEndDate
	for _, v := range rdate {
		var d NEndDate
		d.Date = v.ENDDATE
		nd = append(nd, &d)
	}
	// 查询股东信息
	ldate, err := finchina.NewTQ_SK_OTSHOLDER().GetOtshTop10L(enddate, sc.COMPCODE.String, limit)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}
	var hd []*Holders
	for _, v := range ldate {
		var h Holders
		h.Name = v.SHHOLDERNAME
		h.Holdings = v.HOLDERAMT
		h.Rate = v.PCTOFFLOTSHARES.Float64
		h.Change = v.HOLDERSUMCHGRATE.Float64
		hd = append(hd, &h)
	}
	date.Num = int32(len(hd))
	date.Ndate = nd
	date.HoldersL = hd
	return &date, err
	//}

}
