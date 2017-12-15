// 日资金流向
package szdb

import (
	"haina.com/share/gocraft/dbr"
	"haina.com/share/logging"
	"haina.com/share/models"
)

type SZ_HQ_SECURITYFUNDFLOW struct {
	models.Model     `db:"-" `
	SID              int32           //证券内部ID
	TRADEDATE        int32           //交易日期
	HUGEBUYVALUE     dbr.NullFloat64 //#特大买单成交额(元)
	BIGBUYVALUE      dbr.NullFloat64 //#大买单成交额(元)
	MIDDLEBUYVALUE   dbr.NullFloat64 //#中买单成交额(元)
	SMALLBUYVALUE    dbr.NullFloat64 //#小买单成交额(元)
	HUGESELLVALUE    dbr.NullFloat64 //#特大卖单成交额(元)
	BIGSELLVALUE     dbr.NullFloat64 //#大卖单成交额(元)
	MIDDLESELLVALUE  dbr.NullFloat64 //#中卖单成交额(元)
	SMALLSELLVALUE   dbr.NullFloat64 //#小卖单成交额(元)
	HUGEBUYVOLUME    dbr.NullFloat64 //#特大买单成交量(股)
	BIGBUYVOLUME     dbr.NullFloat64 //#大买单成交量(股)
	MIDDLEBUYVOLUME  dbr.NullFloat64 //#中买单成交量(股)
	SMALLBUYVOLUME   dbr.NullFloat64 //#小买单成交量(股)
	HUGESELLVOLUME   dbr.NullFloat64 //#特大卖单成交量(股)
	BIGSELLVOLUME    dbr.NullFloat64 //#大卖单成交量(股)
	MIDDLESELLVOLUME dbr.NullFloat64 //#中卖单成交量(股)
	SMALLSELLVOLUME  dbr.NullFloat64 //#小卖单成交量(股)
	VALUEOFINFLOW    dbr.NullFloat64 //#最近5日成交总量(股)
	ENTRYDATE        dbr.NullFloat64 //#更新日期
	ENTRYTIME        dbr.NullFloat64 //#更新时间
}

func NewSZ_HQ_SECURITYFUNDFLOW() *SZ_HQ_SECURITYFUNDFLOW {
	return &SZ_HQ_SECURITYFUNDFLOW{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_SZ_HQ_SECURITYFUNDFLOW,
			Db:        models.DBSZ,
		},
	}
}

// 股票历史日资金流向
func (this *SZ_HQ_SECURITYFUNDFLOW) GetHisSecurityFlow(count int32, sid int32) ([]SZ_HQ_SECURITYFUNDFLOW, error) {
	var capflow []SZ_HQ_SECURITYFUNDFLOW

	exps := map[string]interface{}{
		"SID=?": sid,
	}
	builder := this.Db.Select("TRADEDATE,HUGEBUYVALUE,BIGBUYVALUE,HUGESELLVALUE,BIGSELLVALUE").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).OrderBy("TRADEDATE desc").Limit(uint64(count)).LoadStructs(&capflow)
	if err != nil {
		logging.Error("%s", err.Error())
		return capflow, err
	}
	return capflow, err
}

// 股票历史日资金流向
func (this *SZ_HQ_SECURITYFUNDFLOW) GetHisSecurityFlowFull(sid int32) ([]SZ_HQ_SECURITYFUNDFLOW, error) {
	var capflow []SZ_HQ_SECURITYFUNDFLOW

	exps := map[string]interface{}{
		"SID=?": sid,
	}
	builder := this.Db.Select("TRADEDATE,HUGEBUYVALUE,BIGBUYVALUE,MIDDLEBUYVALUE,SMALLBUYVALUE,HUGESELLVALUE,BIGSELLVALUE,MIDDLESELLVALUE,SMALLSELLVALUE").From(this.TableName)
	_, err := this.SelectWhere(builder, exps).OrderBy("TRADEDATE desc").LoadStructs(&capflow)
	if err != nil {
		logging.Error("%s", err.Error())
		return capflow, err
	}
	return capflow, err
}
