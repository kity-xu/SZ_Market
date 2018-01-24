package bigService

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"haina.com/market/f9/models/common"
	"haina.com/share/logging"
)

var (
	symbolParam   string //证券代码（带字母）
	symbol        string //证券代码（不带字母)
	symboType     string //证券类型
	swlevelcode   string //行业代码
	swlevelname   string //行业名称
	compcode      string //公司代码
	compname      string //公司名称
	companyDetail string //公司详情
	//sumcompay     string //该行业的所有公司
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

func GetBigData(scode string) (bigData, error) {
	symbol = scode[2:]
	symbolParam = scode
	if scode[0:2] == "sz" {
		symboType = "001003"
	} else if scode[0:2] == "sh" {
		symboType = "001002"
	}

	var bi bigData

	detail, err := companyDetailModel.NewCompanyDetail().GetCompanyDetail(symbol, symboType)

	//logging.Info("detail==".)

	if err != nil {
		logging.Info("该股票不存在")
		bi.Code = 20000
		bi.Desc = "该股票不存在"
		return bi, err
	}

	if detail.LISTSTATUS != 1 {
		logging.Info("非股票代码")
		bi.Code = 20000
		bi.Desc = "非股票代码"
		return bi, err
	}

	data, _ := GetApi(symbolParam)
	logging.Info("dataApi===%+v", data)
	bi.Result = data.Result.Data.Res.Result
	bi.Market_care = data.Result.Data.Res.Market_care
	bi.Market_emotion = data.Result.Data.Res.Market_emotion
	bi.Emotion_num = data.Result.Data.Res.Emotion_num
	bi.Theme_chart = data.Result.Data.Res.Theme_chart
	bi.Theme = data.Result.Data.Res.Theme
	bi.Recommend_big = data.Result.Data.Res.Recommend_big
	return bi, nil
}

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
	logging.Info(string(body))

	var _param result
	err = json.Unmarshal([]byte(body), &_param)
	logging.Info("===%+v", _param)
	return _param, err
}
