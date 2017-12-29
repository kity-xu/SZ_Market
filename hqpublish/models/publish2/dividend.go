package publish2

import (
	"fmt"

	. "haina.com/market/hqpublish/models"
	"haina.com/share/logging"

	"haina.com/market/hqpublish/models/finchina"
)

type ResDiv struct {
	Num  int             `json:"num"`
	Divs []*DividendJson `json:"divs"`
}

type DividendJson struct {
	DIVIYEAR           string  `json:"date"`      //年度
	PRETAXCASHMAXDVCNY float64 `json:"dividend"`  //分红
	PROBONUSRT         float64 `json:"pro"`       //送股比例(10:X)
	TRANADDRT          float64 `json:"tranAddRt"` //转增比例(10:X)
	BONUSRT            float64 `json:"bonusRt"`   //赠股比例(10:X)
	EQURECORDDATE      string  `json:"regDate"`   //股权登记日 -- > 除权除息日
}

func NewDividendJson() *DividendJson {
	return &DividendJson{}
}

func (DividendJson) GetDividendJson(sid int32) (*ResDiv, error) {
	res := &ResDiv{}
	key := fmt.Sprintf(REDIS_CACHE_DIVIDEND_K, sid)
	if _, err := GetResFromCache(key, res); err == nil {
		return res, nil
	}

	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(sid); err != nil {
		return nil, err
	}
	divs, err := finchina.NewDividendRO().GetDividendRO(sc.COMPCODE.String)
	if err != nil {
		logging.Error("%v", err)
		return nil, err
	}

	dividends := make([]*DividendJson, 0, 4)
	for _, v := range *divs {
		div := &DividendJson{
			DIVIYEAR:           v.DIVIYEAR.String,
			PRETAXCASHMAXDVCNY: v.PRETAXCASHMAXDVCNY.Float64,
			PROBONUSRT:         v.PROBONUSRT.Float64,
			TRANADDRT:          v.TRANADDRT.Float64,
			BONUSRT:            v.BONUSRT.Float64,
			EQURECORDDATE:      v.XDRDATE.String, //分红配股日
		}
		dividends = append(dividends, div)
	}

	res.Num = len(dividends)
	res.Divs = dividends
	SetResToCache(key, res)

	return res, nil
}
