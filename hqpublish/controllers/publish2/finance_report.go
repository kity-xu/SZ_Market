// 财务-报表
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

// 利润表
type ProfitTable struct {
	Date      string  `json:"lDate,omitempty"` // 日期
	OpRe      float64 `json:"lOpRe"`           // 营业收入
	OpPr      float64 `json:"lOpPr"`           // 营业利润
	NetProfit float64 `json:"lNetProfit"`      // 净利润
	OpReRate  float64 `json:"lOpReRate"`       // 营业收入增长率
	OpPrRate  float64 `json:"lOpPrRate"`       // 营业利润增长率
	NetRate   float64 `json:"lNetRate"`        // 净利润增长率
}

// 资产负债表
type DebtTable struct {
	Date  string  `json:"lDate,omitempty"`
	ToAs  float64 `json:"lToAs"`  // 资产合计
	TaLb  float64 `json:"lTaLb"`  // 负债合计
	OESET float64 `json:"lOESET"` // 所有者权益合计
}

// 现金流量表
type FlowTable struct {
	Date  string  `json:"lDate,omitempty"`
	NCFOA float64 `json:"lNCFOA"` // 经营活动产生的现金流量净额
	NCIIA float64 `json:"lNCIIA"` // 投资活动产生的现金流量净额
	NCPFA float64 `json:"lNCPFA"` // 筹资活动产生的现金流量净额
}

type FinanceReportRecord struct {
	Profit ProfitTable `json:"profit"`
	Debt   DebtTable   `json:"debt"`
	Flow   FlowTable   `json:"flow"`
}

type FinanceReport struct {
	Count int `json:"count"`
	//Dates string                 `json:"dates"`
	Rows []*FinanceReportRecord `json:"rows"`
}

func NewFinanceReport() *FinanceReport {
	return &FinanceReport{}
}

func (this *FinanceReport) POST(c *gin.Context) {
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

func (this *FinanceReport) PostJson(c *gin.Context) {
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
func (this *FinanceReport) PostPB(c *gin.Context) {
}

func (this *FinanceReport) jsonProcess(c *gin.Context, sid *Sid, count int) {

	finish := false
	if count == 5 {
		if err := this.readCacheJson(sid.Sid); err == nil {
			lib.WriteString(c, 200, this)
			return
		}
		defer func() {
			if finish {
				this.saveCacheJson(sid.Sid)
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
	ls_debt, err := io_finchina.NewLiabilities().GetList(sid.Symbol, sid.Market, 0, sum, 1)
	if err != nil {
		logging.Error("Err | %v", err)
		lib.WriteString(c, 40002, nil)
		return
	}
	ls_flow, err := io_finchina.NewCashflow().GetList(sid.Symbol, sid.Market, 0, sum, 1)
	if err != nil {
		logging.Error("Err | %v", err)
		lib.WriteString(c, 40002, nil)
		return
	}

	if len(ls) != len(ls_debt) || len(ls) != len(ls_flow) {
		logging.Error("Err | incomplete data")
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
		this.rigger(ls, ls_debt, ls_flow, actual)
	} else {
		this.rigger(ls, ls_debt, ls_flow, count)
	}

	finish = true

	lib.WriteString(c, 200, this)
}

func (this *FinanceReport) rigger(ls []io_finchina.Profits, ls_debt []io_finchina.Liabilities, ls_flow []io_finchina.Cashflow, count int) *FinanceReport {
	logging.Debug("rigger len %v, count %v", len(ls), count)

	this.Rows = make([]*FinanceReportRecord, 0, count)

	for i := 0; i < count; i++ {
		//logging.Debug("i %v", i)
		node := &FinanceReportRecord{
			Profit: ProfitTable{
				Date:      ls[i].ENDDATE.String,
				OpRe:      ls[i].BIZINCO.Float64,   // 营业收入
				OpPr:      ls[i].PERPROFIT.Float64, // 营业利润
				NetProfit: ls[i].NETPROFIT.Float64, // 净利润
			},
			Debt: DebtTable{
				Date:  ls_debt[i].ENDDATE.String,
				ToAs:  ls_debt[i].TOTASSET.Float64, // 资产合计
				TaLb:  ls_debt[i].TOTLIAB.Float64,  // 负债合计
				OESET: ls_debt[i].RIGHAGGR.Float64, // 所有者权益合计
			},
			Flow: FlowTable{
				Date:  ls_flow[i].ENDDATE.String,
				NCFOA: ls_flow[i].MANANETR.Float64,       // 经营活动产生的现金流量净额
				NCIIA: ls_flow[i].INVNETCASHFLOW.Float64, // 投资活动产生的现金流量净额
				NCPFA: ls_flow[i].FINNETCFLOW.Float64,    // 筹资活动产生的现金流量净额
			},
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
				if b.BIZINCO.Float64 != 0 {
					node.Profit.OpReRate = (a.BIZINCO.Float64 - b.BIZINCO.Float64) / b.BIZINCO.Float64
				}
				if b.PERPROFIT.Float64 != 0 {
					node.Profit.OpPrRate = (a.PERPROFIT.Float64 - b.PERPROFIT.Float64) / b.PERPROFIT.Float64
				}
				if b.NETPROFIT.Float64 != 0 {
					node.Profit.NetRate = (a.NETPROFIT.Float64 - b.NETPROFIT.Float64) / b.NETPROFIT.Float64
				}
			} /* else {
				logging.Debug("%v %s - %v %s no pass", ayear, amonth, byear, bmonth)
			} */
		}
		this.Rows = append(this.Rows, node)
	}
	this.Count = count
	return this
}

const FinanceReportKey = "finance:report:%v"

func (this *FinanceReport) readCacheJson(sid int) error {
	key := fmt.Sprintf(FinanceReportKey, sid)
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
func (this *FinanceReport) saveCacheJson(sid int) error {
	key := fmt.Sprintf(FinanceReportKey, sid)
	cache, err := json.Marshal(this)
	if err != nil {
		logging.Debug("Json Marshal Err | %v", err)
		return err
	}
	err = models.SetCache(key, models.TTL.FinanceReport, cache)
	if err != nil {
		logging.Debug("Redis SetCache Err | %v", err)
		return err
	}
	return nil
}
