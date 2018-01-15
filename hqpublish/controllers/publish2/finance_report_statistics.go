package publish2

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"haina.com/share/lib"
	"haina.com/share/logging"

	"haina.com/market/hqpublish/models"

	"fmt"

	"encoding/json"

	"haina.com/market/hqpublish/models/finchina/io_finchina"
	"haina.com/share/garyburd/redigo/redis"
)

type ReportStatistics struct {
	Sid        int         `json:"sid"`
	Period     int         `json:"period"`
	Statistics *Statistics `json:"statistics"`
}

type Statistics struct {
	Sale       int64 `json:"sale"`       // 卖出
	Reduce     int64 `json:"reduce"`     // 减持
	Neutral    int64 `json:"neutral"`    // 中性
	Overweight int64 `json:"overweight"` // 增持
	Buying     int64 `json:"buying"`     // 买入
}

func NewReportStatistics() *ReportStatistics {
	return &ReportStatistics{}
}

func (this *ReportStatistics) POST(c *gin.Context) {
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

func (this *ReportStatistics) PostJson(c *gin.Context) {
	var req struct {
		Sid    int `json:"sid" binding:"required"`
		Period int `json:"period" binding:"required"`
	}
	if err := c.BindJSON(&req); err != nil {
		logging.Debug("BindJson | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}

	logging.Debug("params %+v", req)

	if req.Period < 0 || req.Period > 5 {
		logging.Error("Params Err")
		lib.WriteString(c, 40004, nil)
		return
	}

	this.jsonProcess(c, req.Sid, req.Period)
}
func (this *ReportStatistics) PostPB(c *gin.Context) {
}

func (this *ReportStatistics) jsonProcess(c *gin.Context, sid int, period int) {

	finish := false
	if err := this.ReadCacheJson(sid, period); err == nil {
		lib.WriteString(c, 200, this)
		return
	}
	defer func() {
		if finish {
			this.SaveCacheJson(sid, period)
		}
	}()

	statistics := io_finchina.NewTQ_EXPT_SKSTATRT()

	seg := strconv.Itoa(sid)
	err := statistics.GetSingle(seg[3:], seg[:3], period)
	if err != nil {
		logging.Error("Err | %v", err)
		lib.WriteString(c, 40002, nil)
		return
	}
	logging.Debug("%+v", statistics)

	this.rigger(sid, period, statistics)

	finish = true
	lib.WriteString(c, 200, this)
}

func (this *ReportStatistics) rigger(sid int, period int, statistics *io_finchina.TQ_EXPT_SKSTATRT) *ReportStatistics {
	this.Sid = sid
	this.Period = period
	this.Statistics = &Statistics{
		Sale:       statistics.SELL.Int64,
		Reduce:     statistics.BREDUCENUM.Int64,
		Neutral:    statistics.NEUTRALNUM.Int64,
		Overweight: statistics.ADDBUYNUM.Int64,
		Buying:     statistics.BUY.Int64,
	}

	return this
}

//-------------------------------------------

//const FinanceReportStatisticsKey = "finance:report:statistics:%v:%v"
const FinanceReportStatisticsKey = "finance:report:s:%v:%v"

func (this *ReportStatistics) ReadCacheJson(sid int, period int) error {
	key := fmt.Sprintf(FinanceReportStatisticsKey, sid, period)
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
func (this *ReportStatistics) SaveCacheJson(sid int, period int) error {
	key := fmt.Sprintf(FinanceReportStatisticsKey, sid, period)
	cache, err := json.Marshal(this)
	if err != nil {
		logging.Debug("Json Marshal Err | %v", err)
		return err
	}
	err = models.SetCache(key, models.TTL.FinanceReportStatistics, cache)
	if err != nil {
		logging.Debug("Redis SetCache Err | %v", err)
		return err
	}
	return nil
}
