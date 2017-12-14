package publish2

import (
	"liveshow/share/logging"
	"strconv"

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
	EQURECORDDATE      string  `json:"regDate"`   //股权登记日
}

func NewDividendJson() *DividendJson {
	return &DividendJson{}
}

func (DividendJson) GetDividendJson(sid int32) (*ResDiv, error) {
	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(strconv.Itoa(int(sid))); err != nil {
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
			EQURECORDDATE:      v.EQURECORDDATE.String,
		}
		dividends = append(dividends, div)
	}

	res := &ResDiv{
		Num:  len(dividends),
		Divs: dividends,
	}
	return res, nil
}
