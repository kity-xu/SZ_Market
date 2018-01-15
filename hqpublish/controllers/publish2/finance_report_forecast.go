package publish2

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"haina.com/market/hqpublish/models"

	"encoding/json"
	"fmt"

	"haina.com/market/hqpublish/models/finchina/io_finchina"
	"haina.com/share/garyburd/redigo/redis"
)

type ReportForecast struct {
	Sid int `json:"sid"` // 证券ID

	Tenddate string  `json:"tenddate"` // 预测截止日期（一期）
	Teps     float64 `json:"teps"`     // 一致预期每股收益(一期）
	Tpe      float64 `json:"tpe"`      // 一致预期市盈率(PE)(一期)

	Nenddate string  `json:"nenddate"` // 预测截止日期（二期）
	Neps     float64 `json:"neps"`     // 一致预期每股收益(二期）
	Npe      float64 `json:"npe"`      //  一致预期市盈率(PE)(二期)

	Yanenddate string  `json:"yanenddate"` // 预测截止日期（三期）
	Yaneps     float64 `json:"yaneps"`     // 一致预期每股收益(三期）
	Yanpe      float64 `json:"yanpe"`      //  一致预期市盈率(PE)(三期)
}

func NewReportForecast() *ReportForecast {
	return &ReportForecast{}
}

func (this *ReportForecast) POST(c *gin.Context) {
	replayfmt := c.Query(models.CONTEXT_FORMAT)
	if len(replayfmt) == 0 {
		replayfmt = "json" // 默认
	}

	switch replayfmt {
	case "json":
		this.PostJson(c)
	case "pb":
		this.PostPB(c)
	default:
		return
	}
}

const DefaultForecastCount = 3

func (this *ReportForecast) PostJson(c *gin.Context) {
	var req struct {
		Sid int `json:"sid" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		logging.Debug("BindJson | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	logging.Debug("params %+v", req)

	this.jsonProcess(c, req.Sid)
}
func (this *ReportForecast) PostPB(c *gin.Context) {
}

func (this *ReportForecast) jsonProcess(c *gin.Context, sid int) {

	if err := this.readCacheJson(sid); err == nil {
		lib.WriteString(c, 200, this)
		return
	}
	finish := false
	defer func() {
		if finish {
			this.saveCacheJson(sid)
		}
	}()

	forecast := io_finchina.NewTQ_EXPT_SKSTATN()
	seg := strconv.Itoa(sid)
	err := forecast.GetSingle(seg[3:], seg[:3])
	if err != nil {
		logging.Error("Err | %v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	this.rigger(sid, forecast)

	finish = true
	lib.WriteString(c, 200, this)
}

func (this *ReportForecast) rigger(sid int, forecast *io_finchina.TQ_EXPT_SKSTATN) *ReportForecast {
	this.Sid = sid
	tms := time.Now().Format("2006")
	tm, _ := strconv.Atoi(tms)
	tm1, _ := strconv.Atoi(forecast.TENDDATE.String)
	tm2, _ := strconv.Atoi(forecast.NENDDATE.String)
	tm3, _ := strconv.Atoi(forecast.YANENDDATE.String)

	if tm == tm1/10000 { //今年
		this.Tenddate = forecast.TENDDATE.String
		this.Teps = forecast.TEPS.Float64
		this.Tpe = forecast.TPE.Float64

		this.Nenddate = strconv.Itoa((tm+1)*10000 + 1231)
		this.Neps = forecast.NEPS.Float64
		this.Npe = forecast.NPE.Float64

		this.Yanenddate = strconv.Itoa((tm+2)*10000 + 1231)
		this.Yaneps = forecast.YANEPS.Float64
		this.Yanpe = forecast.YANPE.Float64
	} else if tm == tm2/10000 || tm3 == 0 { //明年
		this.Tenddate = forecast.NENDDATE.String
		this.Teps = forecast.NEPS.Float64
		this.Tpe = forecast.NPE.Float64

		this.Nenddate = strconv.Itoa((tm+1)*10000 + 1231)
		this.Neps = forecast.YANEPS.Float64
		this.Npe = forecast.YANPE.Float64

		this.Yanenddate = strconv.Itoa((tm+2)*10000 + 1231)
		this.Yaneps = 0
		this.Yanpe = 0
	} else if tm == tm3/10000 || tm2 == 0 { //后年
		this.Tenddate = forecast.YANENDDATE.String
		this.Teps = forecast.YANEPS.Float64
		this.Tpe = forecast.YANPE.Float64

		this.Nenddate = strconv.Itoa((tm+1)*10000 + 1231)
		this.Neps = 0
		this.Npe = 0

		this.Yanenddate = strconv.Itoa((tm+2)*10000 + 1231)
		this.Yaneps = 0
		this.Yanpe = 0
	} else {
		this.Tenddate = tms + "1231"
		this.Teps = 0
		this.Tpe = 0

		this.Nenddate = strconv.Itoa((tm+1)*10000 + 1231)
		this.Neps = 0
		this.Npe = 0

		this.Yanenddate = strconv.Itoa((tm+2)*10000 + 1231)
		this.Yaneps = 0
		this.Yanpe = 0
	}
	return this
}

//-------------------------------------------

//const FinanceReportForecastKey = "finance:report:forecast:%v"
const FinanceReportForecastKey = "finance:report:f:%v"

func (this *ReportForecast) readCacheJson(sid int) error {
	key := fmt.Sprintf(FinanceReportForecastKey, sid)
	cache, err := models.GetCache(key)
	if err != nil {
		if err == redis.ErrNil {
			logging.Info("Redis GetCache not found | %v", key)
			return err
		}
		logging.Debug("Redis GetCache Err | %v", err)
		return err
	}
	logging.Debug("hit redis cache %v", key)
	err = json.Unmarshal(cache, this)
	if err != nil {
		logging.Debug("Json Unmarshal Err | %v", err)
		return err
	}
	return nil
}
func (this *ReportForecast) saveCacheJson(sid int) error {
	key := fmt.Sprintf(FinanceReportForecastKey, sid)
	cache, err := json.Marshal(this)
	if err != nil {
		logging.Debug("Json Marshal Err | %v", err)
		return err
	}
	err = models.SetCache(key, models.TTL.FinanceReportForecast, cache)
	if err != nil {
		logging.Debug("Redis SetCache Err | %v", err)
		return err
	}
	return nil
}
