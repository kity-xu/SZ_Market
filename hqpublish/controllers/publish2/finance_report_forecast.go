package publish2

import (
	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"haina.com/market/hqpublish/models"

	"fmt"

	"encoding/json"

	. "haina.com/market/hqpublish/controllers"
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

	s := NewSid(req.Sid)

	this.jsonProcess(c, s)
}
func (this *ReportForecast) PostPB(c *gin.Context) {
}

func (this *ReportForecast) jsonProcess(c *gin.Context, sid *Sid) {

	if err := this.readCacheJson(sid.Sid); err == nil {
		lib.WriteString(c, 200, this)
		return
	}
	finish := false
	defer func() {
		if finish {
			this.saveCacheJson(sid.Sid)
		}
	}()

	forecast := io_finchina.NewTQ_EXPT_SKSTATN()
	err := forecast.GetSingle(sid.Symbol, sid.Market)
	if err != nil {
		logging.Error("Err | %v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	this.rigger(sid, forecast)

	finish = true
	lib.WriteString(c, 200, this)
}

func (this *ReportForecast) rigger(sid *Sid, forecast *io_finchina.TQ_EXPT_SKSTATN) *ReportForecast {
	this.Sid = sid.Sid

	this.Tenddate = forecast.TENDDATE.String
	this.Teps = forecast.TEPS.Float64
	this.Tpe = forecast.TPE.Float64

	this.Nenddate = forecast.NENDDATE.String
	this.Neps = forecast.NEPS.Float64
	this.Npe = forecast.NPE.Float64

	this.Yanenddate = forecast.YANENDDATE.String
	this.Yaneps = forecast.YANEPS.Float64
	this.Yanpe = forecast.YANPE.Float64

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
