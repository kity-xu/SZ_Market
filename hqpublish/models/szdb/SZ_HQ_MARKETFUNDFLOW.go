// 市场分类的资金流向
package szdb

import (
	"liveshow/share/logging"

	"haina.com/share/gocraft/dbr"
	"haina.com/share/models"
)

type SZ_HQ_MARKETFUNDFLOW struct {
	models.Model     `db:"-" `
	MARKETID         int32           //市场ID(暂定:0全市场100000000沪市200000000深市300000000中小400000000创业)
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
	ENTRYDATE        dbr.NullString  //#更新日期
	ENTRYTIME        dbr.NullString  //#更新时间
}

func NewSZ_HQ_MARKETFUNDFLOW() *SZ_HQ_MARKETFUNDFLOW {
	return &SZ_HQ_MARKETFUNDFLOW{
		Model: models.Model{
			CacheKey:  "redis_key",
			TableName: TABLE_SZ_HQ_MARKETFUNDFLOW,
			Db:        models.DBSZ,
		},
	}
}

func (s *SZ_HQ_MARKETFUNDFLOW) GetMarketFundFlow(count int32, marketID int32) ([]SZ_HQ_MARKETFUNDFLOW, error) {
	var capflow []SZ_HQ_MARKETFUNDFLOW

	exps := map[string]interface{}{
		"MARKETID=?": marketID,
	}
	builder := s.Db.Select("*").From(s.TableName)
	_, err := s.SelectWhere(builder, exps).OrderBy("TRADEDATE desc").Limit(uint64(count)).LoadStructs(&capflow)
	if err != nil {
		logging.Error("%s", err.Error())
		return nil, err
	}
	return capflow, err
}
