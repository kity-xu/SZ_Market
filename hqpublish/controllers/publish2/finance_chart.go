// 财务-图表
package publish2

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"haina.com/market/hqpublish/models"
	"haina.com/market/hqpublish/models/finchina/io_finchina"
	"haina.com/share/logging"

	. "haina.com/market/hqpublish/controllers"
	"haina.com/share/lib"

	"haina.com/share/garyburd/redigo/redis"
)

// 每股收益
type NodeEPS struct {
	Date string  `json:"date"` // 年月日(yyyymmdd)
	Eps  float64 `json:"eps"`  // 每股收益
	Rate float64 `json:"rate"`
}

// 营业总收入
type NodeIncome struct {
	Date   string  `json:"date"`   // 年月日(yyyymmdd)
	Income float64 `json:"income"` // 营业总收入
	Rate   float64 `json:"rate"`
}

// 净利润
type NodeNetprofit struct {
	Date      string  `json:"date"`      // 年月日(yyyymmdd)
	Netprofit float64 `json:"netprofit"` // 净利润
	Rate      float64 `json:"rate"`
}

type FinanceChart struct {
	Count     int              `json:"count"`
	EPS       []*NodeEPS       `json:"eps"`
	Income    []*NodeIncome    `json:"income"`
	Netprofit []*NodeNetprofit `json:"netprofit"`
}

func NewFinanceChart() *FinanceChart {
	return &FinanceChart{}
}

func (this *FinanceChart) POST(c *gin.Context) {
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

func (this *FinanceChart) PostJson(c *gin.Context) {
	var req struct {
		Sid   int `json:"sid" binding:"required"`
		Count int `json:"count"`
	}
	if err := c.BindJSON(&req); err != nil {
		logging.Debug("BindJson | %v", err)
		lib.WriteString(c, 40004, nil)
		return
	}
	// 默认5条
	if req.Count == 0 {
		req.Count = 5
	}
	s := NewSid(req.Sid)

	this.jsonProcess(c, s, req.Count)
}

func (this *FinanceChart) PostPB(c *gin.Context) {
}

func (this *FinanceChart) jsonProcess(c *gin.Context, sid *Sid, count int) {

	finish := false
	if count == 5 {
		var op RedisCacheOperator = this
		if err := op.ReadCacheJson(sid.Sid); err == nil {
			lib.WriteString(c, 200, this)
			return
		}
		defer func() {
			if finish {
				op.SaveCacheJson(sid.Sid)
			}
		}()
	}

	sum := 4 + count

	logging.Debug("count %v, sum %v", count, sum)

	ls, err := io_finchina.NewProfits().GetList(sid.Symbol, sid.Market, 0, sum, 1)
	if err != nil {
		logging.Error("Err | %v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	// 计算实际条数，如果数据库里条数不够，计算修正
	actual := count // 实际数量
	if len(ls) < sum {
		if len(ls) < count {
			actual = len(ls) // 实际数量修正
		}
	}

	if actual < count {
		this.rigger(ls, actual)
	} else {
		this.rigger(ls, count)
	}
	finish = true

	lib.WriteString(c, 200, this)
}

func (this *FinanceChart) rigger(ls []io_finchina.Profits, count int) *FinanceChart {
	//logging.Debug("rigger len %v, count %v", len(ls), count)

	this.EPS = make([]*NodeEPS, 0, count)
	this.Income = make([]*NodeIncome, 0, count)
	this.Netprofit = make([]*NodeNetprofit, 0, count)

	for i := 0; i < count; i++ {
		//logging.Debug("i %v", i)
		eps := &NodeEPS{
			Date: ls[i].ENDDATE.String,
			Eps:  ls[i].BASICEPS.Float64,
		}
		income := &NodeIncome{
			Date:   ls[i].ENDDATE.String,
			Income: ls[i].BIZTOTINCO.Float64,
		}
		netprofit := &NodeNetprofit{
			Date:      ls[i].ENDDATE.String,
			Netprofit: ls[i].NETPROFIT.Float64,
		}

		if i+4 < len(ls) && len(ls[i].ENDDATE.String) > 7 && len(ls[i+4].ENDDATE.String) > 7 {

			a := ls[i]
			b := ls[i+4]

			ayear := PackAtoi(a.ENDDATE.String[:4])
			byear := PackAtoi(b.ENDDATE.String[:4])
			amonth := a.ENDDATE.String[4:6]
			bmonth := b.ENDDATE.String[4:6]

			if ayear-1 == byear && amonth == bmonth {
				//logging.Debug("%v %s - %v %s to pass", ayear, amonth, byear, bmonth)
				if b.BASICEPS.Float64 != 0 {
					eps.Rate = (a.BASICEPS.Float64 - b.BASICEPS.Float64) / b.BASICEPS.Float64
				}
				if b.BIZTOTINCO.Float64 != 0 {
					income.Rate = (a.BIZTOTINCO.Float64 - b.BIZTOTINCO.Float64) / b.BIZTOTINCO.Float64
				}
				if b.NETPROFIT.Float64 != 0 {
					netprofit.Rate = (a.NETPROFIT.Float64 - b.NETPROFIT.Float64) / b.NETPROFIT.Float64
				}
			} /* else {
				logging.Debug("%v %s - %v %s no pass", ayear, amonth, byear, bmonth)
			} */
		}
		this.EPS = append(this.EPS, eps)
		this.Income = append(this.Income, income)
		this.Netprofit = append(this.Netprofit, netprofit)
	}
	this.Count = count
	return this
}

//-------------------------------------------

const FinanceChartKey = "finance:chart:%v"

func (this *FinanceChart) ReadCacheJson(sid int) error {
	key := fmt.Sprintf(FinanceChartKey, sid)
	cache, err := models.GetCache(key)
	if err != nil {
		if err == redis.ErrNil {
			logging.Info("Redis GetCache not found | %v", key)
			return err
		}
		logging.Debug("Redis GetCache Err | %v", err)
		return err
	}
	err = json.Unmarshal(cache, this)
	if err != nil {
		logging.Debug("Json Unmarshal Err | %v", err)
		return err
	}
	return nil
}
func (this *FinanceChart) SaveCacheJson(sid int) error {
	key := fmt.Sprintf(FinanceChartKey, sid)
	cache, err := json.Marshal(this)
	if err != nil {
		logging.Debug("Json Marshal Err | %v", err)
		return err
	}
	err = models.SetCache(key, 3600, cache)
	if err != nil {
		logging.Debug("Redis SetCache Err | %v", err)
		return err
	}
	return nil
}
