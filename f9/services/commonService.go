package commonService

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"haina.com/market/f9/models"
	"haina.com/market/f9/models/finchina"
	"haina.com/share/logging"
)

type SecBasicInfo struct {
	SID           string
	Status        int    //证券状态
	Secode        string //证券内码
	SymbolParam   string //证券代码（带字母）
	Symbol        string //证券代码（不带字母)
	ExChange      string //市场类型
	Swlevelcode   string //行业代码
	Swlevelname   string //行业名称
	Compcode      string //公司代码
	Compname      string //公司名称
	CompanyDetail string //公司详情

	//sumcompay     string //该行业的所有公司
}

const (
	CONST_REQUEST_OK                = 0
	CONST_REQUEST_ERROR_OUTTIME int = iota + 9999
	CONST_REQUEST_ERROR_SIGN
	CONST_REQUEST_ERROR_REDIS
	CONST_REQUEST_ERROR_PARAM
	CONST_REQUEST_ERROR_DATA
	CONST_REQUEST_ERROR_CONNET
	CONST_REQUEST_ERROR_NEWSEC
	CONST_REQUEST_ERROR_SUSPEND
	CONST_REQUEST_ERROR_SID
	CONST_REQUEST_ERROR_LISTCOMP
)

func GetCommonData(sid string) (*SecBasicInfo, error) {
	sc := finchina.NewTQ_OA_STCODE()
	if err := sc.GetCompcode(sid); err != nil {
		return nil, err
	}

	detail, err := finchina.NewCompanyDetail().GetCompanyDetail(sc.SECODE.String)
	if err != nil {
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

	basic := &SecBasicInfo{
		SID:         sid,
		Secode:      sc.SECODE.String,
		Status:      detail.LISTSTATUS,
		SymbolParam: exchange + detail.SYMBOL,
		Symbol:      detail.SYMBOL,
		ExChange:    detail.EXCHANGE,
		Swlevelcode: detail.SWLEVEL1CODE,
		Swlevelname: detail.SWLEVEL1NAME,
		Compcode:    sc.COMPCODE.String,
		Compname:    detail.SESNAME,
	}
	return basic, nil
}

//该行业下的所有公司
func IndustryOfAllCompany(swlevelcode string) ([]*finchina.CompanyDetail, error) {
	allCompany, err := finchina.NewCompanyDetail().GetAllCompany(swlevelcode)
	logging.Info("allCompany=====%+v", allCompany)
	return allCompany, err
}

//获取相差时间
func getHourDiffer(start_time, end_time string) int64 {
	var hour int64
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", start_time, time.Local)
	t2, err := time.ParseInLocation("2006-01-02 15:04:05", end_time, time.Local)
	if err == nil && t1.Before(t2) {
		diff := t2.Unix() - t1.Unix() //
		hour = diff / 3600            //返回小时，向下取整
		return hour
	} else {
		return hour
	}
}

//获取post请求的数据
func HttpPostJson(url string, postJson interface{}, _param interface{}) error {
	postJsons, _ := json.Marshal(postJson)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postJsons))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("HTTP POST ERROR| %v", err.Error())
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), _param)
	return err
}

//标准K线结构
type Klist struct {
	Date        int64   `json:"date"`
	Value       int64   `json:"value,omitempty"`
	NSid        int64   `json:"nSID,omitempty"`
	NTime       int64   `json:"nTime,omitempty"`
	NPreCPx     float64 `json:"nPreCPx"`
	NOpenPx     float64 `json:"nOpenPx"`
	NHighPx     float64 `json:"nHighPx"`
	NLastPx     float64 `json:"nLastPx"`
	LlVolume    int64   `json:"llVolume,omitempty"`
	LlValue     int64   `json:"llValue,omitempty"`
	NAvgPx      int64   `json:"nAvgPx,omitempty"`
	NPxChg      float64 `json:"NPxChg,omitempty"`      //int	< 涨跌
	NPxChgRatio float64 `json:"NPxChgRatio,omitempty"` //int	< 涨跌幅*10000
}

//当日K线结构体
type KlistNow struct {
	NSID          int32 //int< 证券ID  4个字节
	NTime         int32 //int< 时间 unix time
	NTradingPhase int32 //unsigned int< 详细见产品交易阶段(TradingPhase) 4个字节
	NPreClosePx   int32 //unsigned int< 昨收价 * 10000
	NOpenPx       int32 //unsigned int< 开盘价 ..
	NHighPx       int32 //unsigned int< 最高价 ..
	NLowPx        int32 //unsigned int< 最低价 ..
	NLastPx       int32 //unsigned int< 最新价 ..
	NHighLimitPx  int32 //unsigned int< 涨停价格(*10000)
	NLowLimitPx   int32 //unsigned int< 跌停价格(*10000)
	LlTradeNum    int64 //long long< 成交笔数   8个字节
	LlVolume      int64 //long long< 成交量
	LlValue       int64 //long long< 成交额(*10000)
	NQuoteSize    int32 //int	< 报价总档数
	NWABidPx      int32 //unsigned int < 加权平均委买均价(*10000)
	NWAOfferPx    int32 //unsigned int< 加权平均委卖均价(*10000)
	LlToBidVol    int64 //long long< 总委买量
	LlToOfferVol  int64 //long long< 总委卖量
	LlInnerVolume int64 //long long< 内盘成交量
	LlOuterVolume int64 //long long	< 外盘成交量
	LlInnerValue  int64 //long long< 内盘成交额
	LlOuterValue  int64 //long long< 外盘成交额
	NPxChg        int32 //int	< 涨跌
	NPxChgRatio   int32 //int	< 涨跌幅*10000
	NPxAmplitude  int32 //int	 < 振幅*10000
	NLiangbi      int32 //int	< 量比*100
	NWeibi        int32 //int	< 委比*10000
	NTurnOver     int32 //int	< 换手率*10000
	NPE           int32 //int	< 动态市盈率*10000
	NPB           int32 //int	< 动态市净率*10000
}

//获取当日K线
func GetKlineNow(KeyNsid string) (Klist, error) {
	tNow := time.Now()
	timeNow := tNow.Format("20060102")
	t, _ := strconv.ParseInt(timeNow, 10, 64)

	//redisData, err := models.RedisStore.Get("hq:st:snap:100600000")
	redisData, err := models.RedisStore.Get("hq:st:snap:" + KeyNsid)

	var kNow Klist
	if err != nil || len(redisData) == 0 {
		kNow.Date = t
		kNow.NPreCPx = 0
		kNow.NHighPx = 0
		kNow.NLastPx = 0
		kNow.LlVolume = 0
		kNow.LlValue = 0
		kNow.NPxChg = 0
		kNow.NPxChgRatio = 0
		logging.Info("==", err.Error())
		return kNow, err
	} else {
		//var kno []KlistNow
		var kn KlistNow
		size := binary.Size(&kn)

		k := &KlistNow{}
		buffer := bytes.NewBuffer([]byte(redisData[0:size]))
		if err := binary.Read(buffer, binary.LittleEndian, k); err != nil && err != io.EOF {
			logging.Info("error===%v", err)
		}
		//kno = append(kno, k)
		logging.Info("kkk====%+v", k)

		//
		kNow.Date = t
		kNow.NPreCPx = float64(k.NPreClosePx)
		kNow.NOpenPx = float64(k.NOpenPx)
		kNow.NHighPx = float64(k.NHighPx)
		kNow.NLastPx = float64(k.NLastPx)
		kNow.LlVolume = k.LlVolume
		kNow.LlValue = k.LlValue
		kNow.NPxChg = float64(k.NPxChg)
		kNow.NPxChgRatio = float64(k.NPxChgRatio)
		return kNow, err
	}
}

//获取历史K线
func GetKlinePast(KeyNsid string) ([]Klist, error) {
	//url := "https://hbmk.0606.com.cn/api/hq/kline?format=json"
	url := "http://47.94.109.175:7280/api/hq/kline?format=json"
	keysid, _ := strconv.Atoi(KeyNsid)

	timestamp := time.Now().Unix()
	tm := time.Unix(timestamp, 0)
	tm1 := tm.Format("20060102")
	tm2, _ := strconv.Atoi(tm1)

	logging.Info("tm1===", tm2)

	posts := map[string]interface{}{
		"SID":       keysid,
		"Type":      1,
		"TimeIndex": tm2,
		"Num":       99,
		"Direct":    0,
	}
	var _param struct {
		Code int64 `json:"code"`
		Data struct {
			SID   int64   `json:"sid"`
			Type  int64   `json:"type"`
			Total int64   `json:"total"`
			Begin int64   `json:"begin"`
			Num   int64   `json:"num"`
			KList []Klist `json:"klist"`
		} `json:"data"`
	}
	err := HttpPostJson(url, posts, &_param)
	if err != nil {
		logging.Info(err.Error())
	}
	return _param.Data.KList, err
}
