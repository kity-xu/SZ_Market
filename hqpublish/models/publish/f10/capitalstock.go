package f10

import (
	"encoding/json"
	"fmt"

	. "haina.com/market/hqpublish/models"

	"haina.com/market/hqpublish/models/finchina"
	"haina.com/share/logging"
)

type CSDate struct {
	Sid      int             `json:"sid"`      // 证券ID
	Num      int             `json:"num"`      // 请求条数(默认10条)
	Capstock []*CapitalStock `json:"capstock"` //
}
type CapitalStock struct {
	CTime      string  `json:"cTime"`      // 日期
	TotalShare float64 `json:"totalShare"` // 总股本（单位:股）
	CircskAmt  float64 `json:"circskAmt"`  // 流通股本（单位:股）
	Cause      string  `json:"cause"`      // 变动原因
}

// 获取股本变动信息
func GetF10CapitalStock(scode int, limit int) (*CSDate, error) {
	var cs CSDate

	key := fmt.Sprintf(REDIS_F10_CAPITALSTOCK, scode)
	data, err := RedisCache.GetBytes(key)
	if err == nil {
		if err = json.Unmarshal(data, &cs); err == nil {
			return &cs, nil
		}
		logging.Debug("股本变动: GetCache error |%v", err)
	}

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

	cs.Sid = scode
	cs.Num = len(csk)
	cs.Capstock = csk

	bys, err := json.Marshal(&cs)
	if err != nil {
		logging.Debug("股本变动: SetCache error")
		return &cs, nil
	}
	RedisCache.Setex(key, TTL.F10HomePage, bys)

	return &cs, nil
}
