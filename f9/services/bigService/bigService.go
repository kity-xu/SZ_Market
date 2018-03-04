package bigService

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"haina.com/market/f9/models/finchina"
	"haina.com/share/logging"
)

type bigData struct {
	RetVal
	Result         string          `json:"result"`
	Market_care    string          `json:"market_care"`
	Market_emotion string          `json:"market_emotion"`
	Emotion_num    Emotion_num     `json:"emotion_num"`
	Theme_chart    []Theme_chart   `json:"theme_chart"`
	Theme          string          `json:"theme"`
	Recommend_big  []Recommend_big `json:"recommend_big"`
}
type RetVal struct {
	Code int    `json:"code,omitempty"`
	Desc string `json:"desc,omitempty"`
}

func GetBigData(sid string) (*bigData, error) {
	var bi bigData
	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(sid); err != nil {
		return nil, err
	}

	detail, err := finchina.NewCompanyDetail().GetCompanyDetail(sc.SECODE.String)
	if err != nil || detail.LISTSTATUS != 1 {
		return nil, err
	}

	var exchange string
	if detail.EXCHANGE == "001002" { //上海证券交易所
		exchange = "sh"
	} else if detail.EXCHANGE == "001003" { //深圳证券交易所
		exchange = "sz"
	} else if detail.EXCHANGE == "001004" { //股份转让市场
		exchange = "zr"
	} else {
		exchange = "no"
	}
	data, _ := GetApi(exchange + sid[2:])
	logging.Info("symbol exchange %v", exchange)
	bi.Result = data.Result.Data.Res.Result
	bi.Market_care = data.Result.Data.Res.Market_care
	bi.Market_emotion = data.Result.Data.Res.Market_emotion
	bi.Emotion_num = data.Result.Data.Res.Emotion_num
	bi.Theme_chart = data.Result.Data.Res.Theme_chart
	bi.Theme = data.Result.Data.Res.Theme
	bi.Recommend_big = data.Result.Data.Res.Recommend_big
	return &bi, nil
}

//--------------------------------------------------------sina api---------------------------------------------------------------//
type result struct {
	Result ResultData `json:"result"`
}
type ResultData struct {
	Status Status `json:"status"`
	Data   Data   `json:"data"`
}
type Status struct {
	Code int64 `json:"code"`
}
type Data struct {
	Code int64 `json:"code"`
	Res  Res   `json:"res"`
}
type Res struct {
	Result         string          `json:"result"`
	Market_care    string          `json:"market_care"`
	Market_emotion string          `json:"market_emotion"`
	Emotion_num    Emotion_num     `json:"emotion_num"`
	Theme_chart    []Theme_chart   `json:"theme_chart"`
	Theme          string          `json:"theme"`
	Recommend_big  []Recommend_big `json:"recommend_big"`
}
type Theme_chart struct {
	THEME_ID string `json:"THEME_ID"`
	Date     string `json:"date"`
	Data_bar string `json:"data_bar"`
}
type Emotion_num struct {
	Positive_num string `json:"positive_num"`
	Negative_num string `json:"negative_num"`
	Neutral_num  int64  `json:"neutral_num"`
}
type Recommend_big struct {
	Symbol      string `json:"symbol"`
	Symbol_name string `json:"symbol_name"`
}

func GetApi(symbol string) (result, error) {
	var url string = "http://touzi.sina.com.cn/api/openapi.php/PerspectiveSzBigService.bigPerspective?symbol=" + symbol
	resp, err := http.Get(url)
	if err != nil {
		logging.Error(err.Error())
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logging.Error(err.Error())
	}

	var _param result
	err = json.Unmarshal([]byte(body), &_param)
	return _param, err
}
