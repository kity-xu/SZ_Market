package company

import (
	"haina.com/market/finance/models/finchina"
	//"haina.com/share/logging"
)

type Top10 struct {
	ID    int64   `json:"-"`     // ID
	ISHIS int     `json:"ISHIS"` // 是否上一报告期存在股东
	Name  string  `json:"Name"`  // 股东名称
	Ovwet float64 `json:"Ovwet"` // 增持股份
	Posi  float64 `json:"Posi"`  // 持股数
	Hshr  float64 `json:"Hshr"`  // 持股数量增减幅度
	Prop  float64 `json:"Prop"`  // 持股数量占总股本比例
	Ptos  float64 `json:"Ptos"`  // 持股数量占流通A股比例
}

type TLSG interface{}
type RetTopInfoJson struct {
	CR float64 `json:"CR"`
	//EndDate string      `json:"EndDate"`
	Rate  float64     `json:"Rate"`
	SCode string      `json:"scode"`
	Sum   float64     `json:"Sum"`
	Edate string      `json:"Edate"`
	TLSG  interface{} `json:"TLSG"`
}

/**
  获取十大流通股东信息
*/
func GetTop10Group(enddate string, scode string, limit int, market string) ([]*RetTopInfoJson, error) {

	// 根据证券代码 开始时间 查询条数 查询日期
	var selwhe = " "
	if enddate != "" {
		selwhe = " and ENDDATE < '" + enddate + "' "
	}
	dataEnd, err := finchina.NewTQ_SK_OTSHOLDER().GetEndDate(scode, selwhe, limit+1, market)

	var whel = ""
	for _, item := range dataEnd {
		whel += "'" + item.ENDDATE + "',"
	}
	if len(whel) > 1 {
		whel = whel[0 : len(whel)-1]
	}

	// 查询9组十大流通股东
	data, err := finchina.NewTQ_SK_OTSHOLDER().GetTop10Group(whel, scode, market)
	// 外循环9个日期 内循环9组十大流通股东数据
	rtpj := []*RetTopInfoJson{}
	for index, itm := range dataEnd {
		var rij RetTopInfoJson
		var sumv1 = 0.0 // 累计
		var crv1 = 0.0  // 较上期变化
		var rat1 = 0.0  // 累计占路通股本比
		objz := []*Top10{}
		for _, item := range data {
			// 上期
			if index < len(dataEnd)-1 {
				if dataEnd[index+1].ENDDATE == item.ENDDATE {
					crv1 += item.HOLDERAMT
				}
			}
			// 最新一期
			if itm.ENDDATE == item.ENDDATE {
				sumv1 += item.HOLDERAMT
				rat1 += item.PCTOFFLOATSHARES
				// ----------
				var jsn Top10
				jsn.ISHIS = item.ISHIS
				jsn.Name = item.SHHOLDERNAME
				jsn.Ovwet = item.HOLDERSUMCHG.Float64
				jsn.Hshr = item.HOLDERSUMCHGRATE.Float64
				jsn.Posi = item.HOLDERAMT
				jsn.Ptos = item.PCTOFFLOATSHARES
				jsn.Prop = item.HOLDERRTO
				objz = append(objz, &jsn)
			}
		}
		crv1 = sumv1 - crv1 // 较上期变化
		// 保存一组数据
		rij.CR = crv1
		rij.Sum = sumv1
		rij.Rate = rat1
		rij.SCode = scode
		rij.Edate = itm.ENDDATE
		rij.TLSG = objz
		rtpj = append(rtpj, &rij)
	}

	return rtpj, err
}
