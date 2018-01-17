package f10

import (
	"encoding/json"
	"fmt"

	. "haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/logging"
)

type ResTop10 struct {
	Sid      int        `json:"sid"`     // 证券ID
	Num      int        `json:"num"`     // 条数
	Htype    int        `json:"htype"`   // 1：股东； 2：流通股东
	Ndate    []int      `json:"ndate"`   // 日期数组
	HoldersL []*Holders `json:"holders"` // 股东信息
}

// 十大股东信息
type Holders struct {
	Date     int     `json:"date"`     // 日期
	Name     string  `json:"name"`     // 股东名称
	Holdings float64 `json:"holdings"` // 持股数量
	Rate     float64 `json:"rate"`     // 占比
	Change   float64 `json:"change"`   // 变动
	IsHis    int     `json:"ishis"`    // 上一期股东是否存在
}

type ShareHolderTop10 struct{}

// 获取十大股东信息
func GetHN_F10_ShareholdersTop10(scode int, htype int, enddate int) (*ResTop10, error) {
	var top ResTop10
	key := fmt.Sprintf(REDIS_F10_SHAREHOLDERSTOP10, htype, scode, enddate)
	data, err := RedisCache.GetBytes(key)
	if err == nil {
		if err = json.Unmarshal(data, &top); err == nil {
			return &top, nil
		}
		logging.Debug("Top10: GetCache error |%v", err)
	}

	top.Sid = scode
	top.Htype = htype
	sc := finchina.NewTQ_OA_STCODE()
	if err = sc.GetCompcode(scode); err != nil {
		return nil, err
	}

	switch htype {
	case 1:
		err = new(ShareHolderTop10).Top10(&top, sc.COMPCODE.String, enddate)
	case 2:
		err = new(ShareHolderTop10).Top10Current(&top, sc.COMPCODE.String, enddate)
	default:
		logging.Error("invalid param type of 'htype'")
		err = fmt.Errorf("invalid param type of 'htype'")
	}
	if err == nil {
		bys, err := json.Marshal(&top)
		if err != nil {
			logging.Debug("Top10: SetCache error")
			return &top, nil
		}
		RedisCache.Setex(key, TTL.F10HomePage, bys)
	}

	return &top, err
}

// 十大股东
func (*ShareHolderTop10) Top10(top *ResTop10, compcode string, enddate int) error {
	// 十大股东发布日期
	times, err := finchina.NewTQ_SK_SHAREHOLDER().GetSharEndDate(compcode)
	if err != nil {
		logging.Error("%v", err)
		return err
	}

	// 十大股东信息
	ldate, err := finchina.NewTQ_SK_SHAREHOLDER().GetSharBaseL(compcode, 10, enddate)
	if err != nil {
		logging.Error("%v", err)
		return err
	}
	var hd []*Holders
	for _, v := range ldate {
		h := &Holders{
			Date:     v.ENDDATE,
			Name:     v.SHHOLDERNAME,
			Holdings: v.HOLDERAMT.Float64,
			Rate:     v.HOLDERRTO.Float64,
			Change:   v.CURCHG.Float64,
			IsHis:    v.ISHIS,
		}
		hd = append(hd, h)
	}
	top.Num = len(hd)
	top.Ndate = times
	top.HoldersL = hd

	return nil
}

// 十大流通股东
func (*ShareHolderTop10) Top10Current(top *ResTop10, compcode string, enddate int) error {
	// 查询日期列表
	times, err := finchina.NewTQ_SK_OTSHOLDER().GetOtshEndDate(compcode)
	if err != nil {
		logging.Error("%v", err)
		return err
	}

	// 查询股东信息
	ldate, err := finchina.NewTQ_SK_OTSHOLDER().GetOtshTop10L(compcode, 10, enddate)
	if err != nil {
		logging.Error("%v", err)
		return err
	}
	var hd []*Holders
	for _, v := range ldate {
		h := &Holders{
			Date:     v.ENDDATE,
			Name:     v.SHHOLDERNAME,
			Holdings: v.HOLDERAMT.Float64,
			Rate:     v.PCTOFFLOTSHARES.Float64,
			Change:   v.HOLDERSUMCHG.Float64,
			IsHis:    v.ISHIS,
		}
		hd = append(hd, h)
	}
	top.Num = len(hd)
	top.Ndate = times
	top.HoldersL = hd
	return nil
}
